package entities

import (
	"bytes"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/golang/protobuf/proto"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MessageToPb(m *models.Message, c *common.Components) (*pb.Message, error) {
	if m.R.Conversation == nil {
		return nil, errors.New("message record has not loaded its conversation")
	}
	if m.R.Author == nil {
		return nil, errors.New("message record has not loaded its author")
	}
	if m.RawAudioID.Valid && m.R.RawAudio == nil {
		return nil, errors.New("message record has not loaded its audio")
	}

	var content pb.MessageContentOneOf
	switch m.Type {
	case models.MessageTypeText:
		content = &pb.Message_TextMessage{TextMessage: &pb.TextMessage{Content: m.Text.String}}
	case models.MessageTypeVoice:
		buf := &bytes.Buffer{}
		if err := c.Storage.Download(m.R.RawAudio.BlobName.String, buf); err != nil {
			return nil, fmt.Errorf("could not download audio file from bucket: %+v", err)
		}

		var pbTranscript pb.AlignedTranscript
		if err := proto.Unmarshal(m.SiriTranscript.Bytes, &pbTranscript); err != nil {
			return nil, fmt.Errorf("could not build protobuf from stored bytearray: %+v", err)
		}
		content = &pb.Message_VoiceMessage{VoiceMessage: &pb.VoiceMessage{RawContent: buf.Bytes(), SiriTranscript: &pbTranscript}}
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
