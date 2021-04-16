package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types/pgeo"

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
	Db  common.DBLogger
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

func (p PgWalkRepository) GetInRadius(pt pgeo.Point, r float64) ([]*models.Walk, error) {
	qs := common.SqlFmt(`
SELECT * FROM "walk" 
WHERE earth_distance(ll_to_earth(walk.start_point[0], walk.start_point[1]),  ll_to_earth($1, $2)) < $3
`)
	walks := []*models.Walk{}

	err := sqlx.Select(p.Db, &walks, qs, pt.X, pt.Y, r)
	if err != nil {
		return nil, err
	}
	return walks, nil
}
