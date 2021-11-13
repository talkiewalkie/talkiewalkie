package repositories

import (
	"context"
	"github.com/jmoiron/sqlx"
	uuid2 "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/repositories/caches"
)

type AssetRepository interface {
	ByIds(...int) ([]*models.Asset, error)
	ByUuids(...uuid2.UUID) ([]*models.Asset, error)

	ById(int) (*models.Asset, error)
	ByUuid(uuid2.UUID) (*models.Asset, error)
}

type AssetRepositoryImpl struct {
	Db        *sqlx.DB
	Context   context.Context
	IdCache   caches.AssetCacheByInt
	UuidCache caches.AssetCacheByUuid
}

func NewAssetRepository(context context.Context, db *sqlx.DB) *AssetRepositoryImpl {
	return &AssetRepositoryImpl{
		Db:      db,
		Context: context,
		IdCache: caches.NewAssetCacheByInt(func(ints []int) ([]*models.Asset, error) {
			return models.Assets(models.ConversationWhere.ID.IN(ints)).All(context, db)
		}, func(value *models.Asset) int {
			return value.ID
		}),
		UuidCache: caches.NewAssetCacheByUuid(func(uuids []uuid2.UUID) ([]*models.Asset, error) {
			return models.Assets(models.AssetWhere.UUID.IN(uuids)).All(context, db)
		}, func(value *models.Asset) uuid2.UUID {
			return value.UUID
		}),
	}
}

func (repository AssetRepositoryImpl) ByIds(ints ...int) ([]*models.Asset, error) {
	assets, err := repository.IdCache.Get(ints)
	if err != nil {
		return nil, err
	}

	repository.UuidCache.Prime(assets...)
	return assets, nil
}

func (repository AssetRepositoryImpl) ByUuids(uuids ...uuid2.UUID) ([]*models.Asset, error) {
	assets, err := repository.UuidCache.Get(uuids)
	if err != nil {
		return nil, err
	}

	repository.IdCache.Prime(assets...)
	return assets, nil
}

func (repository AssetRepositoryImpl) ById(id int) (*models.Asset, error) {
	assets, err := repository.ByIds(id)
	if err != nil {
		return nil, err
	}

	return assets[0], nil
}

func (repository AssetRepositoryImpl) ByUuid(uuid uuid2.UUID) (*models.Asset, error) {
	assets, err := repository.ByUuids(uuid)
	if err != nil {
		return nil, err
	}

	return assets[0], nil
}

var _ AssetRepository = AssetRepositoryImpl{}
