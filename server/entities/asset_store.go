package entities

import (
	"context"
	uuid2 "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
)

type AssetStore interface {
	ByIds(...int) ([]*Asset, error)
	ByUuids(...uuid2.UUID) ([]*Asset, error)
}

type AssetStoreImpl struct {
	*common.Components
	Context   context.Context
	IdCache   AssetCacheByInt
	UuidCache AssetCacheByUuid2UUID
}

func NewAssetStore(context context.Context, components *common.Components) *AssetStoreImpl {
	return &AssetStoreImpl{
		Components: components,
		Context:    context,
		IdCache: NewAssetCacheByInt(func(ints []int) ([]*Asset, error) {
			out := []*Asset{}
			records, err := models.Assets(models.ConversationWhere.ID.IN(ints)).All(context, components.Db)
			if err != nil {
				return nil, err
			}

			for _, record := range records {
				out = append(out, &Asset{record, components})
			}

			return out, nil
		}, func(value *Asset) int {
			return value.Record.ID
		}),
		UuidCache: NewAssetCacheByUuid2UUID(func(uuids []uuid2.UUID) ([]*Asset, error) {
			out := []*Asset{}
			records, err := models.Assets(models.AssetWhere.UUID.IN(uuids)).All(context, components.Db)
			if err != nil {
				return nil, err
			}

			for _, record := range records {
				out = append(out, &Asset{record, components})
			}

			return out, nil
		}, func(value *Asset) uuid2.UUID {
			return value.Record.UUID
		}),
	}
}

func (store AssetStoreImpl) ByIds(ints ...int) ([]*Asset, error) {
	assets, err := store.IdCache.Get(ints)
	if err != nil {
		return nil, err
	}

	store.UuidCache.Prime(assets...)
	return assets, nil
}

func (store AssetStoreImpl) ByUuids(uuids ...uuid2.UUID) ([]*Asset, error) {
	assets, err := store.UuidCache.Get(uuids)
	if err != nil {
		return nil, err
	}

	store.IdCache.Prime(assets...)
	return assets, nil
}

var _ AssetStore = AssetStoreImpl{}
