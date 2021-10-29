package coco

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	uuid2 "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/entities"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"reflect"
	"sort"
	"sync"
	"time"
)

type MessageService struct {
	*common.Components
}

var _ pb.MessageServiceServer = MessageService{}

func NewMessageService(c *common.Components) MessageService {
	return MessageService{Components: c}
}

func (ms MessageService) Incoming(_ *pb.Empty, server pb.MessageService_IncomingServer) error {
	u, err := common.GetUser(server.Context())

	topic := entities.UserPubSubTopic(u)
	log.Printf("established websocket connection [%s]", topic)

	listener, unlisten, err := ms.PgPubSub.Subscribe(topic)

	if err != nil {
		if err := unlisten(); err != nil {
			log.Printf("could not stop listening to topic '%s': %+v", topic, err)
		}
		return status.Error(codes.Internal, fmt.Sprintf("could not subscribe to pubsub topic: %+v", err))
	}

	defer func() {
		if err := unlisten(); err != nil {
			log.Printf("could not stop listening to topic '%s': %+v", topic, err)
		}
	}()

	for {
		psEvent := <-listener.Notify
		var payload common.PubSubEvent
		if err := json.Unmarshal([]byte(psEvent.Extra), &payload); err != nil {
			log.Printf("failed to parse pubsub payload on topic '%s': %s", topic, psEvent.Extra)
		}

		switch payload.Type {
		case common.PubSubEventTypeNewMessage:
			var msg common.NewMessageEvent
			if err = json.Unmarshal([]byte(psEvent.Extra), &msg); err != nil {
				log.Printf("failed to unmarshal pubsub event: %+v", err)
				break
			}

			newMsg, err := models.Messages(
				models.MessageWhere.UUID.EQ(msg.MessageUuid),
				qm.Load(models.MessageRels.Conversation),
				qm.Load(models.MessageRels.Author),
				qm.Load(models.MessageRels.RawAudio),
			).One(server.Context(), ms.Db)
			if err != nil {
				return status.Errorf(codes.Internal, "could not fetch message from db: %+v", err)
			}

			pbMsg, err := entities.MessageToPb(newMsg, ms.Components)
			if err != nil {
				return status.Errorf(codes.Internal, "could not transform message to protobuf: %+v", err)
			}

			err = server.Send(pbMsg)
			if err != nil {
				log.Printf("failed to send message in server stream: %+v", err)
				break
			} else {
				log.Printf("recovered message from pubsub[%s] and forwarded it with success", topic)
			}

		default:
			log.Printf("received unknown pubsub message on topic '%s': (%T) %+v", topic, payload, payload)
		}
	}
}

func (ms MessageService) Send(ctx context.Context, input *pb.MessageSendInput) (*pb.Message, error) {
	me, err := common.GetUser(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	var conv *models.Conversation

	switch input.Recipients.(type) {
	case *pb.MessageSendInput_ConvUuid:
		uid, err := uuid2.FromString(input.GetConvUuid())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		conv, err = models.Conversations(
			models.ConversationWhere.UUID.EQ(uid),
			qm.Load(qm.Rels(models.ConversationRels.UserConversations, models.UserConversationRels.User)),
		).One(ctx, ms.Db)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not fetch conversation: %+v", err)
		}

	case *pb.MessageSendInput_RecipientUuids:
		allUuids := append(input.GetRecipientUuids().Uuids, me.UUID.String())
		uuids := []uuid2.UUID{}
		for _, uidStr := range allUuids {

			uid, err := uuid2.FromString(uidStr)
			if err != nil {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}

			redundant := false
			for _, existingUid := range uuids {
				if existingUid == uid {
					redundant = true
					break
				}
			}
			if !redundant {
				uuids = append(uuids, uid)
			}
		}

		recipients, err := models.Users(models.UserWhere.UUID.IN(uuids)).All(ctx, ms.Db)
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("could not find recipients: %+v", err))
		}
		if len(recipients) != len(uuids) {
			return nil, status.Error(codes.Internal, fmt.Sprintf("some users where not found: provided %d unique uids, found %d users in db", len(uuids), len(recipients)))
		}

		ids := []int{me.ID}
		for _, recipient := range recipients {
			if recipient.ID != me.ID {
				ids = append(ids, recipient.ID)
			}
		}

		ugs, err := models.UserConversations(
			models.UserConversationWhere.UserID.EQ(me.ID),
			qm.Load(qm.Rels(models.UserConversationRels.Conversation, models.ConversationRels.UserConversations)),
		).All(ctx, ms.Db)
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("could not find recipients conversations: %+v", err))
		}

		sort.Ints(ids)
		for _, ug := range ugs {
			conversationIds := []int{}
			for _, ug := range ug.R.Conversation.R.UserConversations {
				// TODO: somehow traversing the dependencies brings redundant rows, e.g. the list we're iterating on can
				// 		yield [115, 115, 116] as user ids.
				//      Relevant issue https://github.com/volatiletech/sqlboiler/issues/457
				redundant := false
				for _, id := range conversationIds {
					if ug.UserID == id {
						redundant = true
						break
					}
				}
				if !redundant {
					conversationIds = append(conversationIds, ug.UserID)
				}
			}
			sort.Ints(conversationIds)
			if reflect.DeepEqual(conversationIds, ids) {
				conv = ug.R.Conversation
				break
			}
		}

		// TODO: use a batch insert method like COPY which would make things faster
		tx, err := ms.Db.BeginTx(ctx, nil)
		if err != nil {
			tx.Rollback()
			return nil, status.Error(codes.Internal, fmt.Sprintf("could not start transaction: %+v", err))
		}

		if conv == nil {
			title := null.StringFrom(input.GetRecipientUuids().Title)
			if title.String == "" {
				title = null.StringFromPtr(nil)
			}

			newConversation := models.Conversation{Name: title}
			if err = newConversation.Insert(ctx, tx, boil.Infer()); err != nil {
				tx.Rollback()
				return nil, status.Error(codes.Internal, fmt.Sprintf("could not create new conversation: %+v", err))
			}

			errs := make(chan error, 1)
			var wg sync.WaitGroup
			for _, id := range ids {
				wg.Add(1)
				uid := id
				go func() {
					ug := models.UserConversation{
						UserID:         uid,
						ConversationID: newConversation.ID,
					}
					if err := ug.Insert(ctx, tx, boil.Infer()); err != nil {
						errs <- err
					}
					wg.Done()
				}()
			}
			wg.Wait()
			close(errs)
			for err := range errs {
				tx.Rollback()
				return nil, status.Error(codes.Internal, fmt.Sprintf("could not add recipient to new conversation: %+v", err))
			}
			conv = &newConversation
		}
	}

	if ok, err := entities.CanAccessConversation(conv, me); !ok || err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "can't send message to this conversation: (err = %+v)", err)
	}

	var msg *models.Message
	switch input.Content.(type) {
	case *pb.MessageSendInput_TextMessage:
		text := input.GetTextMessage().Content
		msg = &models.Message{
			Type:           models.MessageTypeText,
			Text:           null.StringFrom(text),
			AuthorID:       null.IntFrom(me.ID),
			ConversationID: conv.ID,
			CreatedAt:      time.Now(),
		}
	case *pb.MessageSendInput_VoiceMessage:
		vm := input.GetVoiceMessage()

		pbTranscript, err := proto.Marshal(vm.SiriTranscript)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "could not serialize transcript")
		}

		blobUuid, err := ms.Storage.Upload(ctx, bytes.NewReader(vm.RawContent))
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("could not upload voice message content: %+v", err))
		}

		asset := &models.Asset{
			// TODO: filename schema, possibly [authorUuid]-[convUuid]-[timestamp].ogg ?
			FileName: "",
			// TODO: normalize audio with audio service
			MimeType: "audio/*",
			Bucket:   null.StringFrom(ms.Storage.DefaultBucket()),
			BlobName: null.StringFrom(blobUuid.String()),
		}
		if err = asset.Insert(ctx, ms.Db, boil.Infer()); err != nil {
			return nil, status.Errorf(codes.Internal, "could not register asset in db: %+v", err)
		}

		msg = &models.Message{
			Type:           models.MessageTypeVoice,
			SiriTranscript: null.BytesFrom(pbTranscript),
			RawAudioID:     null.IntFrom(asset.ID),

			AuthorID:       null.IntFrom(me.ID),
			ConversationID: conv.ID,
			CreatedAt:      time.Now(),
		}
	default:
		return nil, status.Error(codes.Internal, "unknown content type!")
	}

	if err = msg.Insert(ctx, ms.Db, boil.Infer()); err != nil {
		return nil, status.Errorf(codes.Internal, "could not insert message: %+v", err)
	}

	for _, uc := range conv.R.UserConversations {
		if uc.R.User.UUID == me.UUID {
			continue
		}
		topic := entities.UserPubSubTopic(uc.R.User)
		err = ms.PgPubSub.Publish(topic, common.NewMessageEvent{
			PubSubEvent: common.PubSubEvent{Type: common.PubSubEventTypeNewMessage, Timestamp: time.Now()},
			MessageUuid: msg.UUID,
		})
		if err != nil {
			log.Printf("failed to notify user channel: %+v", err)
		} else {
			log.Printf("sent message on pubsub[%s]!", topic)
		}
	}

	// TODO: remove this and find a way to prime the messageR struct when we already have the objects in order to avoid pointless roundtrips.
	if err = msg.L.LoadConversation(ctx, ms.Db, true, msg, qm.Comment("")); err != nil {
		return nil, status.Errorf(codes.Internal, "failed ot load converstion: %+v", err)
	}
	if err = msg.L.LoadAuthor(ctx, ms.Db, true, msg, qm.Comment("")); err != nil {
		return nil, status.Errorf(codes.Internal, "failed ot load converstion: %+v", err)
	}
	if err = msg.L.LoadRawAudio(ctx, ms.Db, true, msg, qm.Comment("")); err != nil {
		return nil, status.Errorf(codes.Internal, "failed ot load converstion: %+v", err)
	}

	pbm, err := entities.MessageToPb(msg, ms.Components)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to marshal message as pb: %+v", err)
	}

	return pbm, nil
}
