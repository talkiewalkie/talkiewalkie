package testutils

import (
	"context"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/gorilla/websocket"
	"github.com/talkiewalkie/talkiewalkie/clients"
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
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types/pgeo"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
)

func testDbUrl() string {
	return common.DbUri("talkiewalkie-test", "theo", os.Getenv("TEST_DB_PASSWORD"), "localhost", "5432", false)
}

func SetupDb() *sqlx.DB {
	dbUrl := testDbUrl()
	db := sqlx.MustConnect("postgres", dbUrl)
	common.RunMigrations("../migrations", dbUrl)
	return db
}

func TearDownDb(ctx context.Context, db *sqlx.DB) {
	tx := db.MustBegin()
	if _, err := models.Messages().DeleteAll(ctx, tx); err != nil {
		panic(err)
	}
	if _, err := models.UserConversations().DeleteAll(ctx, tx); err != nil {
		panic(err)
	}
	if _, err := models.Conversations().DeleteAll(ctx, tx); err != nil {
		panic(err)
	}
	if _, err := models.Users().DeleteAll(ctx, tx); err != nil {
		panic(err)
	}
	if _, err := models.Events().DeleteAll(ctx, tx); err != nil {
		panic(err)
	}
	if _, err := models.Assets().DeleteAll(ctx, tx); err != nil {
		panic(err)
	}
	if err := tx.Commit(); err != nil {
		log.Panicf("could not reset db state: %+v", err)
	}
}

func NewContext(db *sqlx.DB, t *testing.T) (*common.Components, *models.User, context.Context) {
	me := AddMockUser(db, t)
	components := NewFakeComponents(db)
	ctx := context.WithValue(context.Background(), "components", components)
	ctx = context.WithValue(ctx, "me", me)

	components.ResetEntityStores(ctx)
	return components, me, ctx
}

func NewFakeComponents(db *sqlx.DB) *common.Components {
	pubsubClient := clients.NewPgPubSub(db, testDbUrl())

	return &common.Components{
		Db: db,

		AuthClient:      nil,
		MessagingClient: nil,
		PubSubClient:    pubsubClient,
		StorageClient:   FakeStorageClient{},
		AudioClient:     nil,

		CompressImg: func(s string, i int) (string, error) {
			f, _ := ioutil.TempFile("", "")
			return f.Name(), nil
		}}
}

func AddFakeComponentsToRequest(r *http.Request, u *models.User, db *sqlx.DB) *http.Request {
	twCtx := common.Context{
		Components: NewFakeComponents(db),
		User:       u,
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
				if cnt == n {
					return
				}
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
		PhoneNumber:        faker.Phonenumber(),
		DisplayName:        null.StringFrom(fmt.Sprintf("test-user-%s", rndId)),
		Locales:            []string{"fr"},
		OnboardingFinished: true,
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

func AddMockConversation(db common.DBLogger, t *testing.T, users ...*models.User) *models.Conversation {
	g := &models.Conversation{}
	if err := g.Insert(context.Background(), db, boil.Infer()); err != nil {
		t.Log(err)
		t.Fail()
	}
	for _, user := range users {
		ug := models.UserConversation{ConversationID: g.ID, UserID: user.ID}
		if err := ug.Insert(context.Background(), db, boil.Infer()); err != nil {
			t.Log(err)
			t.Fail()
		}
	}
	return g
}
