// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package caches

import (
	"fmt"
	uuid2 "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/models"
)

type ConversationCacheByInt struct {
	cache   map[int]*models.Conversation
	fetcher func([]int) ([]*models.Conversation, error)
	primer  func(value *models.Conversation) int
}

func NewConversationCacheByInt(
	fetcher func([]int) ([]*models.Conversation, error),
	primer func(value *models.Conversation) int,
) ConversationCacheByInt {
	return ConversationCacheByInt{
		cache:   map[int]*models.Conversation{},
		fetcher: fetcher,
		primer:  primer,
	}
}

func (cache *ConversationCacheByInt) Get(
	identifiers []int,
) ([]*models.Conversation, error) {
	out := make([]*models.Conversation, len(identifiers))
	key2index := map[int]int{}

	for idx, key := range identifiers {
		item, ok := cache.cache[key]
		if ok {
			out[idx] = item
		} else {
			key2index[key] = idx
		}
	}

	if len(key2index) > 0 {
		missingKeys := make([]int, len(key2index))
		for key, _ := range key2index {
			missingKeys = append(missingKeys, key)
		}

		records, err := cache.fetcher(missingKeys)
		if err != nil {
			return nil, err
		}

		for idx := range records {
			key := missingKeys[idx]
			record := records[idx]

			out[key2index[key]] = record
			cache.cache[key] = record
		}
	}

	for index, value := range out {
		if value == nil {
			return nil, fmt.Errorf("[ConversationCacheByInt] error: found nil value at position %d", index)
		}
	}

	return out, nil
}

func (cache ConversationCacheByInt) Clear() {
	cache.cache = nil
}

func (cache ConversationCacheByInt) Prime(values ...*models.Conversation) {
	for _, value := range values {
		cache.cache[cache.primer(value)] = value
	}
}

// ENSURE IMPORTS
var (
	_ uuid2.UUID
	_ models.User
)

type ConversationCacheByUuid struct {
	cache   map[uuid2.UUID]*models.Conversation
	fetcher func([]uuid2.UUID) ([]*models.Conversation, error)
	primer  func(value *models.Conversation) uuid2.UUID
}

func NewConversationCacheByUuid(
	fetcher func([]uuid2.UUID) ([]*models.Conversation, error),
	primer func(value *models.Conversation) uuid2.UUID,
) ConversationCacheByUuid {
	return ConversationCacheByUuid{
		cache:   map[uuid2.UUID]*models.Conversation{},
		fetcher: fetcher,
		primer:  primer,
	}
}

func (cache *ConversationCacheByUuid) Get(
	identifiers []uuid2.UUID,
) ([]*models.Conversation, error) {
	out := make([]*models.Conversation, len(identifiers))
	key2index := map[uuid2.UUID]int{}

	for idx, key := range identifiers {
		item, ok := cache.cache[key]
		if ok {
			out[idx] = item
		} else {
			key2index[key] = idx
		}
	}

	if len(key2index) > 0 {
		missingKeys := make([]uuid2.UUID, len(key2index))
		for key, _ := range key2index {
			missingKeys = append(missingKeys, key)
		}

		records, err := cache.fetcher(missingKeys)
		if err != nil {
			return nil, err
		}

		for idx := range records {
			key := missingKeys[idx]
			record := records[idx]

			out[key2index[key]] = record
			cache.cache[key] = record
		}
	}

	for index, value := range out {
		if value == nil {
			return nil, fmt.Errorf("[ConversationCacheByUuid] error: found nil value at position %d", index)
		}
	}

	return out, nil
}

func (cache ConversationCacheByUuid) Clear() {
	cache.cache = nil
}

func (cache ConversationCacheByUuid) Prime(values ...*models.Conversation) {
	for _, value := range values {
		cache.cache[cache.primer(value)] = value
	}
}

// ENSURE IMPORTS
var (
	_ uuid2.UUID
	_ models.User
)
