package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	uuid2 "github.com/satori/go.uuid"

	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"github.com/talkiewalkie/talkiewalkie/repositories/caches"
)

type UserRepository interface {
	ByIds(...int) (models.UserSlice, error)
	ByUuids(...uuid2.UUID) (models.UserSlice, error)
	ByPhoneNumbers(...string) (models.UserSlice, error)

	ById(int) (*models.User, error)

	FromUserConversations([][]*models.UserConversation) (models.UserSlice, error)
	Clear()
}

type UserRepositoryImpl struct {
	Db      *sqlx.DB
	Context context.Context

	IdCache    caches.UserCacheByInt
	UuidCache  caches.UserCacheByUuid
	PhoneCache caches.UserCacheByString
}

func (repository UserRepositoryImpl) Clear() {
	repository.IdCache.Clear()
	repository.UuidCache.Clear()
	repository.PhoneCache.Clear()
}

func NewUserRepository(context context.Context, db *sqlx.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		Db:      db,
		Context: context,
		IdCache: caches.NewUserCacheByInt(func(ints []int) ([]*models.User, error) {
			return models.Users(models.UserWhere.ID.IN(ints)).All(context, db)
		}, func(value *models.User) int {
			return value.ID
		}),
		UuidCache: caches.NewUserCacheByUuid(func(uuids []uuid2.UUID) ([]*models.User, error) {
			return models.Users(models.UserWhere.UUID.IN(uuids)).All(context, db)
		}, func(value *models.User) uuid2.UUID {
			return value.UUID
		}),
		PhoneCache: caches.NewUserCacheByString(func(phones []string) ([]*models.User, error) {
			return models.Users(models.UserWhere.PhoneNumber.IN(phones)).All(context, db)
		}, func(value *models.User) string {
			return value.PhoneNumber
		}),
	}
}

func (repository UserRepositoryImpl) ByIds(ints ...int) (models.UserSlice, error) {
	users, err := repository.IdCache.Get(ints)
	if err != nil {
		return nil, err
	}

	repository.UuidCache.Prime(users...)
	repository.PhoneCache.Prime(users...)
	return users, nil
}

func (repository UserRepositoryImpl) ByUuids(uuids ...uuid2.UUID) (models.UserSlice, error) {
	users, err := repository.UuidCache.Get(uuids)
	if err != nil {
		return nil, err
	}

	repository.IdCache.Prime(users...)
	repository.PhoneCache.Prime(users...)
	return users, nil
}

func (repository UserRepositoryImpl) ByPhoneNumbers(phones ...string) (models.UserSlice, error) {
	users, err := repository.PhoneCache.Get(phones)
	if err != nil {
		return nil, err
	}

	repository.IdCache.Prime(users...)
	repository.UuidCache.Prime(users...)
	return users, nil
}

func (repository UserRepositoryImpl) ById(id int) (*models.User, error) {
	users, err := repository.ByIds(id)
	if err != nil {
		return nil, err
	}

	return users[0], nil
}

func (repository UserRepositoryImpl) FromUserConversations(ucs [][]*models.UserConversation) (models.UserSlice, error) {
	slices := models.UserConversationSlice{}
	for _, uc := range ucs {
		slices = append(slices, uc...)
	}

	return repository.ByIds(slices.UserIDs()...)
}

var _ UserRepository = UserRepositoryImpl{}

// UTILS

func UserPubSubTopic(user *models.User) string {
	return strings.Replace(fmt.Sprintf("user-conn-%s", user.UUID), "-", "_", -1)
}

func UserToProto(u *models.User) *pb.User {
	return &pb.User{
		DisplayName:   u.DisplayName,
		Handle:        u.Handle,
		Uuid:          u.UUID.String(),
		Conversations: nil,
		Phone:         u.PhoneNumber,
	}
}

func (s Repositories) UserFriends(u *models.User) ([]*models.User, error) {
	convs, err := s.UserConversationRepository.ByUserIds(u.ID)
	if err != nil {
		return nil, err
	}

	allUcs := []*models.UserConversation{}
	for _, conv := range convs {
		for _, uc := range conv {
			allUcs = append(allUcs, uc)
		}
	}

	return s.userConversationsLoadUsers(allUcs)
}

func (s Repositories) UserConversations(u *models.User) ([]*models.Conversation, error) {
	ucs, err := s.UserConversationRepository.ByUserIds(u.ID)
	if err != nil {
		return nil, err
	}

	allUcs := []*models.UserConversation{}
	for _, uc := range ucs {
		for _, ucc := range uc {
			allUcs = append(allUcs, ucc)
		}
	}

	return s.userConversationsLoadConversations(allUcs)
}

func (s Repositories) UserHasAccess(me *models.User, users ...*models.User) (bool, error) {
	friends, err := s.UserFriends(me)
	if err != nil {
		return false, err
	}

	friendIds := map[int]int{}
	for _, friend := range friends {
		friendIds[friend.ID]++
	}

	for _, user := range users {
		if _, ok := friendIds[user.ID]; !ok {
			return false, nil
		}
	}

	return true, nil
}
