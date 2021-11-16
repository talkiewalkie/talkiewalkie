package api

import (
	"context"
	uuid2 "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"github.com/talkiewalkie/talkiewalkie/testutils"
	"log"
	"testing"
)

func TestNewMessage(t *testing.T) {
	db := testutils.SetupDb()
	ctx := context.Background()
	service := EventService{}

	testutils.TearDownDb(ctx, db)
	t.Run("first call, no last event", func(t *testing.T) {
		_, _, ctx := testutils.NewContext(db, t)
		out, err := service.Sync(ctx, &pb.UpSync{Events: nil, LastEventUuid: ""})
		require.NoError(t, err)
		require.NotNil(t, out)
		require.Equal(t, 0, len(out.Events), "no events to load")
		require.Equal(t, "", out.LastEventUuid, "no events to load")
	})

	testutils.TearDownDb(ctx, db)
	t.Run("send new message", func(t *testing.T) {
		_, _, ctx := testutils.NewContext(db, t)
		out, err := service.Sync(ctx, &pb.UpSync{Events: []*pb.Event{
			{
				LocalUuid: uuid2.NewV4().String(),
				Content: &pb.Event_SentNewMessage_{SentNewMessage: &pb.Event_SentNewMessage{
					Message:      &pb.MessageSendInput{Content: &pb.MessageSendInput_TextMessage{TextMessage: &pb.TextMessage{Content: "hello"}}},
					Conversation: &pb.Event_SentNewMessage_NewConversation{NewConversation: &pb.ConversationInput{Title: "new", UserUuids: []string{}}},
				}},
			},
		},
		})
		require.NoError(t, err)
		require.NotNil(t, out)
		log.Println(out.Events)
		log.Println(out.LastEventUuid)
	})
}
