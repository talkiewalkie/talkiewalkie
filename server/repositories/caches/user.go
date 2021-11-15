// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package caches

import (
	uuid2 "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/models"

	"errors"
)

type UserCacheByInt struct {
	cache   map[int]*models.User
	fetcher func([]int) ([]*models.User, error)
	primer  func(value *models.User) int
}

func NewUserCacheByInt(
	fetcher func([]int) ([]*models.User, error),
	primer func(value *models.User) int,
) UserCacheByInt {
	return UserCacheByInt{
		cache:   map[int]*models.User{},
		fetcher: fetcher,
		primer:  primer,
	}
}

func (cache *UserCacheByInt) Get(
	identifiers []int,
) ([]*models.User, error) {
	out := make([]*models.User, len(identifiers))
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

	for _, value := range out {
		if value == nil {
			return nil, errors.New("could not fetch from : found nil value")
		}
	}

	return out, nil
}

func (cache UserCacheByInt) Clear() {
	cache.cache = nil
}

func (cache UserCacheByInt) Prime(values ...*models.User) {
	for _, value := range values {
		cache.cache[cache.primer(value)] = value
	}
}

// ENSURE IMPORTS
var (
	_ uuid2.UUID
	_ models.User
)

type UserCacheByUuid struct {
	cache   map[uuid2.UUID]*models.User
	fetcher func([]uuid2.UUID) ([]*models.User, error)
	primer  func(value *models.User) uuid2.UUID
}

func NewUserCacheByUuid(
	fetcher func([]uuid2.UUID) ([]*models.User, error),
	primer func(value *models.User) uuid2.UUID,
) UserCacheByUuid {
	return UserCacheByUuid{
		cache:   map[uuid2.UUID]*models.User{},
		fetcher: fetcher,
		primer:  primer,
	}
}

func (cache *UserCacheByUuid) Get(
	identifiers []uuid2.UUID,
) ([]*models.User, error) {
	out := make([]*models.User, len(identifiers))
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

	for _, value := range out {
		if value == nil {
			return nil, errors.New("could not fetch from : found nil value")
		}
	}

	return out, nil
}

func (cache UserCacheByUuid) Clear() {
	cache.cache = nil
}

func (cache UserCacheByUuid) Prime(values ...*models.User) {
	for _, value := range values {
		cache.cache[cache.primer(value)] = value
	}
}

// ENSURE IMPORTS
var (
	_ uuid2.UUID
	_ models.User
)

type UserCacheByString struct {
	cache   map[string]*models.User
	fetcher func([]string) ([]*models.User, error)
	primer  func(value *models.User) string
}

func NewUserCacheByString(
	fetcher func([]string) ([]*models.User, error),
	primer func(value *models.User) string,
) UserCacheByString {
	return UserCacheByString{
		cache:   map[string]*models.User{},
		fetcher: fetcher,
		primer:  primer,
	}
}

func (cache *UserCacheByString) Get(
	identifiers []string,
) ([]*models.User, error) {
	out := make([]*models.User, len(identifiers))
	key2index := map[string]int{}

	for idx, key := range identifiers {
		item, ok := cache.cache[key]
		if ok {
			out[idx] = item
		} else {
			key2index[key] = idx
		}
	}

	if len(key2index) > 0 {
		missingKeys := make([]string, len(key2index))
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
			return nil, errors.New("could not fetch from : found nil value")
		}
	}

	return out, nil
}

func (cache UserCacheByString) Clear() {
	cache.cache = nil
}

func (cache UserCacheByString) Prime(values ...*models.User) {
	for _, value := range values {
		cache.cache[cache.primer(value)] = value
	}
}

// ENSURE IMPORTS
var (
	_ uuid2.UUID
	_ models.User
)
