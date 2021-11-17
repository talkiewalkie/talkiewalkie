package events

import (
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func HandleNewEvent(
	components *common.Components,
	me *models.User,
	event *pb.Event,
	// whether or not to publish messages on pubsub if user is connected - typically on Sync we set it to false in order
	// to avoid double ingestion.
	// but maybe double ingestion is not a problem since we have the localuuid?
	sendThroughWire bool,
) (dbSlice models.EventSlice, pbslice EventSlice, err error) {
	switch event.Content.(type) {
	case *pb.Event_SentNewMessage_:
		pbNewEvent, dbNewEvent, err := OnNewMessage(components, me, event, sendThroughWire)
		if err != nil {
			return nil, nil, err
		}

		pbslice = append(pbslice, pbNewEvent)
		dbSlice = append(dbSlice, dbNewEvent)

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
		return nil, nil, status.Errorf(codes.Internal, "unknown kind of event '%T'", event.Content)
	}

	return
}

func HandleIncomingEvents(components *common.Components, me *models.User, server pb.EventService_ConnectServer, errChan chan error) {
	for {
		event, err := server.Recv()
		if err != nil {
			errChan <- err
			return
		}

		log.Printf("processing new event: %T", event.Content)
		_, _, err = HandleNewEvent(components, me, event, true)
		if err != nil {
			errChan <- err
			return
		}
	}
}
