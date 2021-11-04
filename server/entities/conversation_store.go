package entities

import (
	"context"
	uuid2 "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
)

type ConversationStore interface {
	ByIds(...int) (ConversationSlicePtrs, error)
	ByUuids(...uuid2.UUID) (ConversationSlicePtrs, error)
	ById(int) (*Conversation, error)
	ByUuid(uuid uuid2.UUID) (*Conversation, error)
}

type ConversationStoreImpl struct {
	*common.Components
	Context   context.Context
	IdCache   ConversationCacheByInt
	UuidCache ConversationCacheByUuid2UUID
}

func NewConversationStore(context context.Context, components *common.Components) *ConversationStoreImpl {
	return &ConversationStoreImpl{
		Components: components,
		Context:    context,
		IdCache: NewConversationCacheByInt(func(ids []int) ([]*Conversation, error) {
			out := []*Conversation{}
			records, err := models.Conversations(models.ConversationWhere.ID.IN(ids)).All(context, components.Db)
			if err != nil {
				return nil, err
			}

			for _, record := range records {
				out = append(out, &Conversation{record, components})
			}

			return out, nil
		}, func(conv *Conversation) int {
			return conv.Record.ID
		}),
		UuidCache: NewConversationCacheByUuid2UUID(func(uuids []uuid2.UUID) ([]*Conversation, error) {
			out := []*Conversation{}
			records, err := models.Conversations(models.ConversationWhere.UUID.IN(uuids)).All(context, components.Db)
			if err != nil {
				return nil, err
			}

			for _, record := range records {
				out = append(out, &Conversation{record, components})
			}

			return out, nil
		}, func(conv *Conversation) uuid2.UUID {
			return conv.Record.UUID
		}),
	}
}

func (store ConversationStoreImpl) ByIds(ints ...int) (ConversationSlicePtrs, error) {
	convs, err := store.IdCache.Get(ints)
	if err != nil {
		return nil, err
	}

	store.UuidCache.Prime(convs...)
	return convs, err
}

func (store ConversationStoreImpl) ByUuids(uuids ...uuid2.UUID) (ConversationSlicePtrs, error) {
	convs, err := store.UuidCache.Get(uuids)
	if err != nil {
		return nil, err
	}

	store.IdCache.Prime(convs...)
	return convs, err
}

func (store ConversationStoreImpl) ById(id int) (*Conversation, error) {
	convs, err := store.ByIds(id)
	if err != nil {
		return nil, err
	}

	return convs[0], nil
}

func (store ConversationStoreImpl) ByUuid(uuid uuid2.UUID) (*Conversation, error) {
	convs, err := store.ByUuids(uuid)
	if err != nil {
		return nil, err
	}

	return convs[0], nil
}

var _ ConversationStore = ConversationStoreImpl{}
