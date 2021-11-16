package repositories

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/jmoiron/sqlx"
	uuid2 "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"github.com/talkiewalkie/talkiewalkie/pkg/slices"
	"github.com/talkiewalkie/talkiewalkie/repositories/caches"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TimePagination struct {
	Offset   time.Time
	PageSize int
}

type MessageRepository interface {
	ByIds(...int) (models.MessageSlice, error)
	ByUuids(...uuid2.UUID) (models.MessageSlice, error)
	ById(int) (*models.Message, error)
	ByUuid(uuid uuid2.UUID) (*models.Message, error)

	FromConversations(models.ConversationSlice, TimePagination) (models.MessageSlice, error)
	FromConversationsLast(models.ConversationSlice) (models.MessageSlice, error)
	Clear()
}

type MessageRepositoryImpl struct {
	Db      *sqlx.DB
	Context context.Context

	IdCache   caches.MessageCacheByInt
	UuidCache caches.MessageCacheByUuid
}

func (repository MessageRepositoryImpl) Clear() {
	repository.IdCache.Clear()
	repository.UuidCache.Clear()
}

func NewMessageRepository(context context.Context, db *sqlx.DB) *MessageRepositoryImpl {
	return &MessageRepositoryImpl{
		Db:      db,
		Context: context,
		IdCache: caches.NewMessageCacheByInt(func(ints []int) ([]*models.Message, error) {
			return models.Messages(models.MessageWhere.ID.IN(ints)).All(context, db)
		}, func(value *models.Message) int {
			return value.ID
		}),
		UuidCache: caches.NewMessageCacheByUuid(func(uuids []uuid2.UUID) ([]*models.Message, error) {
			return models.Messages(models.MessageWhere.UUID.IN(uuids)).All(context, db)
		}, func(value *models.Message) uuid2.UUID {
			return value.UUID
		}),
	}
}

func (repository MessageRepositoryImpl) ByIds(ints ...int) (models.MessageSlice, error) {
	messages, err := repository.IdCache.Get(ints)
	if err != nil {
		return nil, err
	}

	repository.UuidCache.Prime(messages...)
	return messages, nil
}

func (repository MessageRepositoryImpl) ByUuids(uuids ...uuid2.UUID) (models.MessageSlice, error) {
	messages, err := repository.UuidCache.Get(uuids)
	if err != nil {
		return nil, err
	}

	repository.IdCache.Prime(messages...)
	return messages, nil
}

func (repository MessageRepositoryImpl) ById(id int) (*models.Message, error) {
	messages, err := repository.ByIds(id)
	if err != nil {
		return nil, err
	}
	return messages[0], nil
}

func (repository MessageRepositoryImpl) ByUuid(uuid uuid2.UUID) (*models.Message, error) {
	messages, err := repository.ByUuids(uuid)
	if err != nil {
		return nil, err
	}
	return messages[0], nil
}

func (repository MessageRepositoryImpl) FromConversations(slice models.ConversationSlice, pagination TimePagination) (models.MessageSlice, error) {
	return models.Messages(models.MessageWhere.ConversationID.IN(slice.Ids())).All(repository.Context, repository.Db)
}

func (repository MessageRepositoryImpl) FromConversationsLast(convs models.ConversationSlice) (models.MessageSlice, error) {
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

	if err := models.NewQuery(qm.SQL(fmt.Sprintf(sqlStr, strings.Join(convIds, ",")))).Bind(repository.Context, repository.Db, &mssgs); err != nil {
		return nil, status.Errorf(codes.Internal, "could not fetch last messages for each conv: %+v", err)
	}

	out := []*models.Message{}
	for _, item := range mssgs {
		out = append(out, &item.Message)
	}

	return out, nil
}

var _ MessageRepository = MessageRepositoryImpl{}

// UTILS

func (s Repositories) messagesLoadAuthors(messages []*models.Message) ([]*models.User, error) {
	userIds := slices.IntSlice{}
	for _, message := range messages {
		if message.AuthorID.Valid {
			userIds = append(userIds, message.AuthorID.Int)
		}
	}
	uniqueIds := userIds.UniqueBy(func(i int) interface{} { return i })

	return s.UserRepository.ByIds(uniqueIds...)
}

func (s Repositories) messagesLoadRawAudios(messages models.MessageSlice) ([]*models.Asset, error) {
	var assetIds slices.IntSlice = messages.RawAudioIDs()
	uniqueIds := assetIds.UniqueBy(func(i int) interface{} { return i })

	return s.AssetRepository.ByIds(uniqueIds...)
}

func (s Repositories) messagesLoadConversations(messages []*models.Message) ([]*models.Conversation, error) {
	conversationIds := slices.IntSlice{}
	for _, message := range messages {
		conversationIds = append(conversationIds, message.ConversationID)
	}
	uniqueIds := conversationIds.UniqueBy(func(i int) interface{} { return i })

	return s.ConversationRepository.ByIds(uniqueIds...)
}

type PbMessageSlice []*pb.Message

func (s PbMessageSlice) UuidMap() map[uuid2.UUID]*pb.Message {
	out := make(map[uuid2.UUID]*pb.Message, len(s))
	for _, pbm := range s {
		uid, _ := uuid2.FromString(pbm.Uuid)
		out[uid] = pbm
	}

	return out
}

func (s Repositories) MessagesToProto(messages models.MessageSlice) (PbMessageSlice, error) {
	if _, err := s.UserRepository.ByIds(messages.AuthorIDs()...); err != nil {
		return nil, err
	}
	if _, err := s.AssetRepository.ByIds(messages.RawAudioIDs()...); err != nil {
		return nil, err
	}
	if _, err := s.ConversationRepository.ByIds(messages.ConversationIDs()...); err != nil {
		return nil, err
	}

	out := []*pb.Message{}
	for _, message := range messages {
		var content pb.MessageContentOneOf
		switch message.Type {
		case models.MessageTypeText:
			content = &pb.Message_TextMessage{TextMessage: &pb.TextMessage{Content: message.Text.String}}
		case models.MessageTypeVoice:
			rawAudio, err := s.AssetRepository.ById(message.RawAudioID.Int)
			if err != nil {
				return nil, fmt.Errorf("could not fetch associated audio asset: %+v", err)
			}

			buf := &bytes.Buffer{}
			if err := s.CloudStorage.Download(rawAudio.BlobName.String, buf); err != nil {
				return nil, fmt.Errorf("could not download audio file from bucket: %+v", err)
			}

			var pbTranscript pb.AlignedTranscript
			if err := proto.Unmarshal(message.SiriTranscript.Bytes, &pbTranscript); err != nil {
				return nil, fmt.Errorf("could not build protobuf from repositoryd bytearray: %+v", err)
			}
			content = &pb.Message_VoiceMessage{VoiceMessage: &pb.VoiceMessage{RawContent: buf.Bytes(), SiriTranscript: &pbTranscript}}
		default:
			// TODO
			return nil, errors.New("message content type other than text or voice unhandled atm.")
		}

		conv, err := s.ConversationRepository.ById(message.ConversationID)
		if err != nil {
			return nil, err
		}

		var author *pb.User
		if message.AuthorID.Valid {
			dbAuthor, err := s.UserRepository.ById(message.AuthorID.Int)
			if err != nil {
				return nil, err
			}
			author = UserToProto(dbAuthor)
		}

		out = append(out, &pb.Message{
			Uuid:      message.UUID.String(),
			ConvUuid:  conv.UUID.String(),
			Content:   content,
			Author:    author,
			CreatedAt: timestamppb.New(message.CreatedAt),
		})
	}

	return out, nil
}
