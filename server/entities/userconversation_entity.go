package entities

import (
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserConversation struct {
	Record     *models.UserConversation
	Components *common.Components
}

func (uc UserConversation) User() (*User, error) {
	return uc.Components.UserStore.ById(uc.Record.UserID)
}

func (uc UserConversation) Conversation() (*Conversation, error) {
	return uc.Components.ConversationStore.ById(uc.Record.ConversationID)
}

func (uc UserConversation) ToPb() (*pb.UserConversation, error) {
	user, err := uc.User()
	if err != nil {
		return nil, err
	}
	return &pb.UserConversation{User: user.ToPb(), ReadUntil: timestamppb.New(uc.Record.ReadUntil)}, nil
}
