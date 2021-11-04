package coco

import (
	"context"
	"fmt"
	uuid2 "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/entities"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/pb"
	_ "github.com/talkiewalkie/talkiewalkie/pkg/slices"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ConversationService struct {
}

var _ pb.ConversationServiceServer = ConversationService{}

func NewConversationService() ConversationService {
	return ConversationService{}
}

func (c ConversationService) Get(ctx context.Context, input *pb.ConversationGetInput) (*pb.Conversation, error) {
	components, me, err := WithAuthedContext(ctx)
	if err != nil {
		return nil, err
	}

	uid, err := uuid2.FromString(input.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("bad uuid: %+v", err))
	}

	conv, err := components.ConversationStore.ByUuid(uid)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("could not get conversation: %+v", err))
	}

	if ok, err := conv.HasAccess(me); !ok || err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "cannot access this conversation: %+v", err)
	}

	messageRecords, err := models.Messages(
		models.MessageWhere.ConversationID.EQ(conv.Record.ID),
		qm.Limit(20),
		qm.OrderBy(fmt.Sprintf("%s DESC", models.MessageColumns.CreatedAt)),
	).All(ctx, components.Db)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get conv's messages: %+v", err)
	}

	var messages entities.MessageSlicePtrs
	for _, record := range messageRecords {
		messages = append(messages, &entities.Message{record, components})
	}

	if _, err := messages.LoadAuthors(); err != nil {
		return nil, err
	}
	if _, err := messages.LoadRawAudio(); err != nil {
		return nil, err
	}

	pbConv, err := conv.ToPb(messages)
	if err != nil {
		return nil, err
	}

	return pbConv, nil
}

// This is an expensive call - we're not paginating.
func (c ConversationService) List(input *pb.ConversationListInput, server pb.ConversationService_ListServer) error {
	components, me, err := WithAuthedContext(server.Context())
	if err != nil {
		return err
	}

	myConversations, err := components.UserConversationStore.ByUserIds(me.Record.ID)
	if err != nil {
		return err
	}

	if len(myConversations) == 0 {
		return nil
	}

	// Eager load attached items
	convs, err := components.UserConversationStore.LoadConversationsFromResult(myConversations)
	if err != nil {
		return err
	}
	if _, err := components.UserConversationStore.LoadUsersFromResult(myConversations); err != nil {
		return err
	}

	// Fetch last messages for display
	lastMessages, err := convs.LoadLastMessages(server.Context())

	for _, conv := range convs {
		pbConv, err := conv.ToPb(lastMessages)
		if err != nil {
			return err
		}

		if err = server.Send(pbConv); err != nil {
			return err
		}
	}

	return nil
}
