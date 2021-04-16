package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types/pgeo"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/repository"
	"github.com/talkiewalkie/talkiewalkie/testutils"
)

// https://stackoverflow.com/a/40748810
func TestWalkRepository(t *testing.T) {
	db := testutils.SetupDb()
	ql := common.NewDbLogger(db)
	repo := repository.PgWalkRepository{
		Components: &common.Components{},
		Db:         ql,
		Ctx:        context.Background(),
	}

	t.Run("can create walk", createWalkTest(repo))
	testutils.TearDownDb(db)
	t.Run("can list walk", listWalkTest(repo))
	testutils.TearDownDb(db)
	t.Run("can list walk near point", listWalkInRadiusTest(repo))
	testutils.TearDownDb(db)

}

var (
	HOME            = pgeo.Point{X: 48.8450234, Y: 2.3997529}
	IPPUDO          = pgeo.Point{X: 48.8645814, Y: 2.3425034} // yumyum tasty ramens
	REF_DISTANCE_KM = 4.73
)

func createWalkTest(repo repository.PgWalkRepository) func(t *testing.T) {
	return func(t *testing.T) {
		u := testutils.AddMockUser(repo.Db, t)
		w := &models.Walk{
			Title:      "some title",
			CoverID:    null.Int{Valid: false},
			AudioID:    null.Int{Valid: false},
			AuthorID:   u.ID,
			StartPoint: IPPUDO,
			EndPoint:   IPPUDO,
		}
		err := w.Insert(context.Background(), repo.Db, boil.Infer())
		if err != nil {
			t.Log(err)
			t.Fail()
		}
	}
}

func listWalkTest(repo repository.PgWalkRepository) func(t *testing.T) {
	return func(t *testing.T) {
		u := testutils.AddMockUser(repo.Db, t)
		_ = testutils.AddMockWalk(u, repo.Db, t)
		walks, err := repo.GetAll()
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		assert.Equal(t, 1, len(walks))

		near, err := repo.GetInRadius(HOME, REF_DISTANCE_KM*1000+500)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		assert.Equal(t, 1, len(near))
	}
}
func listWalkInRadiusTest(repo repository.PgWalkRepository) func(t *testing.T) {
	return func(t *testing.T) {
		u := testutils.AddMockUser(repo.Db, t)
		_ = testutils.AddMockWalk(u, repo.Db, t)

		near, err := repo.GetInRadius(HOME, REF_DISTANCE_KM*1000+500)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		assert.Equal(t, 1, len(near))
	}
}
