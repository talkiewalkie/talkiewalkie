// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package caches

import (
	"errors"

	uuid2 "github.com/satori/go.uuid"

	"github.com/talkiewalkie/talkiewalkie/models"
)

var UserConversationMultiCacheByIntErrNotFound = errors.New("UserConversationMultiCacheByInt error did not find values for keys")

type UserConversationMultiCacheByInt struct {
	cache   map[int][]*models.UserConversation
	fetcher func([]int) ([]*models.UserConversation, error)
	primer  func(value *models.UserConversation) int
}

func NewUserConversationMultiCacheByInt(
	fetcher func([]int) ([]*models.UserConversation, error),
	primer func(value *models.UserConversation) int,
) UserConversationMultiCacheByInt {
	return UserConversationMultiCacheByInt{
		cache:   map[int][]*models.UserConversation{},
		fetcher: fetcher,
		primer:  primer,
	}
}

func (cache *UserConversationMultiCacheByInt) Get(identifiers []int) ([][]*models.UserConversation, error) {
	out := make([][]*models.UserConversation, len(identifiers))
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

		for _, record := range records {
			key := cache.primer(record)

			out[key2index[key]] = append(out[key2index[key]], record)
		}

		for key, index := range key2index {
			cache.cache[key] = out[index]
		}
	}

	return out, nil
}

func (cache UserConversationMultiCacheByInt) Clear() {
	cache.cache = nil
}

func (cache UserConversationMultiCacheByInt) Prime(values ...[]*models.UserConversation) {
	for _, value := range values {
		if len(value) > 0 {
			cache.cache[cache.primer(value[0])] = value
		}
	}
}

var (
	_ models.User
	_ uuid2.UUID
)
