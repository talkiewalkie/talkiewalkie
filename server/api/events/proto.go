package events

import (
	uuid2 "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EventSlice []*pb.Event

func (slice EventSlice) UuidMap() map[string]*pb.Event {
	out := make(map[string]*pb.Event, len(slice))
	for _, event := range slice {
		out[event.Uuid] = event
	}

	return out
}

func (slice EventSlice) LocalUuidMap() map[string]*pb.Event {
	out := make(map[string]*pb.Event, len(slice))
	for _, event := range slice {
		if event.LocalUuid != "" {
			out[event.LocalUuid] = event
		}
	}

	return out
}

func EventsToProto(components *common.Components, events models.EventSlice) (out EventSlice, err error) {
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
