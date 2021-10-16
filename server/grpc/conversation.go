package coco

import (
	"context"
	"fmt"
	uuid2 "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/entities"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ConversationService struct {
	*common.Components
}

var _ pb.ConversationServiceServer = ConversationService{}

func NewConversationService(c *common.Components) ConversationService {
	return ConversationService{Components: c}
}

func (c ConversationService) Get(ctx context.Context, input *pb.ConversationGetInput) (*pb.Conversation, error) {
	uid, err := uuid2.FromString(input.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("bad uuid: %+v", err))
	}

	conv, err := models.Conversations(
		models.ConversationWhere.UUID.EQ(uid),
		qm.Load(
			qm.Rels(models.ConversationRels.Messages, models.MessageRels.Author),
			qm.Limit(50), qm.OrderBy(fmt.Sprintf("%s DESC", models.ConversationColumns.CreatedAt))),
	).One(ctx, c.Db)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("could not get conversation: %+v", err))
	}

	msgs := []*pb.Message{}
	for _, message := range conv.R.Messages {
		pbm, err := entities.MessageToPb(message)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		msgs = append(msgs, pbm)
	}

	return &pb.Conversation{
		Uuid:     conv.UUID.String(),
		Title:    entities.ConversationDisplay(conv),
		Messages: msgs,
	}, nil
}

func (c ConversationService) List(input *pb.ConversationListInput, server pb.ConversationService_ListServer) error {
	u, err := common.GetUser(server.Context())
	if err != nil {
		return status.Error(codes.PermissionDenied, err.Error())
	}

	myConvs, err := models.UserConversations(
		models.UserConversationWhere.UserID.EQ(u.ID),
		qm.Load(qm.Rels(models.UserConversationRels.Conversation, models.ConversationRels.UserConversations, models.UserConversationRels.User)),
		qm.Load(qm.Rels(models.UserConversationRels.Conversation, models.ConversationRels.Messages, models.MessageRels.Author), qm.Limit(10), qm.OrderBy(fmt.Sprintf("%s DESC", models.MessageColumns.CreatedAt))),
		qm.Limit(20), qm.Offset(int(input.Page)),
	).All(server.Context(), c.Db)

	for _, uc := range myConvs {
		conv := uc.R.Conversation
		title := entities.ConversationDisplay(conv)
		messages := []*pb.Message{}
		for _, m := range conv.R.Messages {
			pbm, err := entities.MessageToPb(m)
			if err != nil {
				return err
			}
			messages = append(messages, pbm)
		}

		if err = server.Send(&pb.Conversation{
			Uuid:  conv.UUID.String(),
			Title: title,
		}); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}

	return nil
}
