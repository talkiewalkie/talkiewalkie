package api

import (
	"context"
	"fmt"
	"sort"

	uuid2 "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/talkiewalkie/talkiewalkie/api/events"
	"github.com/talkiewalkie/talkiewalkie/repositories"

	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/pb"
)

type EventService struct{}

func (e EventService) Connect(server pb.EventService_ConnectServer) error {
	components, me, err := WithAuthedContext(server.Context())
	if err != nil {
		return err
	}

	topic := repositories.UserPubSubTopic(me)
	mq, cancel, err := components.PubSubClient.Subscribe(topic)
	if err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("could not subscribe to pubsub topic: %+v", err))
	}
	defer cancel()

	incomingErrs := make(chan error, 1)
	go events.HandleIncomingEvents(components, me, server, incomingErrs)

	for {
		select {
		case m := <-mq:
			components.ResetEntityStores(server.Context())
			event := &pb.Event{}
			if err := protojson.Unmarshal([]byte(m.Extra), event); err != nil {
				return status.Errorf(codes.Internal, "could not process pubsub message: %+v", err)
			}

			if err := server.Send(event); err != nil {
				return status.Errorf(codes.Internal, "failed to send new event: %+v", err)
			}

		case <-server.Context().Done():
			return nil

		case ierr := <-incomingErrs:
			return ierr
		}
	}
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
		dbEvs, pbEvs, err := events.HandleNewEvent(components, me, event)
		if err != nil {
			return nil, err
		}

		pbNewEvents = append(pbNewEvents, pbEvs...)
		dbNewEvents = append(dbNewEvents, dbEvs...)
	}

	out = append(out, pbNewEvents...)

	//
	// -- LAST UUID COMPUTE
	//
	var uid string
	dbEvents := append(catchUp, dbNewEvents...)
	sort.Slice(dbEvents, func(i, j int) bool { return dbEvents[i].ID > dbEvents[j].ID })
	if len(dbEvents) > 0 {
		uid = dbEvents[len(dbEvents)-1].UUID.String()
	}

	return &pb.DownSync{Events: out, LastEventUuid: uid}, nil
}

var _ pb.EventServiceServer = EventService{}
