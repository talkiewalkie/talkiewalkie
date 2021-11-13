package api

import (
	"context"
	"fmt"
	uuid2 "github.com/satori/go.uuid"
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

	conv, err := components.ConversationRepository.ByUuid(uid)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("could not get conversation: %+v", err))
	}

	if ok, err := components.ConversationHasAccess(me, conv); !ok || err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "cannot access this conversation: %+v", err)
	}

	messages, err := models.Messages(
		models.MessageWhere.ConversationID.EQ(conv.ID),
		qm.Limit(20),
		qm.OrderBy(fmt.Sprintf("%s DESC", models.MessageColumns.CreatedAt)),
	).All(ctx, components.Db)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get conv's messages: %+v", err)
	}

	mPbs, err := components.MessagesToProto(messages)
	if err != nil {
		return nil, err
	}

	pbConvs, err := components.ConversationsToProto([]*models.Conversation{conv})
	if err != nil {
		return nil, err
	}

	pbConv := pbConvs[0]
	pbConv.Messages = mPbs

	return pbConv, nil
}

// This is an expensive call - we're not paginating.
func (c ConversationService) List(input *pb.ConversationListInput, server pb.ConversationService_ListServer) error {
	components, me, err := WithAuthedContext(server.Context())
	if err != nil {
		return err
	}

	myConversations, err := components.UserConversationRepository.ByUserIds(me.ID)
	if err != nil {
		return err
	}

	if len(myConversations) == 0 {
		return nil
	}

	// Eager load attached items
	convs, err := components.ConversationRepository.FromUserConversations(myConversations)
	if err != nil {
		return err
	}
	if _, err := components.UserRepository.FromUserConversations(myConversations); err != nil {
		return err
	}

	// Fetch last messages for display
	lastMessages, err := components.MessageRepository.FromConversationsLast(convs)
	pbMsgs, err := components.MessagesToProto(lastMessages)
	if err != nil {
		return err
	}

	pbConvs, err := components.ConversationsToProto(convs)
	if err != nil {
		return err
	}

	for idx, pbConv := range pbConvs {
		pbConv.Messages = []*pb.Message{pbMsgs[idx]}
		if err = server.Send(pbConv); err != nil {
			return err
		}
	}

	return nil
}
