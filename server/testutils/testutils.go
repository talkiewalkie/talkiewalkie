package testutils

import (
	"context"
	"fmt"
	"github.com/bxcodec/faker/v3"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/talkiewalkie/talkiewalkie/clients"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
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

func TearDownDb(db *sqlx.DB) {
	ctx := context.Background()

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

// NewFakeServer will mock a live server for use when testing streaming calls.
// It uses `bufconn` which allows us to simulate the server through buffers instead of opening ports.
func NewFakeServer(
	db *sqlx.DB,
	t *testing.T,
	withService func(*grpc.Server),
) (
	*common.Components,
	*models.User,
	context.Context,
	*grpc.ClientConn,
	context.CancelFunc,
	*bufconn.Listener, // this is returned only to bypass the GC so that the connection is not closed.
) {
	components, me, ctx := NewContext(db, t)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

	s := grpc.NewServer(grpc.StreamInterceptor(func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		newCtx := context.WithValue(ss.Context(), "components", components)
		newCtx = context.WithValue(newCtx, "me", me)
		components.ResetEntityStores(newCtx)

		newSS := grpc_middleware.WrapServerStream(ss)
		newSS.WrappedContext = newCtx
		return handler(srv, newSS)
	}))

	withService(s)

	listener := bufconn.Listen(1024 * 1024)
	go func() {
		go func() {
			if err := s.Serve(listener); err != nil {
				t.Fatal(err)
			}
		}()
		select {
		case <-ctx.Done():
			s.Stop()
			t.Fatal("deadline reached")
		}
	}()

	bufDialer := func(context.Context, string) (net.Conn, error) { return listener.Dial() }
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("could not dial bufnet: %v", err)
	}
	return components, me, ctx, conn, cancel, listener
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
		t.Fatal(err)
	}
	return u
}

func AddMockConversation(db common.DBLogger, t *testing.T, users ...*models.User) *models.Conversation {
	g := &models.Conversation{}
	if err := g.Insert(context.Background(), db, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	for _, user := range users {
		ug := models.UserConversation{ConversationID: g.ID, UserID: user.ID}
		if err := ug.Insert(context.Background(), db, boil.Infer()); err != nil {
			t.Fatal(err)
		}
	}
	return g
}

func AddMockMessage(db common.DBLogger, t *testing.T, message *models.Message) *models.Message {
	message.Type = models.MessageTypeText
	message.Text = null.StringFrom("coucou")

	if err := message.Insert(context.Background(), db, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	return message
}
