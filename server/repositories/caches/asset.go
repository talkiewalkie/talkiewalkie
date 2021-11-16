// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package caches

import (
	uuid2 "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/models"

	"errors"
)

type AssetCacheByInt struct {
	cache   map[int]*models.Asset
	fetcher func([]int) ([]*models.Asset, error)
	primer  func(value *models.Asset) int
}

func NewAssetCacheByInt(
	fetcher func([]int) ([]*models.Asset, error),
	primer func(value *models.Asset) int,
) AssetCacheByInt {
	return AssetCacheByInt{
		cache:   map[int]*models.Asset{},
		fetcher: fetcher,
		primer:  primer,
	}
}

func (cache *AssetCacheByInt) Get(
	identifiers []int,
) ([]*models.Asset, error) {
	out := make([]*models.Asset, len(identifiers))
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

func (cache AssetCacheByInt) Clear() {
	cache.cache = nil
}

func (cache AssetCacheByInt) Prime(values ...*models.Asset) {
	for _, value := range values {
		cache.cache[cache.primer(value)] = value
	}
}

// ENSURE IMPORTS
var (
	_ uuid2.UUID
	_ models.User
)

type AssetCacheByUuid struct {
	cache   map[uuid2.UUID]*models.Asset
	fetcher func([]uuid2.UUID) ([]*models.Asset, error)
	primer  func(value *models.Asset) uuid2.UUID
}

func NewAssetCacheByUuid(
	fetcher func([]uuid2.UUID) ([]*models.Asset, error),
	primer func(value *models.Asset) uuid2.UUID,
) AssetCacheByUuid {
	return AssetCacheByUuid{
		cache:   map[uuid2.UUID]*models.Asset{},
		fetcher: fetcher,
		primer:  primer,
	}
}

func (cache *AssetCacheByUuid) Get(
	identifiers []uuid2.UUID,
) ([]*models.Asset, error) {
	out := make([]*models.Asset, len(identifiers))
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

func (cache AssetCacheByUuid) Clear() {
	cache.cache = nil
}

func (cache AssetCacheByUuid) Prime(values ...*models.Asset) {
	for _, value := range values {
		cache.cache[cache.primer(value)] = value
	}
}

// ENSURE IMPORTS
var (
	_ uuid2.UUID
	_ models.User
)
