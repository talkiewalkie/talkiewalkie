package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/testutils"
)

func TestWalkRepository(t *testing.T) {
	db := testutils.SetupDb()

	t.Run("can create walk", createWalkTest(db))
	testutils.TearDownDb(db)
	//t.Run("can list walk", listWalkTest(repo))
	//testutils.TearDownDb(db)
	//t.Run("can list walk near point", listWalkInRadiusTest(repo))
	//testutils.TearDownDb(db)
}

func createWalkTest(db *sqlx.DB) func(t *testing.T) {
	return func(t *testing.T) {
		u := testutils.AddMockUser(db, t)

		w := &httptest.ResponseRecorder{}
		bb, _ := json.Marshal(CreateWalkInput{
			Title:        "test walk",
			Description:  "",
			CoverArtUuid: uuid.UUID{},
			AudioUuid:    uuid.UUID{},
		})
		r := httptest.NewRequest(http.MethodGet, "/walk", bytes.NewReader(bb))
		mctx := common.Context{
			Components: &common.Components{Db: db},
			User:       u,
		}
		ctx := context.WithValue(r.Context(), "context", mctx)
		r = r.WithContext(ctx)

		CreateWalk(w, r)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		//w := &models.Walk{
		//	Title:      "some title",
		//	CoverID:    null.Int{Valid: false},
		//	AudioID:    null.Int{Valid: false},
		//	AuthorID:   u.ID,
		//	StartPoint: IPPUDO,
		//	EndPoint:   IPPUDO,
		//}
		//err := w.Insert(context.Background(), repo.Db, boil.Infer())
		//if err != nil {
		//	t.Log(err)
		//	t.Fail()
		//}
	}
}
