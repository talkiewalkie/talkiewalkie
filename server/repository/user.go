package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type User struct {
	Id          int            `json:"-"`
	Uuid        string         `json:"uuid"`
	Handle      string         `json:"handle"`
	Email       string         `json:"email"`
	FirebaseUid string         `json:"-" db:"firebase_uid"`
	EmailToken  sql.NullString `json:"-" db:"email_token"`
}

type UserRepository interface {
	GetUserByHandle(handle string) (*User, error)
	GetUserByUid(uid string) (*User, error)
	CreateUser(email, firebaseUid, emailToken string) (*User, error)
}

var _ UserRepository = PgUserRepository{}

type PgUserRepository struct {
	Db *sqlx.DB
}

func (p PgUserRepository) GetUserByHandle(handle string) (*User, error) {
	var u User
	if err := p.Db.QueryRowx(`SELECT * FROM "user" WHERE handle = $1;`, handle).StructScan(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (p PgUserRepository) GetUserByUid(uid string) (*User, error) {
	var u User
	if err := p.Db.QueryRowx(`SELECT * FROM "user" WHERE "firebase_uid" = $1;`, uid).StructScan(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (p PgUserRepository) CreateUser(email, firebaseUid, emailToken string) (*User, error) {
	var u User
	err := p.Db.QueryRowx(`
	INSERT INTO "user" 
	    (handle, email, firebase_uid, email_token) VALUES 
		(?, ?, ?, ?) RETURNING (id, uuid, handle, email)`,
		email, email, firebaseUid, emailToken,
	).StructScan(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
