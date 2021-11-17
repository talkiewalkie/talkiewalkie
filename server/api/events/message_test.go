package events

import (
	"context"
	"testing"
	"time"

	uuid2 "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/talkiewalkie/talkiewalkie/pb"
	"github.com/talkiewalkie/talkiewalkie/repositories"
	"github.com/talkiewalkie/talkiewalkie/testutils"
)

func TestNewMessage(t *testing.T) {
	db := testutils.SetupDb()

	testutils.TearDownDb(db)
	t.Run("send new message notifies others and self", func(t *testing.T) {
		components, me, ctx := testutils.NewContext(db, t)

		u1 := testutils.AddMockUser(db, t)
		u2 := testutils.AddMockUser(db, t)
		conv := testutils.AddMockConversation(db, t, me, u1, u2)

		u1topic := repositories.UserPubSubTopic(u1)
		u2topic := repositories.UserPubSubTopic(u2)
		metopic := repositories.UserPubSubTopic(me)
		u1mq, u1cancel, err := components.PubSubClient.Subscribe(u1topic)
		u2mq, u2cancel, err := components.PubSubClient.Subscribe(u2topic)
		memq, mecancel, err := components.PubSubClient.Subscribe(metopic)
		defer u1cancel()
		defer u2cancel()
		defer mecancel()

		localUuid := uuid2.NewV4()
		event := &pb.Event{
			LocalUuid: localUuid.String(),
			Content: &pb.Event_SentNewMessage_{SentNewMessage: &pb.Event_SentNewMessage{
				Message:      &pb.MessageSendInput{Content: &pb.MessageSendInput_TextMessage{TextMessage: &pb.TextMessage{Content: "hello"}}},
				Conversation: &pb.Event_SentNewMessage_ConvUuid{ConvUuid: conv.UUID.String()}},
			}}
		pbE, _, err := OnNewMessage(components, me, event)
		if err != nil {
			t.Fatal(err)
		}
		require.Equal(t, localUuid.String(), pbE.LocalUuid)

		deadline, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()

		select {
		case m := <-u1mq:
			event := &pb.Event{}
			if err := protojson.Unmarshal([]byte(m.Extra), event); err != nil {
				t.Fatal(err)
			}
			require.IsType(t, &pb.Event_ReceivedNewMessage_{}, event.Content)

		case m := <-u2mq:
			event := &pb.Event{}
			if err := protojson.Unmarshal([]byte(m.Extra), event); err != nil {
				t.Fatal(err)
			}
			require.IsType(t, &pb.Event_ReceivedNewMessage_{}, event.Content)

		case m := <-memq:
			event := &pb.Event{}
			if err := protojson.Unmarshal([]byte(m.Extra), event); err != nil {
				t.Fatal(err)
			}
			require.Equal(t, localUuid.String(), event.LocalUuid)
			require.IsType(t, &pb.Event_ReceivedNewMessage_{}, event.Content)

		case <-deadline.Done():
			t.Log("did not receive messages under the deadline")
			cancel()
			t.Fail()
		}

	})
}
