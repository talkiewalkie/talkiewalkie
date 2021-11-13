package repositories

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/talkiewalkie/talkiewalkie/clients"
)

type Repositories struct {
	context context.Context
	db      *sqlx.DB

	CloudStorage clients.StorageClient
	PubSubClient clients.PubSubClient

	AssetRepository            AssetRepository
	ConversationRepository     ConversationRepository
	UserRepository             UserRepository
	MessageRepository          MessageRepository
	UserConversationRepository UserConversationRepository
}

func New(
	ctx context.Context,
	db *sqlx.DB,
	storage clients.StorageClient,
	pubsub clients.PubSubClient,
) Repositories {
	return Repositories{
		context: ctx,
		db:      db,

		// TODO
		CloudStorage: storage,
		PubSubClient: pubsub,

		AssetRepository:            NewAssetRepository(ctx, db),
		ConversationRepository:     NewConversationRepository(ctx, db),
		UserRepository:             NewUserRepository(ctx, db),
		MessageRepository:          NewMessageRepository(ctx, db),
		UserConversationRepository: NewUserConversationRepository(ctx, db),
	}
}
