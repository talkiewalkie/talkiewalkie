package testutils

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types/pgeo"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
)

func SetupDb() *sqlx.DB {
	dbUrl := common.DbUrl("talkiewalkie-test", "theo", os.Getenv("TEST_DB_PASSWORD"), "localhost", "5432", false)
	db := sqlx.MustConnect("postgres", dbUrl)
	common.RunMigrations("../migrations", dbUrl)
	return db
}

func TearDownDb(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec(`DELETE FROM "user";`)
	tx.MustExec(`DELETE FROM "walk";`)
	tx.MustExec(`DELETE FROM "user_walk";`)
	tx.MustExec(`DELETE FROM "asset";`)
	if err := tx.Commit(); err != nil {
		log.Panicf("could not reset db state: %+v", err)
	}
}

type FakeStorageClient struct{}

func (f FakeStorageClient) AssetUrl(asset *models.Asset) (string, error) {
	return "https://some.fake.url/123", nil
}

func (f FakeStorageClient) DefaultBucket() string {
	return "test-bucket"
}

func (f FakeStorageClient) Download(blobName string, writer io.Writer) error {
	_, err := writer.Write([]byte("hello this is test content"))
	return err
}

func (f FakeStorageClient) Upload(ctx context.Context, blob io.Reader) (*uuid.UUID, error) {
	uid := uuid.NewV4()
	return &uid, nil
}

func (f FakeStorageClient) SignedUrl(bucket, blobName string) (string, error) {
	return "https://some.fake.url/123", nil
}

var _ common.StorageClient = FakeStorageClient{}

func MakeRequest(u *models.User, db *sqlx.DB, method, target string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, target, body)

	twCtx := common.Context{
		Components: &common.Components{
			Db:      db,
			Storage: FakeStorageClient{},
			CompressImg: func(s string, i int) (string, error) {
				f, _ := ioutil.TempFile("", "")
				return f.Name(), nil
			}},
		User: u,
	}
	ctx := context.WithValue(req.Context(), "context", twCtx)

	return req.WithContext(ctx)
}

func AddMockUser(db common.DBLogger, t *testing.T) *models.User {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rnd := fmt.Sprintf("%f", rng.Float32())
	rndId := rnd[len(rnd)-6:]
	u := &models.User{
		Handle:     fmt.Sprintf("test-user-%s", rndId),
		Email:      fmt.Sprintf("test-user-%s@gmail.com", rndId),
		Password:   []byte("abc123"),
		EmailToken: null.String{},
	}
	if err := u.Insert(context.Background(), db, boil.Infer()); err != nil {
		t.Log(err)
		t.Fail()
	}
	return u
}

var (
	IPPUDO = pgeo.Point{X: 48.8645814, Y: 2.3425034} // yumyum tasty ramens
)

func AddMockWalk(u *models.User, db common.DBLogger, t *testing.T) *models.Walk {
	w := &models.Walk{
		Title:      "some title",
		CoverID:    null.Int{Valid: false},
		AudioID:    null.Int{Valid: false},
		AuthorID:   u.ID,
		StartPoint: IPPUDO,
	}
	if err := w.Insert(context.Background(), db, boil.Infer()); err != nil {
		t.Log(err)
		t.Fail()
	}
	return w
}

func AddMockAsset(mimeType string, db common.DBLogger, t *testing.T) *models.Asset {
	a := &models.Asset{MimeType: mimeType}
	if err := a.Insert(context.Background(), db, boil.Infer()); err != nil {
		t.Log(err)
		t.Fail()
	}
	return a
}
