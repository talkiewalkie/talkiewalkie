// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package caches

import (
	"errors"
	uuid2 "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/models"
)

var MessageCacheByIntErrNotFound = errors.New("MessageCacheByInt error did not find values for keys")

type MessageCacheByInt struct {
	cache   map[int]*models.Message
	fetcher func([]int) ([]*models.Message, error)
	primer  func(value *models.Message) int
}

func NewMessageCacheByInt(
	fetcher func([]int) ([]*models.Message, error),
	primer func(value *models.Message) int,
) MessageCacheByInt {
	return MessageCacheByInt{
		cache:   map[int]*models.Message{},
		fetcher: fetcher,
		primer:  primer,
	}
}

func (cache *MessageCacheByInt) Get(
	identifiers []int,
) ([]*models.Message, error) {
	out := make([]*models.Message, len(identifiers))
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
		missingKeys := []int{}
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

	for _, value := range out {
		if value == nil {
			return nil, MessageCacheByIntErrNotFound
		}
	}

	return out, nil
}

func (cache MessageCacheByInt) Clear() {
	cache.cache = nil
}

func (cache MessageCacheByInt) Prime(values ...*models.Message) {
	for _, value := range values {
		cache.cache[cache.primer(value)] = value
	}
}

// ENSURE IMPORTS
var (
	_ uuid2.UUID
	_ models.User
)

var MessageCacheByUuidErrNotFound = errors.New("MessageCacheByUuid error did not find values for keys")

type MessageCacheByUuid struct {
	cache   map[uuid2.UUID]*models.Message
	fetcher func([]uuid2.UUID) ([]*models.Message, error)
	primer  func(value *models.Message) uuid2.UUID
}

func NewMessageCacheByUuid(
	fetcher func([]uuid2.UUID) ([]*models.Message, error),
	primer func(value *models.Message) uuid2.UUID,
) MessageCacheByUuid {
	return MessageCacheByUuid{
		cache:   map[uuid2.UUID]*models.Message{},
		fetcher: fetcher,
		primer:  primer,
	}
}

func (cache *MessageCacheByUuid) Get(
	identifiers []uuid2.UUID,
) ([]*models.Message, error) {
	out := make([]*models.Message, len(identifiers))
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
		missingKeys := []uuid2.UUID{}
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

	for _, value := range out {
		if value == nil {
			return nil, MessageCacheByUuidErrNotFound
		}
	}

	return out, nil
}

func (cache MessageCacheByUuid) Clear() {
	cache.cache = nil
}

func (cache MessageCacheByUuid) Prime(values ...*models.Message) {
	for _, value := range values {
		cache.cache[cache.primer(value)] = value
	}
}

// ENSURE IMPORTS
var (
	_ uuid2.UUID
	_ models.User
)
