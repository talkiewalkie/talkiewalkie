package entities

import (
	"github.com/friendsofgo/errors"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MessageToPb(m *models.Message) (*pb.Message, error) {
	if m.R.Conversation == nil {
		return nil, errors.New("message record has not loaded its conversation")
	}
	if m.R.Author == nil {
		return nil, errors.New("message record has not loaded its author")
	}

	var content pb.MessageContentOneOf
	switch m.Type {
	case models.MessageTypeText:
		content = &pb.Message_TextMessage{TextMessage: &pb.TextMessage{Content: m.Text.String}}
	case models.MessageTypeVoice:
		content = &pb.Message_VoiceMessage{VoiceMessage: &pb.VoiceMessage{Url: "todo"}}
	default:
		// TODO
		return nil, errors.New("message content type other than text or voice unhandled atm.")
	}

	return &pb.Message{
		Uuid:     m.UUID.String(),
		ConvUuid: m.R.Conversation.UUID.String(),
		Content:  content,
		Author: &pb.User{
			Uuid:        m.R.Author.UUID.String(),
			DisplayName: UserDisplayName(m.R.Author),
			Phone:       m.R.Author.PhoneNumber,
		},
		CreatedAt: timestamppb.New(m.CreatedAt),
	}, nil
}
