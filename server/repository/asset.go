package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/satori/go.uuid"
)

type Asset struct {
	Id         int       `db:"id" json:"-"`
	Uuid       string    `db:"uuid" json:"uuid"`
	FileName   string    `db:"file_name" json:"fileName"`
	MimeType   string    `db:"mime_type" json:"mimeType"`
	UploadedAt time.Time `db:"uploaded_at" json:"uploadedAt"`
}

type AssetRepository interface {
	GetByUuid(uuid.UUID) (*Asset, error)
	GetAllByUuid(uuids []uuid.UUID) ([]Asset, error)
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

func (p PgAssetRepository) GetAllByUuid(uuids []uuid.UUID) ([]Asset, error) {
	assets := []Asset{}
	uuidsStr := []string{}
	for _, u := range uuids {
		uuidsStr = append(uuidsStr, u.String())
	}
	if err := p.Db.Select(&assets, `SELECT * FROM "asset" WHERE "uuid" = ANY($1);`, pq.Array(uuidsStr)); err != nil {
		return nil, err
	}
	return assets, nil
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
