package testutils

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
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

func testDbUrl() string {
	return common.DbUrl("talkiewalkie-test", "theo", os.Getenv("TEST_DB_PASSWORD"), "localhost", "5432", false)
}

func SetupDb() *sqlx.DB {
	dbUrl := testDbUrl()
	db := sqlx.MustConnect("postgres", dbUrl)
	common.RunMigrations("../migrations", dbUrl)
	return db
}

func TearDownDb(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec(`DELETE FROM "message";`)
	tx.MustExec(`DELETE FROM "user_group";`)
	tx.MustExec(`DELETE FROM "group";`)
	tx.MustExec(`DELETE FROM "walk";`)
	tx.MustExec(`DELETE FROM "user_walk";`)
	tx.MustExec(`DELETE FROM "asset";`)
	tx.MustExec(`DELETE FROM "user";`)
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

func AddFakeComponentsToRequest(r *http.Request, u *models.User, db *sqlx.DB) *http.Request {
	pgps := common.NewPgPubSub(db, testDbUrl())

	twCtx := common.Context{
		Components: &common.Components{
			Db:       db,
			PgPubSub: &pgps,
			Storage:  FakeStorageClient{},
			CompressImg: func(s string, i int) (string, error) {
				f, _ := ioutil.TempFile("", "")
				return f.Name(), nil
			}},
		User: u,
	}

	return r.WithContext(context.WithValue(r.Context(), "context", twCtx))
}

func MakeRequest(u *models.User, db *sqlx.DB, method, target string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, target, body)

	return AddFakeComponentsToRequest(req, u, db)
}

func RequireReceive(t *testing.T, ws *websocket.Conn, timeout time.Duration, n int, checker func(m []byte) bool) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	msgs := make(chan []byte)
	go func() {
		for {
			mt, m, err := ws.ReadMessage()
			if err != nil {
				log.Printf("(testing) error reading message: (%d) %+v", mt, err)
			}
			msgs <- m
		}
	}()
	cnt := 0
TIMEOUT:
	for {
		select {
		case m := <-msgs:
			if checker(m) {
				cnt += 1
			} else {
				log.Printf("DEBUG:did not pass checker function: '%s'", string(m))
			}
		case <-ctx.Done():
			break TIMEOUT
		}
	}
	if cnt != n {
		t.Fatal(fmt.Errorf("expected to receive message conforming checker %d times but matched %d times", n, cnt))
	}
}

func AddMockUser(db common.DBLogger, t *testing.T) *models.User {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rnd := fmt.Sprintf("%f", rng.Float32())
	rndId := rnd[len(rnd)-6:]
	u := &models.User{
		Handle: fmt.Sprintf("test-user-%s", rndId),
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

func AddMockGroup(db common.DBLogger, t *testing.T, users ...*models.User) *models.Group {
	g := &models.Group{}
	if err := g.Insert(context.Background(), db, boil.Infer()); err != nil {
		t.Log(err)
		t.Fail()
	}
	for _, user := range users {
		ug := models.UserGroup{GroupID: g.ID, UserID: user.ID}
		if err := ug.Insert(context.Background(), db, boil.Infer()); err != nil {
			t.Log(err)
			t.Fail()
		}
	}
	return g
}
