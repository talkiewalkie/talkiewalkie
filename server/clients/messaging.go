package clients

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"log"
)

type MessageInput struct {
	Topic string
	Data  map[string]string

	Title string
	Body  string
}

type MessagingClient interface {
	Send(context.Context, MessageInput) (string, error)
	SendAll(context.Context, []MessageInput) (*messaging.BatchResponse, error)
}

type FirebaseMessagingClientImpl struct {
	*messaging.Client
}

func (f FirebaseMessagingClientImpl) Send(ctx context.Context, input MessageInput) (string, error) {
	return f.Client.Send(ctx, &messaging.Message{
		Data:         input.Data,
		Notification: &messaging.Notification{Title: input.Title, Body: input.Body},
		Topic:        input.Topic,
	})
}

func (f FirebaseMessagingClientImpl) SendAll(ctx context.Context, input []MessageInput) (*messaging.BatchResponse, error) {
	messages := []*messaging.Message{}
	for _, messageInput := range input {
		messages = append(messages, &messaging.Message{
			Data:         messageInput.Data,
			Notification: &messaging.Notification{Title: messageInput.Title, Body: messageInput.Body},
			Topic:        messageInput.Topic,
		})
	}

	return f.Client.SendAll(ctx, messages)
}

var _ MessagingClient = FirebaseMessagingClientImpl{}

func NewFirebaseMessagingClient(app *firebase.App) *FirebaseMessagingClientImpl {
	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Panicf("%+v", err)
	}

	return &FirebaseMessagingClientImpl{client}
}
