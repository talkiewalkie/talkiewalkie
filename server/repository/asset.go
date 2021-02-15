package repository

import (
	"github.com/docker/distribution/uuid"
	"github.com/jmoiron/sqlx"
)

type Asset struct {
	Id       int    `json:"-" db:"id"`
	Uuid     string `json:"uuid" db:"uuid"`
	FileName string `json:"fileName" db:"file_name"`
	MimeType string `json:"mimeType" db:"mime_type"`
	Url      string `json:"url" db:"-"`
}

type AssetRepository interface {
	GetByUuid(uid uuid.UUID) (*Asset, error)
	Create(uid uuid.UUID, fileName, mimeType string) (*Asset, error)
}

var _ AssetRepository = PgAssetRepository{}

type PgAssetRepository struct {
	Db *sqlx.DB
}

func (p PgAssetRepository) GetByUuid(uid uuid.UUID) (*Asset, error) {
	var a Asset
	if err := p.Db.QueryRowx(
		`SELECT * FROM "asset" WHERE uuid = $1;`, uid.String(),
	).StructScan(&a); err != nil {
		return nil, err
	}
	return &a, nil
}

func (p PgAssetRepository) Create(uid uuid.UUID, fileName, mimeType string) (*Asset, error) {
	var a Asset
	if err := p.Db.QueryRowx(
		`INSERT INTO "asset" (uuid, file_name, mime_type) VALUES ($1, $2, $3) RETURNING *`,
		uid.String(), fileName, mimeType,
	).StructScan(&a); err != nil {
		return nil, err
	}
	return &a, nil
}
