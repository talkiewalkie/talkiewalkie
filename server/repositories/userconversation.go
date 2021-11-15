package repositories

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"github.com/talkiewalkie/talkiewalkie/pkg/slices"
	"github.com/talkiewalkie/talkiewalkie/repositories/caches"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserConversationRepository interface {
	ByConversationIds(...int) ([][]*models.UserConversation, error)
	ByUserIds(...int) ([][]*models.UserConversation, error)

	//LoadConversationsFromResult([][]*models.UserConversation) (models.ConversationSlicePtrs, error)
	//LoadUsersFromResult([][]*models.UserConversation) (models.UserSlicePtrs, error)

	Clear()
}

type UserConversationRepositoryImpl struct {
	Db                  *sqlx.DB
	Context             context.Context
	ConversationIdCache caches.UserConversationMultiCacheByInt
	UserIdCache         caches.UserConversationMultiCacheByInt
}

func (repository UserConversationRepositoryImpl) Clear() {
	repository.ConversationIdCache.Clear()
}

func NewUserConversationRepository(context context.Context, db *sqlx.DB) *UserConversationRepositoryImpl {
	return &UserConversationRepositoryImpl{
		Db:      db,
		Context: context,
		ConversationIdCache: caches.NewUserConversationMultiCacheByInt(func(ints []int) ([]*models.UserConversation, error) {
			return models.UserConversations(models.UserConversationWhere.ConversationID.IN(ints)).All(context, db)
		}, func(value *models.UserConversation) int {
			return value.ConversationID
		}),
		UserIdCache: caches.NewUserConversationMultiCacheByInt(func(ints []int) ([]*models.UserConversation, error) {
			return models.UserConversations(models.UserConversationWhere.UserID.IN(ints)).All(context, db)
		}, func(value *models.UserConversation) int {
			return value.UserID
		}),
	}
}

func (repository UserConversationRepositoryImpl) ByConversationIds(convIds ...int) ([][]*models.UserConversation, error) {
	out, err := repository.ConversationIdCache.Get(convIds)
	if err != nil {
		return nil, err
	}

	repository.UserIdCache.Prime(out...)
	return out, nil
}
func (repository UserConversationRepositoryImpl) ByUserIds(userIds ...int) ([][]*models.UserConversation, error) {
	out, err := repository.UserIdCache.Get(userIds)
	if err != nil {
		return nil, err
	}

	repository.ConversationIdCache.Prime(out...)
	return out, nil
}

//func (repository UserConversationRepositoryImpl) LoadConversationsFromResult(in [][]*models.UserConversation) (models.ConversationSlicePtrs, error) {
//	convIds := []int{}
//	for _, convParticipants := range in {
//		if len(convParticipants) >= 0 {
//			convIds = append(convIds, convParticipants[0].Record.ConversationID)
//		}
//	}
//
//	return repository.ConversationRepository.ByIds(convIds...)
//}
//
//func (repository UserConversationRepositoryImpl) LoadUsersFromResult(in [][]*models.UserConversation) (models.UserSlicePtrs, error) {
//	userIds := []int{}
//	for _, convParticipants := range in {
//		for _, uc := range convParticipants {
//			userIds = append(userIds, uc.Record.UserID)
//		}
//	}
//
//	return repository.UserRepository.ByIds(userIds...)
//}

var _ UserConversationRepository = UserConversationRepositoryImpl{}

// UTILS

func (s Repositories) userConversationsLoadUsers(ucs []*models.UserConversation) ([]*models.User, error) {
	userIds := slices.IntSlice{}
	for _, uc := range ucs {
		userIds = append(userIds, uc.UserID)
	}
	uniqueIds := userIds.UniqueBy(func(i int) interface{} { return i })

	return s.UserRepository.ByIds(uniqueIds...)
}

func (s Repositories) userConversationsLoadConversations(ucs []*models.UserConversation) ([]*models.Conversation, error) {
	convIds := slices.IntSlice{}
	for _, uc := range ucs {
		convIds = append(convIds, uc.ConversationID)
	}
	uniqueIds := convIds.UniqueBy(func(i int) interface{} { return i })

	return s.ConversationRepository.ByIds(uniqueIds...)
}

func (s Repositories) UserConversationsToProto(ucs []*models.UserConversation) ([]*pb.UserConversation, error) {
	if _, err := s.userConversationsLoadUsers(ucs); err != nil {
		return nil, err
	}

	out := []*pb.UserConversation{}
	for _, uc := range ucs {
		user, _ := s.UserRepository.ById(uc.UserID)
		out = append(out, &pb.UserConversation{
			User:      UserToProto(user),
			ReadUntil: timestamppb.New(uc.ReadUntil),
		})
	}

	return out, nil
}
