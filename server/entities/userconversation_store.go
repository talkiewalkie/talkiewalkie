package entities

import (
	"context"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
)

type UserConversationStore interface {
	ByConversationIds(...int) ([][]*UserConversation, error)
	ByUserIds(...int) ([][]*UserConversation, error)

	LoadConversationsFromResult([][]*UserConversation) (ConversationSlicePtrs, error)
	LoadUsersFromResult([][]*UserConversation) (UserSlicePtrs, error)

	Clear()
}

type UserConversationStoreImpl struct {
	*common.Components
	Context             context.Context
	ConversationIdCache UserConversationMultiCacheByInt
	UserIdCache         UserConversationMultiCacheByInt
}

func (store UserConversationStoreImpl) Clear() {
	store.ConversationIdCache.Clear()
}

func NewUserConversationStore(context context.Context, components *common.Components) *UserConversationStoreImpl {
	return &UserConversationStoreImpl{
		Components: components,
		Context:    context,
		ConversationIdCache: NewUserConversationMultiCacheByInt(func(ints []int) ([]*UserConversation, error) {
			out := []*UserConversation{}
			records, err := models.UserConversations(models.UserConversationWhere.ConversationID.IN(ints)).All(context, components.Db)
			if err != nil {
				return nil, err
			}

			for _, record := range records {
				out = append(out, &UserConversation{record, components})
			}

			return out, nil
		}, func(value *UserConversation) int {
			return value.Record.ConversationID
		}),
		UserIdCache: NewUserConversationMultiCacheByInt(func(ints []int) ([]*UserConversation, error) {
			out := []*UserConversation{}
			records, err := models.UserConversations(models.UserConversationWhere.UserID.IN(ints)).All(context, components.Db)
			if err != nil {
				return nil, err
			}

			for _, record := range records {
				out = append(out, &UserConversation{record, components})
			}

			return out, nil
		}, func(value *UserConversation) int {
			return value.Record.UserID
		}),
	}
}

func (store UserConversationStoreImpl) ByConversationIds(convIds ...int) ([][]*UserConversation, error) {
	out, err := store.ConversationIdCache.Get(convIds)
	if err != nil {
		return nil, err
	}

	store.UserIdCache.Prime(out...)
	return out, nil
}
func (store UserConversationStoreImpl) ByUserIds(userIds ...int) ([][]*UserConversation, error) {
	out, err := store.UserIdCache.Get(userIds)
	if err != nil {
		return nil, err
	}

	store.ConversationIdCache.Prime(out...)
	return out, nil
}

func (store UserConversationStoreImpl) LoadConversationsFromResult(in [][]*UserConversation) (ConversationSlicePtrs, error) {
	convIds := []int{}
	for _, convParticipants := range in {
		if len(convParticipants) >= 0 {
			convIds = append(convIds, convParticipants[0].Record.ConversationID)
		}
	}

	return store.ConversationStore.ByIds(convIds...)
}

func (store UserConversationStoreImpl) LoadUsersFromResult(in [][]*UserConversation) (UserSlicePtrs, error) {
	userIds := []int{}
	for _, convParticipants := range in {
		for _, uc := range convParticipants {
			userIds = append(userIds, uc.Record.UserID)
		}
	}

	return store.UserStore.ByIds(userIds...)
}

var _ UserConversationStore = UserConversationStoreImpl{}
