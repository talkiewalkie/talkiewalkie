package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
)

type AssetRepository interface {
	GetByUuid(uuid.UUID) (*models.Asset, error)
	GetAllByUuid(uuids []uuid.UUID) ([]*models.Asset, error)
	Create(uid uuid.UUID, fileName, mimeType string) (*models.Asset, error)
}

var _ AssetRepository = PgAssetRepository{}

type PgAssetRepository struct {
	*common.Components
	Db  *sqlx.DB
	Ctx context.Context
}

func (p PgAssetRepository) GetByUuid(uid uuid.UUID) (*models.Asset, error) {
	a, err := models.Assets(models.AssetWhere.UUID.EQ(uid)).One(p.Ctx, p.Db)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (p PgAssetRepository) GetAllByUuid(uuids []uuid.UUID) ([]*models.Asset, error) {
	a, err := models.Assets(models.AssetWhere.UUID.IN(uuids)).All(p.Ctx, p.Db)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (p PgAssetRepository) Create(uid uuid.UUID, fileName, mimeType string) (*models.Asset, error) {
	a := models.Asset{UUID: uid, FileName: fileName, MimeType: mimeType}
	if err := a.Insert(p.Ctx, p.Db, boil.Infer()); err != nil {
		return nil, err
	}
	return &a, nil
}
