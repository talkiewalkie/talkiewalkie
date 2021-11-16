// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package caches

import (
	"fmt"
	uuid2 "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/models"
)

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
			return nil, fmt.Errorf("[MessageCacheByInt] error: found nil value at position %d", index)
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
			return nil, fmt.Errorf("[MessageCacheByUuid] error: found nil value at position %d", index)
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
