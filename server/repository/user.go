package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type User struct {
	Id         int            `json:"-" db:"id"`
	Uuid       string         `json:"uuid" db:"uuid"`
	Handle     string         `json:"handle" db:"handle"`
	Email      string         `json:"email" db:"email"`
	Password   []byte         `json:"-" db:"password"`
	EmailToken sql.NullString `json:"-" db:"email_token"`
}

type UserRepository interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByUuid(uid string) (*User, error)
	CreateUser(email string, password []byte, emailToken string) (*User, error)
}

var _ UserRepository = PgUserRepository{}

type PgUserRepository struct {
	Db *sqlx.DB
}

func (p PgUserRepository) GetUserByEmail(email string) (*User, error) {
	var u User
	if err := p.Db.QueryRowx(`SELECT * FROM "user" WHERE email = $1;`, email).StructScan(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (p PgUserRepository) GetUserByUuid(uid string) (*User, error) {
	var u User
	if err := p.Db.QueryRowx(`SELECT * FROM "user" WHERE "uuid" = $1;`, uid).StructScan(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (p PgUserRepository) CreateUser(email string, password []byte, emailToken string) (*User, error) {
	var u User
	err := p.Db.QueryRowx(`
		INSERT INTO "user" (handle, email, password, email_token) 
		VALUES ($1, $2, $3, $4) RETURNING *;
		`, email, email, password, emailToken,
	).StructScan(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
