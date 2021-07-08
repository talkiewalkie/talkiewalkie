package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/testutils"
)

type FakeStorageClient struct{}

func (f FakeStorageClient) Upload(ctx context.Context, blob io.Reader) (*uuid.UUID, error) {
	uid := uuid.NewV4()
	return &uid, nil
}

func (f FakeStorageClient) Url(dest string) (string, error) {
	return "https://some.fake.url/123", nil
}

var _ common.StorageClient = FakeStorageClient{}

func TestAssets(t *testing.T) {
	db := testutils.SetupDb()

	t.Run("can create assets", createAssetTest(db))
	testutils.TearDownDb(db)
	//t.Run("can list walk", listWalkTest(repo))
	//testutils.TearDownDb(db)
	//t.Run("can list walk near point", listWalkInRadiusTest(repo))
	//testutils.TearDownDb(db)
}

func createAssetTest(db *sqlx.DB) func(t *testing.T) {
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
			Components: &common.Components{Db: db, Storage: FakeStorageClient{}},
			User:       u,
		}
		ctx := context.WithValue(r.Context(), "context", mctx)
		r = r.WithContext(ctx)

		UploadHandler(w, r)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		buf := new(bytes.Buffer)
		buf.ReadFrom(w.Result().Body)
		newStr := buf.String()
		t.Log(w.Result().Status)
		t.Log(newStr)
		t.Log(w.Body)

		assets, _ := models.Assets().All(r.Context(), mctx.Components.Db)
		assert.NotEmptyf(t, assets, "no asset in db")
	}
}
