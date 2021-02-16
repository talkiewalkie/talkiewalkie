package repository

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type Walk struct {
	Id        int           `db:"id" json:"-"`
	Uuid      string        `db:"uuid" json:"uuid"`
	Title     string        `db:"title" json:"title"`
	CoverId   sql.NullInt32 `db:"cover_id" json:"-"`
	AudioId   sql.NullInt32 `db:"audio_id" json:"-"`
	CreatedAt time.Time     `db:"created_at" json:"createdAt"`
	AuthorId  int           `db:"author_id" json:"-"`
}

type WalkRepository interface {
	GetAll() ([]*Walk, error)
}

var _ WalkRepository = PgWalkRepository{}

type PgWalkRepository struct {
	Db *sqlx.DB
}

func (p PgWalkRepository) GetAll() ([]*Walk, error) {
	walks := []*Walk{}
	rows, err := p.Db.Queryx(`SELECT * FROM "walk" LIMIT 50;`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var w Walk
		err = rows.StructScan(&w)
		if err != nil {
			return nil, err
		}
		walks = append(walks, &w)
	}
	return walks, nil
}
