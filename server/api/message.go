package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	uuid2 "github.com/satori/go.uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/talkiewalkie/talkiewalkie/clients"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"github.com/talkiewalkie/talkiewalkie/pkg/slices"
	"github.com/talkiewalkie/talkiewalkie/repositories"
)

type MessageService struct {
}

var _ pb.MessageServiceServer = MessageService{}

func NewMessageService() MessageService {
	return MessageService{}
}

func (ms MessageService) Incoming(_ *pb.Empty, server pb.MessageService_IncomingServer) error {
	components, me, err := WithAuthedContext(server.Context())
	if err != nil {
		return err
	}

	topic := repositories.UserPubSubTopic(me)
	log.Printf("established websocket connection [%s]", topic)

	listener, unlisten, err := components.PubSubClient.Subscribe(topic)

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
		psEvent := <-listener
		// On each new event we restore the request cache to mitigate caching issues
		components.ResetEntityStores(server.Context())

		var payload clients.PubSubEvent
		if err := json.Unmarshal([]byte(psEvent.Extra), &payload); err != nil {
			log.Printf("failed to parse pubsub payload on topic '%s': %s", topic, psEvent.Extra)
		}

		switch payload.Type {
		case clients.PubSubEventTypeNewMessage:
			var msg clients.NewMessageEvent
			if err = json.Unmarshal([]byte(psEvent.Extra), &msg); err != nil {
				log.Printf("failed to unmarshal pubsub event: %+v", err)
				break
			}

			newMsg, err := components.MessageRepository.ByUuid(msg.MessageUuid)
			if err != nil {
				return err
			}

			pbMsgs, err := components.MessagesToProto([]*models.Message{newMsg})
			if err != nil {
				return status.Errorf(codes.Internal, "could not transform message to protobuf: %+v", err)
			}

			err = server.Send(pbMsgs[0])
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
	components, me, err := WithAuthedContext(ctx)
	if err != nil {
		return nil, err
	}

	var conv *models.Conversation

	switch input.Recipients.(type) {
	case *pb.MessageSendInput_ConvUuid:
		uid, err := uuid2.FromString(input.GetConvUuid())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		dbConv, err := components.ConversationRepository.ByUuid(uid)
		if err != nil {
			return nil, err
		}

		conv = dbConv

	case *pb.MessageSendInput_RecipientUuids:
		allUuids := append(input.GetRecipientUuids().Uuids, me.UUID.String())
		var uuids slices.Uuid2UUIDSlice
		for _, uidStr := range allUuids {
			uid, err := uuid2.FromString(uidStr)
			if err != nil {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}

			uuids = append(uuids, uid)
		}
		uuids = uuids.UniqueBy(func(uuid uuid2.UUID) interface{} { return uuid.String() })

		recipients, err := components.UserRepository.ByUuids(uuids...)
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("could not find recipients: %+v", err))
		}

		var recipientIds slices.IntSlice
		recipientIds = recipients.Ids()

		ugs, err := components.UserConversationRepository.ByUserIds(recipientIds...)
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("could not find recipients conversations: %+v", err))
		}

		for _, convUCs := range ugs {
			slice := models.UserConversationSlice(convUCs)
			participantIds := slice.UserIDs()

			if recipientIds.SameAs(participantIds) {
				dbConv, err := components.ConversationRepository.ById(convUCs[0].ConversationID)
				if err != nil {
					return nil, err
				}

				conv = dbConv
				break
			}
		}

		// TODO: use a batch insert method like COPY which would make things faster
		tx, err := components.Db.BeginTx(ctx, nil)
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
			for _, id := range recipientIds {
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

	if ok, err := components.ConversationHasAccess(me, conv); !ok || err != nil {
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

		blobUuid, err := components.StorageClient.Upload(ctx, bytes.NewReader(vm.RawContent))
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("could not upload voice message content: %+v", err))
		}

		asset := &models.Asset{
			// TODO: filename schema, possibly [authorUuid]-[convUuid]-[timestamp].ogg ?
			FileName: "",
			// TODO: normalize audio with audio service
			MimeType: "audio/*",
			Bucket:   null.StringFrom(components.StorageClient.DefaultBucket()),
			BlobName: null.StringFrom(blobUuid.String()),
		}
		if err = asset.Insert(ctx, components.Db, boil.Infer()); err != nil {
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

	if err = msg.Insert(ctx, components.Db, boil.Infer()); err != nil {
		return nil, status.Errorf(codes.Internal, "could not insert message: %+v", err)
	}

	ucs, err := components.ConversationUsers(conv)
	for _, uc := range ucs {
		if uc.UserID == me.ID {
			continue
		}

		user, err := components.UserRepository.ById(uc.UserID)
		if err != nil {
			return nil, err
		}

		topic := repositories.UserPubSubTopic(user)
		err = components.PubSubClient.Publish(topic, clients.NewMessageEvent{
			PubSubEvent: clients.PubSubEvent{Type: clients.PubSubEventTypeNewMessage, Timestamp: time.Now()},
			MessageUuid: msg.UUID,
		})
		if err != nil {
			log.Printf("failed to notify user channel: %+v", err)
		} else {
			log.Printf("sent message on pubsub[%s]!", topic)
		}
	}

	pbMsgs, err := components.MessagesToProto([]*models.Message{msg})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to marshal message as pb: %+v", err)
	}

	return pbMsgs[0], nil
}
