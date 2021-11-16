package generics

import (
	"errors"
	uuid2 "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/models"
)

var CacheValueMultiCacheByCacheKeyErrNotFound = errors.New("CacheValueMultiCacheByCacheKey error did not find values for keys")

type CacheValueMultiCacheByCacheKey struct {
	cache   map[CacheKey][]*CacheValue
	fetcher func([]CacheKey) ([]*CacheValue, error)
	primer  func(value *CacheValue) CacheKey
}

func NewCacheValueMultiCacheByCacheKey(
	fetcher func([]CacheKey) ([]*CacheValue, error),
	primer func(value *CacheValue) CacheKey,
) CacheValueMultiCacheByCacheKey {
	return CacheValueMultiCacheByCacheKey{
		cache:   map[CacheKey][]*CacheValue{},
		fetcher: fetcher,
		primer:  primer,
	}
}

func (cache *CacheValueMultiCacheByCacheKey) Get(identifiers []CacheKey) ([][]*CacheValue, error) {
	out := make([][]*CacheValue, len(identifiers))
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

func (cache CacheValueMultiCacheByCacheKey) Clear() {
	cache.cache = nil
}

func (cache CacheValueMultiCacheByCacheKey) Prime(values ...[]*CacheValue) {
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
