package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/talkiewalkie/talkiewalkie/models"
)

type WalkRepository interface {
	GetAll() ([]*models.Walk, error)
}

var _ WalkRepository = PgWalkRepository{}

type PgWalkRepository struct {
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
