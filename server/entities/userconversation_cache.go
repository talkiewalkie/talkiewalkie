// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package entities

import (
	"errors"
	"log"
)

type UserConversationMultiCacheByInt struct {
	cache   map[int][]*UserConversation
	fetcher func([]int) ([]*UserConversation, error)
	primer  func(value *UserConversation) int
}

func NewUserConversationMultiCacheByInt(
	fetcher func([]int) ([]*UserConversation, error),
	primer func(value *UserConversation) int,
) UserConversationMultiCacheByInt {
	return UserConversationMultiCacheByInt{
		cache:   map[int][]*UserConversation{},
		fetcher: fetcher,
		primer:  primer,
	}
}

func (cache *UserConversationMultiCacheByInt) Get(identifiers []int) ([][]*UserConversation, error) {
	out := make([][]*UserConversation, len(identifiers))
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

	for _, value := range out {
		if value == nil {
			return nil, errors.New("could not fetch from : found nil value")
		}
	}

	return out, nil
}

func (cache UserConversationMultiCacheByInt) Clear() {
	cache.cache = nil
}

func (cache UserConversationMultiCacheByInt) Prime(values ...[]*UserConversation) {
	for _, value := range values {
		if len(value) > 0 {
			cache.cache[cache.primer(value[0])] = value
		} else {
			log.Printf("error?: sending empty lists to prime the cache")
		}
	}
}
