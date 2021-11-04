package entities

import (
	"context"
	uuid2 "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
)

type MessageStore interface {
	ByIds(...int) (MessageSlicePtrs, error)
	ByUuids(...uuid2.UUID) (MessageSlicePtrs, error)
	ById(int) (*Message, error)
	ByUuid(uuid uuid2.UUID) (*Message, error)
}

type MessageStoreImpl struct {
	*common.Components
	Context   context.Context
	IdCache   MessageCacheByInt
	UuidCache MessageCacheByUuid2UUID
}

func NewMessageStore(context context.Context, components *common.Components) *MessageStoreImpl {
	return &MessageStoreImpl{
		Components: components,
		Context:    context,
		IdCache: NewMessageCacheByInt(func(ints []int) ([]*Message, error) {
			out := []*Message{}
			records, err := models.Messages(models.MessageWhere.ID.IN(ints)).All(context, components.Db)
			if err != nil {
				return nil, err
			}

			for _, record := range records {
				out = append(out, &Message{record, components})
			}

			return out, nil
		}, func(value *Message) int {
			return value.Record.ID
		}),
		UuidCache: NewMessageCacheByUuid2UUID(func(uuids []uuid2.UUID) ([]*Message, error) {
			out := []*Message{}
			records, err := models.Messages(models.MessageWhere.UUID.IN(uuids)).All(context, components.Db)
			if err != nil {
				return nil, err
			}

			for _, record := range records {
				out = append(out, &Message{record, components})
			}

			return out, nil
		}, func(value *Message) uuid2.UUID {
			return value.Record.UUID
		}),
	}
}

func (store MessageStoreImpl) ByIds(ints ...int) (MessageSlicePtrs, error) {
	messages, err := store.IdCache.Get(ints)
	if err != nil {
		return nil, err
	}

	store.UuidCache.Prime(messages...)
	return messages, nil
}

func (store MessageStoreImpl) ByUuids(uuids ...uuid2.UUID) (MessageSlicePtrs, error) {
	messages, err := store.UuidCache.Get(uuids)
	if err != nil {
		return nil, err
	}

	store.IdCache.Prime(messages...)
	return messages, nil
}

func (store MessageStoreImpl) ById(id int) (*Message, error) {
	messages, err := store.ByIds(id)
	if err != nil {
		return nil, err
	}
	return messages[0], nil
}

func (store MessageStoreImpl) ByUuid(uuid uuid2.UUID) (*Message, error) {
	messages, err := store.ByUuids(uuid)
	if err != nil {
		return nil, err
	}
	return messages[0], nil
}

var _ MessageStore = MessageStoreImpl{}
