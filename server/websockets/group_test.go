package websockets

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/routes"
	"github.com/talkiewalkie/talkiewalkie/testutils"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestMessageRoutes(t *testing.T) {
	db := testutils.SetupDb()

	testutils.TearDownDb(db)
	t.Run("can reach websocket", websocketCanSendAndReceive(db))
}

func fakeMiddleware(u *models.User, db *sqlx.DB, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerFunc(w, testutils.AddFakeComponentsToRequest(r, u, db))
	}
}

func setupWs(t *testing.T, u *models.User, db *sqlx.DB, handlerFunc http.HandlerFunc) (*httptest.Server, *websocket.Conn) {
	s := httptest.NewServer(fakeMiddleware(u, db, handlerFunc))
	wsURL, _ := url.Parse(s.URL)

	switch wsURL.Scheme {
	case "http":
		wsURL.Scheme = "ws"
	case "https":
		wsURL.Scheme = "wss"
	}

	ws, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		t.Fatal(err)
	}
	return s, ws
}

func websocketCanSendAndReceive(db *sqlx.DB) func(t *testing.T) {
	return func(t *testing.T) {
		userA := testutils.AddMockUser(db, t)
		sA, wsA := setupWs(t, userA, db, GroupWebsocketHandler)
		defer sA.Close()
		defer wsA.Close()

		userB := testutils.AddMockUser(db, t)
		sB, wsB := setupWs(t, userB, db, GroupWebsocketHandler)
		defer sB.Close()
		defer wsB.Close()

		userC := testutils.AddMockUser(db, t)
		sC, wsC := setupWs(t, userC, db, GroupWebsocketHandler)
		defer sC.Close()
		defer wsC.Close()

		msg := routes.MessageInput{Handles: []string{userB.Handle, userA.Handle, userC.Handle}, Text: "hey its me userb"}
		body, _ := json.Marshal(msg)
		req := testutils.MakeRequest(userB, db, http.MethodPost, "/message", bytes.NewBuffer(body))
		w := &httptest.ResponseRecorder{}
		routes.Message(w, req)

		buf := new(strings.Builder)
		io.Copy(buf, w.Result().Body)
		require.Equal(t, http.StatusOK, w.Result().StatusCode, "msg: '%s'", buf.String())

		msgChecker := func(m []byte) bool {
			var inp common.NewMessageEvent
			err := json.Unmarshal(m, &inp)
			if err != nil {
				return false
			}
			return inp.Message == "hey its me userb"
		}
		testutils.RequireReceive(t, wsB, time.Second, 0, msgChecker)
		testutils.RequireReceive(t, wsA, time.Second, 1, msgChecker)
		testutils.RequireReceive(t, wsC, time.Second, 1, msgChecker)
	}
}
