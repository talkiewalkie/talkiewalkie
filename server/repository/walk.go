package repository

import (
	"context"

	"github.com/cridenour/go-postgis"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
)

type WalkRepository interface {
	GetAll() ([]*models.Walk, error)
	GetByUuid(uuid.UUID) (*models.Walk, error)
}

var _ WalkRepository = PgWalkRepository{}

type PgWalkRepository struct {
	*common.Components
	Db  *sqlx.DB
	Ctx context.Context
}

func (p PgWalkRepository) GetAll() ([]*models.Walk, error) {
	walks, err := models.Walks(qm.Limit(50)).All(p.Ctx, p.Db)
	if err != nil {
		return nil, err
	}
	return walks, nil
}

func (p PgWalkRepository) GetByUuid(uid uuid.UUID) (*models.Walk, error) {
	w, err := models.Walks(models.WalkWhere.UUID.EQ(uid)).One(p.Ctx, p.Db)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (p PgWalkRepository) GetInRadius(pt postgis.Point, r float32) ([]*models.Walk, error) {
	_, err := p.Db.Query(`SELECT * FROM walk WHERE st_distance(start_at, $1) < $2`, pt, r)
	return nil, err
}
