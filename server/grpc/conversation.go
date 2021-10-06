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
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
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
		qm.Load(models.ConversationRels.Messages), qm.Limit(50), qm.OrderBy(fmt.Sprintf("%s DESC", models.ConversationColumns.CreatedAt)),
		qm.Load(qm.Rels(models.ConversationRels.Messages, models.MessageRels.Author)),
	).One(ctx, c.Db)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("could not get conversation: %+v", err))
	}

	msgs := []*pb.Message{}
	for _, message := range conv.R.Messages {
		msgs = append(msgs, &pb.Message{
			ConvUuid:   conv.UUID.String(),
			Content:    &pb.Message_TextMessage{TextMessage: &pb.TextMessage{Content: message.Text}},
			AuthorUuid: message.R.Author.UUID.String(),
			CreatedAt:  timestamppb.New(message.CreatedAt),
		})
	}

	return &pb.Conversation{
		Uuid:     conv.UUID.String(),
		Title:    entities.ConversationDisplay(*conv),
		Messages: msgs,
	}, nil
}

type convListJoin1 struct {
	models.Conversation     `boil:",bind"`
	models.UserConversation `boil:",bind"`
}

type convListJoin2 struct {
	models.User             `boil:",bind"`
	models.UserConversation `boil:",bind"`
}

func (c ConversationService) List(input *pb.ConversationListInput, server pb.ConversationService_ListServer) error {
	u, err := common.GetUser(server.Context())
	if err != nil {
		return status.Error(codes.PermissionDenied, err.Error())
	}

	var convs []*convListJoin1
	err = models.NewQuery(
		qm.Select("*"),
		qm.From(models.TableNames.Conversation),
		qm.InnerJoin(fmt.Sprintf("%s on %s = %s", models.TableNames.UserConversation, models.ConversationTableColumns.ID, models.UserConversationTableColumns.ConversationID)),
		models.UserConversationWhere.UserID.EQ(u.ID),
		qm.Limit(20), qm.Offset(int(input.Page)),
	).Bind(server.Context(), c.Db, &convs)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	convIds := []int{}
	for _, joinRes := range convs {
		convIds = append(convIds, joinRes.Conversation.ID)
	}

	var ucs []*convListJoin2
	err = models.NewQuery(
		qm.Select("*"),
		qm.From(models.TableNames.UserConversation),
		qm.InnerJoin(fmt.Sprintf("\"%s\" on \"%s\".\"%s\" = %s", models.TableNames.User, models.TableNames.User, models.UserColumns.ID, models.UserConversationTableColumns.UserID)),
		models.UserConversationWhere.ConversationID.IN(convIds),
	).Bind(server.Context(), c.Db, &ucs)

	cid2users := map[int][]*models.User{}
	for _, uc := range ucs {
		cid2users[uc.UserConversation.ConversationID] = append(cid2users[uc.UserConversation.ConversationID], &uc.User)
	}

	for _, joinRes := range convs {
		conv := joinRes.Conversation
		title := conv.Name.String
		if !conv.Name.Valid {
			participants, ok := cid2users[conv.ID]
			if !ok {
				return status.Error(codes.Internal, "conversation participants not found")
			}
			handles := []string{}
			for _, participant := range participants {
				handles = append(handles, participant.Handle)
			}
			title = strings.Join(handles, ", ")
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
