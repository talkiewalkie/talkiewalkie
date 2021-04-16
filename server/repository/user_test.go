package repository_test

import (
	"context"
	"testing"

	_ "github.com/lib/pq"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/repository"
	"github.com/talkiewalkie/talkiewalkie/testutils"
)

func TestUserRepository(t *testing.T) {
	db := testutils.SetupDb()
	ql := common.NewDbLogger(db)
	repo := repository.PgUserRepository{
		Components: &common.Components{},
		Db:         ql,
		Ctx:        context.Background(),
	}

	t.Run("can create user", createUserTest(repo))
	testutils.TearDownDb(db)
}

func createUserTest(repo repository.PgUserRepository) func(t *testing.T) {
	return func(t *testing.T) {
		_, err := repo.CreateUser("bibou", "bibou.doux@doudou.dou", "sometoken", []byte("kd"))
		if err != nil {
			t.Log(err)
			t.Fail()
		}
	}
}
