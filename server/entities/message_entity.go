package entities

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Message struct {
	Record     *models.Message
	Components *common.Components
}

func (message Message) ToPb() (*pb.Message, error) {
	var content pb.MessageContentOneOf
	switch message.Record.Type {
	case models.MessageTypeText:
		content = &pb.Message_TextMessage{TextMessage: &pb.TextMessage{Content: message.Record.Text.String}}
	case models.MessageTypeVoice:
		rawAudios, err := message.Components.AssetStore.ByIds(message.Record.RawAudioID.Int)
		if err != nil || len(rawAudios) == 0 {
			return nil, fmt.Errorf("could not fetch associated audio asset: %+v", err)
		}

		buf := &bytes.Buffer{}
		if err := message.Components.Storage.Download(rawAudios[0].Record.BlobName.String, buf); err != nil {
			return nil, fmt.Errorf("could not download audio file from bucket: %+v", err)
		}

		var pbTranscript pb.AlignedTranscript
		if err := proto.Unmarshal(message.Record.SiriTranscript.Bytes, &pbTranscript); err != nil {
			return nil, fmt.Errorf("could not build protobuf from stored bytearray: %+v", err)
		}
		content = &pb.Message_VoiceMessage{VoiceMessage: &pb.VoiceMessage{RawContent: buf.Bytes(), SiriTranscript: &pbTranscript}}
	default:
		// TODO
		return nil, errors.New("message content type other than text or voice unhandled atm.")
	}

	convs, err := message.Components.ConversationStore.ByIds(message.Record.ConversationID)
	if err != nil {
		return nil, err
	}

	var author *pb.User
	if message.Record.AuthorID.Valid {
		authors, err := message.Components.UserStore.ByIds(message.Record.AuthorID.Int)
		if err != nil {
			return nil, err
		}
		author = authors[0].ToPb()
	}

	return &pb.Message{
		Uuid:      message.Record.UUID.String(),
		ConvUuid:  convs[0].Record.UUID.String(),
		Content:   content,
		Author:    author,
		CreatedAt: timestamppb.New(message.Record.CreatedAt),
	}, nil
}

func (ms MessageSlicePtrs) LoadRawAudio() ([]*Asset, error) {
	ids := []int{}
	for _, m := range ms {
		if m.Record.RawAudioID.Valid {
			ids = append(ids, m.Record.RawAudioID.Int)
		}
	}

	if len(ids) > 0 {
		return ms[0].Components.AssetStore.ByIds(ids...)
	}

	return nil, nil
}

func (ms MessageSlicePtrs) LoadAuthors() ([]*User, error) {
	ids := []int{}
	for _, m := range ms {
		if m.Record.AuthorID.Valid {
			ids = append(ids, m.Record.AuthorID.Int)
		}
	}

	if len(ids) > 0 {
		return ms[0].Components.UserStore.ByIds(ids...)
	}

	return nil, nil
}

func (ms MessageSlicePtrs) LoadConversations() ([]*Conversation, error) {
	ids := []int{}
	for _, m := range ms {
		ids = append(ids, m.Record.ConversationID)
	}

	return ms[0].Components.ConversationStore.ByIds(ids...)
}
