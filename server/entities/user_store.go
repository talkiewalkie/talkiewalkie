package entities

import (
	"context"
	uuid2 "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
)

type UserStore interface {
	ByIds(...int) (UserSlicePtrs, error)
	ByUuids(...uuid2.UUID) (UserSlicePtrs, error)
	ByPhoneNumbers(...string) (UserSlicePtrs, error)

	ById(int) (*User, error)
}

type UserStoreImpl struct {
	*common.Components
	Context    context.Context
	IdCache    UserCacheByInt
	UuidCache  UserCacheByUuid2UUID
	PhoneCache UserCacheByString
}

func NewUserStore(context context.Context, components *common.Components) *UserStoreImpl {
	return &UserStoreImpl{
		Components: components,
		Context:    context,
		IdCache: NewUserCacheByInt(func(ints []int) ([]*User, error) {
			out := []*User{}
			records, err := models.Users(models.UserWhere.ID.IN(ints)).All(context, components.Db)
			if err != nil {
				return nil, err
			}

			for _, record := range records {
				out = append(out, &User{record, components})
			}

			return out, nil
		}, func(value *User) int {
			return value.Record.ID
		}),
		UuidCache: NewUserCacheByUuid2UUID(func(uuids []uuid2.UUID) ([]*User, error) {
			out := []*User{}
			records, err := models.Users(models.UserWhere.UUID.IN(uuids)).All(context, components.Db)
			if err != nil {
				return nil, err
			}

			for _, record := range records {
				out = append(out, &User{record, components})
			}

			return out, nil
		}, func(value *User) uuid2.UUID {
			return value.Record.UUID
		}),
		PhoneCache: NewUserCacheByString(func(phones []string) ([]*User, error) {
			out := []*User{}

			records, err := models.Users(models.UserWhere.PhoneNumber.IN(phones)).All(context, components.Db)
			if err != nil {
				return nil, err
			}

			for _, record := range records {
				out = append(out, &User{record, components})
			}

			return out, nil
		}, func(value *User) string {
			return value.Record.PhoneNumber
		}),
	}
}

func (store UserStoreImpl) ByIds(ints ...int) (UserSlicePtrs, error) {
	users, err := store.IdCache.Get(ints)
	if err != nil {
		return nil, err
	}

	store.UuidCache.Prime(users...)
	store.PhoneCache.Prime(users...)
	return users, nil
}

func (store UserStoreImpl) ByUuids(uuids ...uuid2.UUID) (UserSlicePtrs, error) {
	users, err := store.UuidCache.Get(uuids)
	if err != nil {
		return nil, err
	}

	store.IdCache.Prime(users...)
	store.PhoneCache.Prime(users...)
	return users, nil
}

func (store UserStoreImpl) ByPhoneNumbers(phones ...string) (UserSlicePtrs, error) {
	users, err := store.PhoneCache.Get(phones)
	if err != nil {
		return nil, err
	}

	store.IdCache.Prime(users...)
	store.UuidCache.Prime(users...)
	return users, nil
}

func (store UserStoreImpl) ById(id int) (*User, error) {
	users, err := store.ByIds(id)
	if err != nil {
		return nil, err
	}

	return users[0], nil
}

var _ UserStore = UserStoreImpl{}
