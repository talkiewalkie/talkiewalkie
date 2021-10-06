package coco

import (
	"context"
	"encoding/json"
	"fmt"
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
	"google.golang.org/protobuf/types/known/timestamppb"
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

		case "newmessage":
			var msg common.NewMessageEvent
			if err = json.Unmarshal([]byte(psEvent.Extra), &msg); err != nil {
				log.Printf("failed to unmarshal pubsub event: %+v", err)
				break
			}

			err = server.Send(&pb.Message{
				ConvUuid:   msg.ConversationUuid,
				Content:    &pb.Message_TextMessage{TextMessage: &pb.TextMessage{Content: msg.Text}},
				AuthorUuid: msg.AuthorUuid,
				CreatedAt:  timestamppb.New(msg.Timestamp),
			})
			if err != nil {
				log.Printf("failed to send message in server stream: %+v", err)
				break
			}

		default:
			log.Printf("received unknown pubsub message on topic '%s': (%T) %+v", topic, payload, payload)

		}
	}
}

func (ms MessageService) Send(ctx context.Context, input *pb.MessageSendInput) (*pb.Empty, error) {

	u, err := common.GetUser(ctx)
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

		conv, err = models.Conversations(models.ConversationWhere.UUID.EQ(uid)).One(ctx, ms.Db)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

	case *pb.MessageSendInput_Handles:
		rawHandles := append(input.GetHandles().Handles, u.Handle)
		handles := make(map[string]int)
		for _, handle := range rawHandles {
			handles[handle] += 1
		}

		uniqueHandles := []string{}
		for handle, _ := range handles {
			uniqueHandles = append(uniqueHandles, handle)
		}

		recipients, err := models.Users(models.UserWhere.Handle.IN(uniqueHandles)).All(ctx, ms.Db)
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("could not find recipients: %+v", err))
		}
		if len(recipients) != len(uniqueHandles) {
			return nil, status.Error(codes.Internal, "some users where not found")
		}

		ids := []int{u.ID}
		for _, recipient := range recipients {
			if recipient.ID != u.ID {
				ids = append(ids, recipient.ID)
			}
		}
		// TODO: crash if len(recipients) != handles

		ugs, err := models.UserConversations(
			models.UserConversationWhere.UserID.EQ(u.ID),
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
			newConversation := models.Conversation{}
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

	var text string
	switch input.Content.(type) {
	case *pb.MessageSendInput_TextMessage:
		text = input.GetTextMessage().Content
	default:
		return nil, status.Error(codes.Internal, "unknown content type!")
	}

	msg := &models.Message{
		Text:           text,
		AuthorID:       null.IntFrom(u.ID),
		ConversationID: conv.ID,
		CreatedAt:      time.Now(),
	}
	if err = msg.Insert(ctx, ms.Db, boil.Infer()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Empty{}, nil
}