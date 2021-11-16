package api

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/talkiewalkie/talkiewalkie/clients"
	"github.com/talkiewalkie/talkiewalkie/pkg/slices"
	"github.com/talkiewalkie/talkiewalkie/repositories"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"log"
	"sort"
	"time"

	sq "github.com/Masterminds/squirrel"
	uuid2 "github.com/satori/go.uuid"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/pb"
)

type EventService struct{}

func (e EventService) Connect(server pb.EventService_ConnectServer) error {
	panic("implement me")
}

func (e EventService) Sync(ctx context.Context, sync *pb.UpSync) (*pb.DownSync, error) {
	components, me, err := WithAuthedContext(ctx)
	if err != nil {
		return nil, err
	}

	var out []*pb.Event
	var catchUp []*models.Event
	if sync.LastEventUuid == "" {
	} else {

		evUid, err := uuid2.FromString(sync.LastEventUuid)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "couldn't parse uuid: %+v", evUid)
		}

		lastEvent, err := models.Events(models.EventWhere.UUID.EQ(evUid)).One(ctx, components.Db)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		catchUp, err := models.Events(
			models.EventWhere.RecipientID.EQ(me.ID),
			models.EventWhere.ID.GT(lastEvent.ID),
		).All(ctx, components.Db)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		out, err = EventsToProto(components, catchUp)
		if err != nil {
			return nil, err
		}
	}
	in := []*models.Event{}
	for _, event := range sync.Events {
		switch event.Content.(type) {
		case *pb.Event_SentNewMessage_:
			newEvents, err := OnNewMessage(components, me, event.GetSentNewMessage())
			if err != nil {
				return nil, err
			}

			myNewEvents := []*models.Event{}
			for _, newEvent := range newEvents {
				if newEvent.RecipientID == me.ID {
					myNewEvents = append(myNewEvents, newEvent)
				}
			}
			in = append(in, myNewEvents...)

		case *pb.Event_ReceivedNewMessage_:
			break

		case *pb.Event_DeletedMessage_:
			break

		case *pb.Event_ChangedPicture_:
			break

		case *pb.Event_JoinedConversation_:
			break

		case *pb.Event_LeftConversation_:
			break

		case *pb.Event_ConversationTitleChanged_:
			break

		default:
			return nil, fmt.Errorf("unknown kind of event '%q'", event.Content)
		}
	}

	newEvents, err := EventsToProto(components, in)

	out = append(out, newEvents...)
	events := append(catchUp, in...)
	sort.Slice(events, func(i, j int) bool {
		return events[i].ID > events[j].ID
	})
	uid := ""
	if len(events) > 0 {
		uid = events[len(events)-1].UUID.String()
	}

	return &pb.DownSync{
		Events:        out,
		LastEventUuid: uid,
	}, nil
}

var _ pb.EventServiceServer = EventService{}

func EventsToProto(components *common.Components, events models.EventSlice) (out []*pb.Event, err error) {
	// EAGER LOADS
	messages, err := components.MessageRepository.ByIds(events.MessageIDs()...)
	if err != nil {
		return nil, err
	}
	pbMessages, err := components.MessagesToProto(messages)
	if err != nil {
		return nil, err
	}
	uuid2pbMessage := pbMessages.UuidMap()
	mid2uuid := messages.IntToUuidMap()

	convs, err := components.ConversationRepository.ByIds(append(events.ConversationIDs(), messages.ConversationIDs()...)...)
	if err != nil {
		return nil, err
	}
	pbConvs, err := components.ConversationsToProto(convs)
	if err != nil {
		return nil, err
	}
	uuid2pbConv := pbConvs.UuidMap()

	for _, event := range events {
		var content pb.EventContentOneOf

		switch event.Type {
		case models.EventTypeNewMessage:
			muid := mid2uuid[event.MessageID.Int]
			msg := uuid2pbMessage[muid]
			cuid, _ := uuid2.FromString(msg.ConvUuid)
			content = &pb.Event_ReceivedNewMessage_{ReceivedNewMessage: &pb.Event_ReceivedNewMessage{
				Message:      msg,
				Conversation: uuid2pbConv[cuid],
			}}

		case models.EventTypeDeletedMessage:
			content = &pb.Event_DeletedMessage_{DeletedMessage: &pb.Event_DeletedMessage{Uuid: event.DeletedMessageUUID.String}}

		case models.EventTypeChangedPicture:
			break

		case models.EventTypeJoinedConversation:
			break

		case models.EventTypeLeftConversation:
			break

		case models.EventTypeConversationTitleChanged:
			break

		default:
			return nil, status.Errorf(codes.Internal, "unknown event type: %q", event.Type)
		}

		out = append(out, &pb.Event{
			Uuid:    event.UUID.String(),
			Content: content,
		})
	}
	return out, nil
}

func OnNewMessage(components *common.Components, me *models.User, nm *pb.Event_SentNewMessage) (out []*models.Event, err error) {
	var conv *models.Conversation
	switch nm.Conversation.(type) {
	case *pb.Event_SentNewMessage_ConvUuid:
		uid, err := uuid2.FromString(nm.GetConvUuid())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "bad conversation uuid: %+v", err)
		}

		conv, err = components.ConversationRepository.ByUuid(uid)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not find conversation: %+v", err)
		}

	case *pb.Event_SentNewMessage_NewConversation:
		allUuids := append(nm.GetNewConversation().UserUuids, me.UUID.String())
		var uuids slices.UuidSlice
		for _, uidStr := range allUuids {
			uid, err := uuid2.FromString(uidStr)
			if err != nil {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}

			uuids = append(uuids, uid)
		}
		uuids = uuids.Unique()

		conv, err = getOrCreateConvForUsers(components, uuids, nm.GetNewConversation().GetTitle())
		if err != nil {
			return nil, err
		}

	default:
		return nil, status.Errorf(codes.Internal, "received unknown conversation input: %q", nm.Conversation)
	}

	if ok, err := components.ConversationHasAccess(me, conv); !ok || err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "can't send message to this conversation: (err = %+v)", err)
	}

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
			return nil, status.Error(codes.InvalidArgument, "could not serialize transcript")
		}

		blobUuid, err := components.StorageClient.Upload(components.Ctx, bytes.NewReader(vm.RawContent))
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
		if err = asset.Insert(components.Ctx, components.Db, boil.Infer()); err != nil {
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

	if err := msg.Insert(components.Ctx, components.Db, boil.Infer()); err != nil {
		return nil, status.Errorf(codes.Internal, "could not insert message: %+v", err)
	}

	//
	// -- STORE EVENTS IN DB FOR ALL USERS
	//
	ucs, err := components.ConversationUsers(conv)
	if err != nil {
		return nil, err
	}
	participants, err := components.UserRepository.FromUserConversations([][]*models.UserConversation{ucs})
	q := sq.Insert(models.TableNames.Event).Columns(
		models.EventColumns.Type,
		models.EventColumns.RecipientID,
		models.EventColumns.MessageID, models.EventColumns.ConversationID,
	)
	for _, p := range participants {
		q = q.Values(models.EventTypeNewMessage, p.ID, nil, conv.ID)
	}
	query, args, _ := q.Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	var dbEvs []*models.Event
	if err = queries.Raw(query, args...).Bind(components.Ctx, components.Db, &dbEvs); err != nil {
		return nil, err
	}

	//
	// -- SEND EVENTS TO CONNECTED USERS
	//
	for _, uc := range ucs {
		if uc.UserID == me.ID {
			// TODO: Send an event to acknowledge that the message was effectively sent
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

	return dbEvs, nil
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
	if _, err = components.Db.Exec(query, args); err != nil {
		return nil, status.Errorf(codes.Internal, "could not insert new members to conv: %+v", err)
	}

	return conv, nil
}
