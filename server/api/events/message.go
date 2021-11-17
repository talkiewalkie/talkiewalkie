package events

import (
	"bytes"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	uuid2 "github.com/satori/go.uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"github.com/talkiewalkie/talkiewalkie/pkg/slices"
	"github.com/talkiewalkie/talkiewalkie/repositories"
)

func OnNewMessage(
	components *common.Components,
	me *models.User,
	event *pb.Event,
) (
	ev *pb.Event,
	dbEv *models.Event,
	err error,
) {
	nm := event.GetSentNewMessage()

	var conv *models.Conversation
	switch nm.Conversation.(type) {
	case *pb.Event_SentNewMessage_ConvUuid:
		uid, err := uuid2.FromString(nm.GetConvUuid())
		if err != nil {
			return nil, nil, status.Errorf(codes.InvalidArgument, "bad conversation uuid: %+v", err)
		}

		conv, err = components.ConversationRepository.ByUuid(uid)
		if err != nil {
			return nil, nil, status.Errorf(codes.Internal, "could not find conversation: %+v", err)
		}

		if ok, err := components.ConversationHasAccess(me, conv); !ok || err != nil {
			return nil, nil, status.Errorf(codes.PermissionDenied, "can't send message to this conversation: (err = %+v)", err)
		}

	case *pb.Event_SentNewMessage_NewConversation:
		allUuids := append(nm.GetNewConversation().UserUuids, me.UUID.String())
		var uuids slices.UuidSlice
		for _, uidStr := range allUuids {
			uid, err := uuid2.FromString(uidStr)
			if err != nil {
				return nil, nil, status.Error(codes.InvalidArgument, err.Error())
			}

			uuids = append(uuids, uid)
		}
		uuids = uuids.Unique()

		conv, err = getOrCreateConvForUsers(components, uuids, nm.GetNewConversation().GetTitle())
		if err != nil {
			return nil, nil, err
		}

	default:
		return nil, nil, status.Errorf(codes.Internal, "received unknown conversation input: %q", nm.Conversation)
	}

	components.ResetEntityStores(components.Ctx)

	//
	// -- MAKE MESSAGE
	//
	var msg *models.Message
	switch nm.Message.Content.(type) {
	case *pb.MessageSendInput_TextMessage:
		text := nm.Message.GetTextMessage().Content
		msg = &models.Message{
			Type:           models.MessageTypeText,
			Text:           null.StringFrom(text),
			AuthorID:       null.IntFrom(me.ID),
			ConversationID: conv.ID,
			CreatedAt:      time.Now(),
		}
	case *pb.MessageSendInput_VoiceMessage:
		vm := nm.Message.GetVoiceMessage()

		pbTranscript, err := proto.Marshal(vm.SiriTranscript)
		if err != nil {
			return nil, nil, status.Error(codes.InvalidArgument, "could not serialize transcript")
		}

		blobUuid, err := components.StorageClient.Upload(components.Ctx, bytes.NewReader(vm.RawContent))
		if err != nil {
			return nil, nil, status.Error(codes.Internal, fmt.Sprintf("could not upload voice message content: %+v", err))
		}

		asset := &models.Asset{
			// TODO: filename schema, possibly [authorUuid]-[convUuid]-[timestamp].ogg ?
			FileName: "",
			// TODO: normalize audio with audio service
			MimeType: "audio/*",
			Bucket:   null.StringFrom(components.StorageClient.DefaultBucket()),
			BlobName: null.StringFrom(blobUuid.String()),
		}
		if err = asset.Insert(components.Ctx, components.Db, boil.Infer()); err != nil {
			return nil, nil, status.Errorf(codes.Internal, "could not register asset in db: %+v", err)
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
		return nil, nil, status.Error(codes.Internal, "unknown content type!")
	}

	if err := msg.Insert(components.Ctx, components.Db, boil.Infer()); err != nil {
		return nil, nil, status.Errorf(codes.Internal, "could not insert message: %+v", err)
	}

	//
	// -- STORE EVENTS IN DB FOR ALL USERS
	//
	ucs, err := components.ConversationUsers(conv)
	if err != nil {
		return nil, nil, err
	}
	participants, err := components.UserRepository.FromUserConversations([][]*models.UserConversation{ucs})

	dbEvInputs := models.EventSlice{}
	for _, p := range participants {
		dbEvInputs = append(dbEvInputs, &models.Event{
			Type:        models.EventTypeNewMessage,
			RecipientID: p.ID,
			MessageID:   null.IntFrom(msg.ID),
		})
	}
	dbEvs, err := BatchInsert(components.Ctx, components.Db, dbEvInputs)
	if err != nil {
		return nil, nil, status.Errorf(codes.Internal, "could not insert new events: %+v", err)
	}

	//
	// -- SEND EVENTS TO CONNECTED USERS
	//
	pbNewEvents, err := EventsToProto(components, dbEvs)
	if err != nil {
		return nil, nil, err
	}
	userToEvents := dbEvs.GroupByRecipientIDs()

	for _, uc := range ucs {
		user, err := components.UserRepository.ById(uc.UserID)
		if err != nil {
			return nil, nil, err
		}

		topic := repositories.UserPubSubTopic(user)

		userEvents := userToEvents[uc.UserID]
		for _, dbevent := range userEvents {
			pbEvent := pbNewEvents.UuidMap()[dbevent.UUID.String()]
			if uc.UserID == me.ID {
				pbEvent.LocalUuid = event.LocalUuid
			}

			if err = components.PubSubClient.Publish(topic, pbEvent); err != nil {
				log.Printf("failed to notify user channel: %+v", err)
			} else {
				log.Printf("sent message on pubsub[%s]!", topic)
			}
		}

	}

	//
	// -- OUTPUT PROTO
	//
	var mydbEvent *models.Event
	for _, event := range dbEvs {
		if event.RecipientID == me.ID {
			mydbEvent = event
		}
	}

	mypbEvents, err := EventsToProto(components, []*models.Event{mydbEvent})
	if err != nil {
		return nil, nil, err
	}

	mypbEvent := mypbEvents[0]
	mypbEvent.LocalUuid = event.LocalUuid
	return mypbEvent, mydbEvent, nil
}

func getOrCreateConvForUsers(
	components *common.Components,
	uuids slices.UuidSlice,
	title string,
) (conv *models.Conversation, err error) {
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
			conv, err = components.ConversationRepository.ById(convUCs[0].ConversationID)
			if err != nil {
				return nil, err
			}
			return conv, nil
		}
	}

	conv = &models.Conversation{Name: null.StringFrom(title)}
	if err = conv.Insert(components.Ctx, components.Db, boil.Infer()); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("could not create new conversation: %+v", err))
	}

	q := sq.Insert(models.TableNames.UserConversation).Columns(models.UserConversationColumns.UserID, models.UserConversationColumns.ConversationID)
	for _, id := range recipientIds {
		q = q.Values(id, conv.ID)
	}
	query, args, _ := q.PlaceholderFormat(sq.Dollar).ToSql()
	if _, err = components.Db.Exec(query, args...); err != nil {
		return nil, status.Errorf(codes.Internal, "could not insert new members to conv: %+v", err)
	}
	components.UserConversationRepository.Clear()

	return conv, nil
}
