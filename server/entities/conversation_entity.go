package entities

import (
	"context"
	"fmt"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"strings"
)

type Conversation struct {
	Record     *models.Conversation
	Components *common.Components
}

func (conv Conversation) Participants() ([]*UserConversation, error) {
	ucs, err := conv.Components.UserConversationStore.ByConversationIds(conv.Record.ID)
	if err != nil {
		return nil, err
	}

	// Eager load user records by default
	userIds := []int{}
	for _, uc := range ucs[0] {
		userIds = append(userIds, uc.Record.UserID)
	}

	_, err = conv.Components.UserStore.ByIds(userIds...)
	if err != nil {
		return nil, err
	}

	return ucs[0], nil
}

func (conv Conversation) HasAccess(u *User) (bool, error) {
	ucs, err := conv.Components.UserConversationStore.ByUserIds(u.Record.ID)
	if err != nil {
		return false, err
	}

	for _, uc := range ucs[0] {
		if uc.Record.ConversationID == conv.Record.ID {
			return true, nil
		}
	}

	return false, nil
}

func (conv Conversation) ToPb(withMessages MessageSlicePtrs) (*pb.Conversation, error) {
	pbMessages := []*pb.Message{}
	for _, message := range withMessages {
		pbm, err := message.ToPb()
		if err != nil {
			return nil, err
		}
		pbMessages = append(pbMessages, pbm)
	}

	participants, err := conv.Participants()
	if err != nil {
		return nil, err
	}

	pbUsers := []*pb.UserConversation{}
	for _, user := range participants {
		pbu, err := user.ToPb()
		if err != nil {
			return nil, err
		}
		pbUsers = append(pbUsers, pbu)
	}

	return &pb.Conversation{
		Uuid:         conv.Record.UUID.String(),
		Title:        conv.Record.Name.String,
		Messages:     pbMessages,
		Participants: pbUsers,
	}, nil
}

// Slice utils

func (cs ConversationSlice) LoadParticipants() ([]*User, error) {
	if len(cs) == 0 {
		return nil, nil
	}

	convIds := []int{}
	for _, conv := range cs {
		convIds = append(convIds, conv.Record.ID)
	}

	convsRot, err := cs[0].Components.UserConversationStore.ByConversationIds(convIds...)
	if err != nil {
		return nil, err
	}

	userIdsSet := map[int]bool{}
	for _, convPats := range convsRot {
		for _, uc := range convPats {
			userIdsSet[uc.Record.UserID] = true
		}
	}

	userIds := []int{}
	for id, _ := range userIdsSet {
		userIds = append(userIds, id)
	}

	users, err := cs[0].Components.UserStore.ByIds(userIds...)
	return users, err
}

func (cs ConversationSlicePtrs) LoadLastMessages(ctx context.Context) (MessageSlicePtrs, error) {
	if len(cs) == 0 {
		return nil, nil
	}
	components := cs[0].Components

	type lastMessage struct {
		Discard1       int `boil:"discard_1"`
		Discard2       int `boil:"discard_2"`
		models.Message `boil:",bind"`
	}
	var mssgs []*lastMessage

	sqlStr := `
		SELECT 
			DISTINCT ON (conversation_id) conversation_id discard_1, 
			last_value("id") OVER (PARTITION BY "conversation_id" ORDER BY created_at DESC) "discard_2", 
			*
		FROM "message"
		WHERE "conversation_id" in (%s);
	`
	convIds := cs.MapToString(func(c *Conversation) string { return strconv.Itoa(c.Record.ID) })

	if err := models.NewQuery(qm.SQL(fmt.Sprintf(sqlStr, strings.Join(convIds, ",")))).Bind(ctx, components.Db, &mssgs); err != nil {
		return nil, status.Errorf(codes.Internal, "could not fetch last messages for each conv: %+v", err)
	}

	out := []*Message{}
	for _, item := range mssgs {
		out = append(out, &Message{&item.Message, components})
	}

	return out, nil
}
