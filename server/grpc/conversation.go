package coco

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	uuid2 "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/entities"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"strconv"
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
	me, err := common.GetUser(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	uid, err := uuid2.FromString(input.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("bad uuid: %+v", err))
	}

	conv, err := models.Conversations(
		models.ConversationWhere.UUID.EQ(uid),
		qm.Load(qm.Rels(models.ConversationRels.UserConversations, models.UserConversationRels.User)),
	).One(ctx, c.Db)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("could not get conversation: %+v", err))
	}

	if ok, err := entities.CanAccessConversation(conv, me); !ok || err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "cannot access this conversation: %+v", err)
	}

	messages, err := models.Messages(
		models.MessageWhere.ConversationID.EQ(conv.ID),
		qm.Load(models.MessageRels.RawAudio),
		qm.Load(models.MessageRels.Author),
		qm.Limit(20),
		qm.OrderBy(fmt.Sprintf("%s DESC", models.MessageColumns.CreatedAt)),
	).All(ctx, c.Db)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get conv's messages: %+v", err)
	}

	msgs := []*pb.Message{}
	for _, message := range messages {
		message.R.Conversation = conv
		pbm, err := entities.MessageToPb(message, c.Components)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		msgs = append(msgs, pbm)
	}

	title, err := entities.ConversationDisplay(conv, me)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not compute conversation title: %+v", err)
	}

	return &pb.Conversation{
		Uuid:         conv.UUID.String(),
		Title:        title,
		Messages:     msgs,
		Participants: nil, // TODO
	}, nil
}

func (c ConversationService) List(input *pb.ConversationListInput, server pb.ConversationService_ListServer) error {
	me, err := common.GetUser(server.Context())
	if err != nil {
		return status.Error(codes.PermissionDenied, err.Error())
	}

	myConvs, err := models.UserConversations(
		models.UserConversationWhere.UserID.EQ(me.ID),
		qm.Load(qm.Rels(models.UserConversationRels.Conversation, models.ConversationRels.UserConversations, models.UserConversationRels.User)),
		qm.Limit(20), qm.Offset(int(input.Page)),
	).All(server.Context(), c.Db)

	convIds := []string{}
	for _, conv := range myConvs {
		convIds = append(convIds, strconv.Itoa(conv.ConversationID))
	}

	if len(convIds) == 0 {
		return nil
	}

	type lastMessage struct {
		Discard1       int `boil:"discard_1"`
		Discard2       int `boil:"discard_2"`
		models.Message `boil:",bind"`
	}
	var mssgs []*lastMessage

	sqlStr := `
SELECT 
	DISTINCT ON (conversation_id) conversation_id discard_1, 
	first_value("id") OVER (PARTITION BY "conversation_id" ORDER BY created_at) "discard_2", 
	*
FROM "message"
WHERE "conversation_id" in (%s);
`
	if err = models.NewQuery(qm.SQL(fmt.Sprintf(sqlStr, strings.Join(convIds, ",")))).Bind(server.Context(), c.Db, &mssgs); err != nil {
		return err
	}

	for _, uc := range myConvs {
		conv := uc.R.Conversation
		title, err := entities.ConversationDisplay(conv, me)
		if err != nil {
			return status.Errorf(codes.Internal, "could not display conversation title: %+v", err)
		}

		participants := []*pb.User{}
		id2p := map[int]*models.User{}
		for _, uc := range conv.R.UserConversations {
			user := uc.R.User
			participants = append(participants, &pb.User{
				DisplayName: entities.UserDisplayName(user),
				Uuid:        user.UUID.String(),
				Phone:       user.PhoneNumber,
			})
			id2p[user.ID] = user
		}

		messages := []*pb.Message{}
		for _, m := range mssgs {
			if m.ConversationID == conv.ID {
				var content pb.MessageContentOneOf
				switch m.Type {
				case models.MessageTypeText:
					content = &pb.Message_TextMessage{TextMessage: &pb.TextMessage{Content: m.Text.String}}
				case models.MessageTypeVoice:
					var pbTranscript *pb.AlignedTranscript
					if err := proto.Unmarshal(m.SiriTranscript.Bytes, pbTranscript); err != nil {
						return status.Errorf(codes.Internal, "could not build protobuf from stored bytearray: %+v", err)
					}
					content = &pb.Message_VoiceMessage{VoiceMessage: &pb.VoiceMessage{SiriTranscript: pbTranscript}}
				}
				var author *pb.User
				if m.AuthorID.Valid {
					if tt, ok := id2p[m.AuthorID.Int]; ok {
						author = &pb.User{
							DisplayName: entities.UserDisplayName(tt),
							Uuid:        tt.UUID.String(),
							Phone:       tt.PhoneNumber,
						}
					}
				}
				messages = append(messages, &pb.Message{
					Uuid:      m.UUID.String(),
					ConvUuid:  conv.UUID.String(),
					Content:   content,
					Author:    author,
					CreatedAt: timestamppb.New(m.CreatedAt),
				})
			}
		}

		log.Printf("%d msgs found for conv %s", len(messages), conv.UUID.String())
		if err = server.Send(&pb.Conversation{
			Uuid:         conv.UUID.String(),
			Title:        title,
			Messages:     messages,
			Participants: participants,
		}); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}

	return nil
}
