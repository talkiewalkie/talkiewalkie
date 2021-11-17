package api

import (
	"testing"

	uuid2 "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"google.golang.org/grpc"

	"github.com/talkiewalkie/talkiewalkie/api/events"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"github.com/talkiewalkie/talkiewalkie/repositories"
	"github.com/talkiewalkie/talkiewalkie/testutils"
)

func TestSync(t *testing.T) {
	db := testutils.SetupDb()
	service := EventService{}

	testutils.TearDownDb(db)
	t.Run("send nothing, no last event", func(t *testing.T) {
		_, _, ctx := testutils.NewContext(db, t)

		out, err := service.Sync(ctx, &pb.UpSync{Events: nil, LastEventUuid: ""})

		require.NoError(t, err)
		require.NotNil(t, out)
		require.Equal(t, 0, len(out.Events), "no events to load")
		require.Equal(t, "", out.LastEventUuid, "no events to load")
	})

	testutils.TearDownDb(db)
	t.Run("send new message to new conversation, no last event", func(t *testing.T) {
		_, _, ctx := testutils.NewContext(db, t)

		localUuid := uuid2.NewV4()
		out, err := service.Sync(ctx, &pb.UpSync{Events: []*pb.Event{
			{
				LocalUuid: localUuid.String(),
				Content: &pb.Event_SentNewMessage_{SentNewMessage: &pb.Event_SentNewMessage{
					Message:      &pb.MessageSendInput{Content: &pb.MessageSendInput_TextMessage{TextMessage: &pb.TextMessage{Content: "hello"}}},
					Conversation: &pb.Event_SentNewMessage_NewConversation{NewConversation: &pb.ConversationInput{Title: "new", UserUuids: []string{}}},
				}},
			},
		},
		})

		require.NoError(t, err)
		require.NotNil(t, out)
		require.Equal(t, 1, len(out.Events))
		require.Equal(t, localUuid.String(), out.Events[0].LocalUuid)
	})

	testutils.TearDownDb(db)
	t.Run("send new message with existing events to catch up", func(t *testing.T) {
		_, me, ctx := testutils.NewContext(db, t)

		u1 := testutils.AddMockUser(db, t)
		u2 := testutils.AddMockUser(db, t)
		conv := testutils.AddMockConversation(db, t, me, u1, u2)
		messages := []*models.Message{
			testutils.AddMockMessage(db, t, &models.Message{AuthorID: null.IntFrom(u1.ID), ConversationID: conv.ID}),
			testutils.AddMockMessage(db, t, &models.Message{AuthorID: null.IntFrom(u1.ID), ConversationID: conv.ID}),
			testutils.AddMockMessage(db, t, &models.Message{AuthorID: null.IntFrom(u2.ID), ConversationID: conv.ID}),
		}
		evs, err := events.BatchInsert(ctx, db, []*models.Event{
			{Type: models.EventTypeNewMessage, MessageID: null.IntFrom(messages[0].ID), RecipientID: me.ID},
			{Type: models.EventTypeNewMessage, MessageID: null.IntFrom(messages[1].ID), RecipientID: me.ID},
			{Type: models.EventTypeNewMessage, MessageID: null.IntFrom(messages[2].ID), RecipientID: me.ID},
		})
		if err != nil {
			t.Fatal(err)
		}

		sent := newMsgEvent(conv.UUID.String())
		out, err := service.Sync(ctx,
			&pb.UpSync{
				Events:        []*pb.Event{sent},
				LastEventUuid: evs[0].UUID.String(),
			})

		require.NoError(t, err)
		require.NotNil(t, out)
		require.Equal(t, 3, len(out.Events))
		require.IsType(t, &pb.Event_ReceivedNewMessage_{}, out.Events[0].Content)
		require.IsType(t, &pb.Event_ReceivedNewMessage_{}, out.Events[1].Content)
		require.IsType(t, &pb.Event_ReceivedNewMessage_{}, out.Events[2].Content)

		require.Equal(t, "", out.Events[0].LocalUuid) // non "local" events
		require.Equal(t, "", out.Events[1].LocalUuid) // non "local" events
		require.Equal(t, sent.LocalUuid, out.Events[2].LocalUuid)
	})
}

func TestConnect(t *testing.T) {
	db := testutils.SetupDb()

	testutils.TearDownDb(db)
	t.Run("connect and receive items", func(t *testing.T) {
		components, me, ctx, conn, timeout, _ := testutils.NewFakeServer(
			db, t,
			func(s *grpc.Server) { pb.RegisterEventServiceServer(s, EventService{}) },
		)

		defer timeout()
		defer conn.Close()

		client := pb.NewEventServiceClient(conn)
		connection, err := client.Connect(ctx)
		if err != nil {
			t.Fatal(err)
		}

		u1 := testutils.AddMockUser(db, t)
		u2 := testutils.AddMockUser(db, t)
		conv := testutils.AddMockConversation(db, t, me, u1, u2)

		sent := newMsgEvent(conv.UUID.String())
		if err = connection.Send(sent); err != nil {
			t.Fatal(err)
		}

		event, err := connection.Recv()
		if err != nil {
			t.Fatal(err)
		}
		require.Equal(t, sent.LocalUuid, event.LocalUuid)

		components.ResetEntityStores(ctx)
		msgs, err := components.MessageRepository.FromConversations([]*models.Conversation{conv}, repositories.TimePagination{})
		require.NoError(t, err)
		require.Equal(t, 1, len(msgs))
	})
}

//
// -- UTILS
//

func newMsgEvent(convUuid string) *pb.Event {
	return &pb.Event{
		LocalUuid: uuid2.NewV4().String(),
		Content: &pb.Event_SentNewMessage_{SentNewMessage: &pb.Event_SentNewMessage{
			Message:      &pb.MessageSendInput{Content: &pb.MessageSendInput_TextMessage{TextMessage: &pb.TextMessage{Content: "hello"}}},
			Conversation: &pb.Event_SentNewMessage_ConvUuid{ConvUuid: convUuid}},
		}}
}
