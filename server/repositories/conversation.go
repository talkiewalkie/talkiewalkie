package repositories

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	uuid2 "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"github.com/talkiewalkie/talkiewalkie/repositories/caches"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"strings"
)

type ConversationRepository interface {
	ByIds(...int) (models.ConversationSlice, error)
	ByUuids(...uuid2.UUID) (models.ConversationSlice, error)
	ById(int) (*models.Conversation, error)
	ByUuid(uuid uuid2.UUID) (*models.Conversation, error)

	FromUserConversations([][]*models.UserConversation) (models.ConversationSlice, error)

	Clear()
}

type ConversationRepositoryImpl struct {
	Db      *sqlx.DB
	Context context.Context

	IdCache   caches.ConversationCacheByInt
	UuidCache caches.ConversationCacheByUuid
}

func (repository ConversationRepositoryImpl) Clear() {
	repository.IdCache.Clear()
	repository.UuidCache.Clear()
}

func NewConversationRepository(context context.Context, db *sqlx.DB) *ConversationRepositoryImpl {
	return &ConversationRepositoryImpl{
		Db:      db,
		Context: context,
		IdCache: caches.NewConversationCacheByInt(func(ids []int) ([]*models.Conversation, error) {
			return models.Conversations(models.ConversationWhere.ID.IN(ids)).All(context, db)
		}, func(conv *models.Conversation) int {
			return conv.ID
		}),
		UuidCache: caches.NewConversationCacheByUuid(func(uuids []uuid2.UUID) ([]*models.Conversation, error) {
			return models.Conversations(models.ConversationWhere.UUID.IN(uuids)).All(context, db)
		}, func(conv *models.Conversation) uuid2.UUID {
			return conv.UUID
		}),
	}
}

func (repository ConversationRepositoryImpl) ByIds(ints ...int) (models.ConversationSlice, error) {
	convs, err := repository.IdCache.Get(ints)
	if err != nil {
		return nil, err
	}

	repository.UuidCache.Prime(convs...)
	return convs, err
}

func (repository ConversationRepositoryImpl) ByUuids(uuids ...uuid2.UUID) (models.ConversationSlice, error) {
	convs, err := repository.UuidCache.Get(uuids)
	if err != nil {
		return nil, err
	}

	repository.IdCache.Prime(convs...)
	return convs, err
}

func (repository ConversationRepositoryImpl) ById(id int) (*models.Conversation, error) {
	convs, err := repository.ByIds(id)
	if err != nil {
		return nil, err
	}

	return convs[0], nil
}

func (repository ConversationRepositoryImpl) ByUuid(uuid uuid2.UUID) (*models.Conversation, error) {
	convs, err := repository.ByUuids(uuid)
	if err != nil {
		return nil, err
	}

	return convs[0], nil
}

func (repository ConversationRepositoryImpl) FromUserConversations(ucs [][]*models.UserConversation) (models.ConversationSlice, error) {
	slices := models.UserConversationSlice{}
	for _, uc := range ucs {
		slices = append(slices, uc...)
	}

	return repository.ByIds(slices.ConversationIDs()...)
}

var _ ConversationRepository = ConversationRepositoryImpl{}

// UTILS

type PbConversationSlice []*pb.Conversation

func (s PbConversationSlice) UuidMap() map[uuid2.UUID]*pb.Conversation {
	out := make(map[uuid2.UUID]*pb.Conversation, len(s))
	for _, item := range s {
		uid, _ := uuid2.FromString(item.Uuid)
		out[uid] = item
	}

	return out
}

func (s Repositories) ConversationsToProto(convs models.ConversationSlice) (PbConversationSlice, error) {
	convIds := []int{}
	for _, conv := range convs {
		convIds = append(convIds, conv.ID)
	}

	pp, err := s.UserConversationRepository.ByConversationIds(convIds...)
	if err != nil {
		return nil, err
	}

	userIds := []int{}
	for _, users := range pp {
		for _, user := range users {
			userIds = append(userIds, user.UserID)
		}
	}
	if _, err := s.UserRepository.ByIds(userIds...); err != nil {
		return nil, err
	}

	out := []*pb.Conversation{}
	for idx := range convs {
		conv := convs[idx]
		users := pp[idx]

		pbs, err := s.UserConversationsToProto(users)
		if err != nil {
			return nil, err
		}

		out = append(out, &pb.Conversation{
			Uuid:         conv.UUID.String(),
			Title:        conv.Name.String,
			Messages:     nil,
			Participants: pbs,
		})
	}
	return out, nil
}

func (s Repositories) ConversationsLoadLastMessages(convs []*models.Conversation) ([]*models.Message, error) {
	if len(convs) == 0 {
		return nil, nil
	}

	type lastMessage struct {
		Discard1       int `boil:"discard_1"`
		Discard2       int `boil:"discard_2"`
		models.Message `boil:",bind"`
	}
	var mssgs []*lastMessage

	sqlStr := `
		SELECT 
			DISTINCT ON (conversation_id) conversation_id discard_1, 
			last_value("id") OVER (PARTITION BY "conversation_id" ORDER BY created_at DESC) "discard_2", 
			*
		FROM "message"
		WHERE "conversation_id" in (%s);
	`
	convIds := []string{}
	for _, conv := range convs {
		convIds = append(convIds, strconv.Itoa(conv.ID))
	}

	if err := models.NewQuery(qm.SQL(fmt.Sprintf(sqlStr, strings.Join(convIds, ",")))).Bind(s.context, s.db, &mssgs); err != nil {
		return nil, status.Errorf(codes.Internal, "could not fetch last messages for each conv: %+v", err)
	}

	out := []*models.Message{}
	for _, item := range mssgs {
		out = append(out, &item.Message)
	}

	return out, nil
}

func (s Repositories) ConversationUsers(conv *models.Conversation) ([]*models.UserConversation, error) {
	ucs, err := s.UserConversationRepository.ByConversationIds(conv.ID)
	if err != nil {
		return nil, err
	}

	return ucs[0], nil
}

func (s Repositories) ConversationHasAccess(me *models.User, convs ...*models.Conversation) (bool, error) {
	// TODO: fix this: when we go through this we populate the UserConversation cache with false info (all my
	// 		 conversations will have only me as participant because that's how the query is made.)
	myConvs, err := s.UserConversations(me)
	if err != nil {
		return false, err
	}

	myConvIds := map[int]int{}
	for _, conv := range myConvs {
		myConvIds[conv.ID]++
	}

	for _, conv := range convs {
		if _, ok := myConvIds[conv.ID]; !ok {
			return false, nil
		}
	}
	return true, nil
}
