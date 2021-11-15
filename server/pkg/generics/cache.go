package generics

import (
	"errors"
	uuid2 "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/models"

	"github.com/cheekybits/genny/generic"
)

type CacheKey generic.Type
type CacheValue generic.Type

type CacheValueCacheByCacheKey struct {
	cache   map[CacheKey]*CacheValue
	fetcher func([]CacheKey) ([]*CacheValue, error)
	primer  func(value *CacheValue) CacheKey
}

func NewCacheValueCacheByCacheKey(
	fetcher func([]CacheKey) ([]*CacheValue, error),
	primer func(value *CacheValue) CacheKey,
) CacheValueCacheByCacheKey {
	return CacheValueCacheByCacheKey{
		cache:   map[CacheKey]*CacheValue{},
		fetcher: fetcher,
		primer:  primer,
	}
}

func (cache *CacheValueCacheByCacheKey) Get(
	identifiers []CacheKey,
) ([]*CacheValue, error) {
	out := make([]*CacheValue, len(identifiers))
	key2index := map[CacheKey]int{}

	for idx, key := range identifiers {
		item, ok := cache.cache[key]
		if ok {
			out[idx] = item
		} else {
			key2index[key] = idx
		}
	}

	if len(key2index) > 0 {
		missingKeys := make([]CacheKey, len(key2index))
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

func (cache CacheValueCacheByCacheKey) Clear() {
	cache.cache = nil
}

func (cache CacheValueCacheByCacheKey) Prime(values ...*CacheValue) {
	for _, value := range values {
		cache.cache[cache.primer(value)] = value
	}
}

// ENSURE IMPORTS
var (
	_ uuid2.UUID
	_ models.User
)
