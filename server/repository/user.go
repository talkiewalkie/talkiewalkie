package repository

import (
	"context"

	uuid "github.com/satori/go.uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
)

type UserRepository interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUuid(uuid.UUID) (*models.User, error)
	CreateUser(handle, email, emailToken string, password []byte) (*models.User, error)
}

var _ UserRepository = PgUserRepository{}

type PgUserRepository struct {
	*common.Components
	Db  common.DBLogger
	Ctx context.Context
}

func (p PgUserRepository) GetUserByEmail(email string) (*models.User, error) {
	u, err := models.Users(models.UserWhere.Email.EQ(email)).One(p.Ctx, p.Db)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (p PgUserRepository) GetUserByUuid(uid uuid.UUID) (*models.User, error) {
	u, err := models.Users(models.UserWhere.UUID.EQ(uid)).One(p.Ctx, p.Db)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (p PgUserRepository) CreateUser(handle, email, emailToken string, password []byte) (*models.User, error) {
	u := &models.User{Handle: handle, Email: email, Password: password, EmailToken: null.NewString(emailToken, true)}
	if err := u.Insert(p.Ctx, p.Db, boil.Infer()); err != nil {
		return nil, err
	}
	return u, nil
}
