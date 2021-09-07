package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"github.com/talkiewalkie/talkiewalkie/models"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/testutils"
)

func TestWalkRepository(t *testing.T) {
	db := testutils.SetupDb()

	testutils.TearDownDb(db)
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
			Title:       "test walk",
			Description: "test walk for testing purposes",
			StartPoint:  Coords{Lat: testutils.IPPUDO.X, Lng: testutils.IPPUDO.Y},
		})
		var formData bytes.Buffer
		formWriter := multipart.NewWriter(&formData)
		_ = formWriter.WriteField("payload", string(bb))
		cf, _ := formWriter.CreateFormFile("cover", "cover.test.png")
		cf.Write([]byte("fake file"))
		af, _ := formWriter.CreateFormFile("walk", "audio.test.mp3")
		af.Write([]byte("fake file"))
		_ = formWriter.Close()

		r := httptest.NewRequest(http.MethodGet, "/walk", &formData)
		r.Header.Set("Content-Type", formWriter.FormDataContentType())
		mctx := common.Context{
			Components: &common.Components{Db: db, Storage: testutils.FakeStorageClient{}},
			User:       u,
		}
		ctx := context.WithValue(r.Context(), "context", mctx)
		r = r.WithContext(ctx)

		CreateWalk(w, r)
		require.Equal(t, http.StatusOK, w.Code, "hello", w.Body.String())
		cnt, err := models.Walks().Count(ctx, mctx.Components.Db)
		require.Equal(t, int64(1), cnt)
		if err != nil {
			t.Fatalf("could not count walks: %+v", err)
		}

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
