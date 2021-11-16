package api

import (
	"context"
	"fmt"
	uuid2 "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/api/events"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sort"

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
	if sync.LastEventUuid != "" {
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

		out, err = events.EventsToProto(components, catchUp)
		if err != nil {
			return nil, err
		}
	}

	var pbNewEvents events.EventSlice
	var dbNewEvents models.EventSlice
	for _, event := range sync.Events {
		switch event.Content.(type) {
		case *pb.Event_SentNewMessage_:
			pbNewEvent, dbNewEvent, err := events.OnNewMessage(components, me, event)
			if err != nil {
				return nil, err
			}

			pbNewEvents = append(pbNewEvents, pbNewEvent)
			dbNewEvents = append(dbNewEvents, dbNewEvent)

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

	//
	// -- FILTER TO ONLY OUTPUT RELEVANT EVENTS
	//
	out = append(out, pbNewEvents...)

	//
	// -- LAST UUID COMPUTE
	//
	var uid string
	events := append(catchUp, dbNewEvents...)
	sort.Slice(events, func(i, j int) bool { return events[i].ID > events[j].ID })
	if len(events) > 0 {
		uid = events[len(events)-1].UUID.String()
	}

	return &pb.DownSync{Events: out, LastEventUuid: uid}, nil
}

var _ pb.EventServiceServer = EventService{}
