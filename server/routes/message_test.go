package routes

import (
	"bytes"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/testutils"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMessageRoutes(t *testing.T) {
	db := testutils.SetupDb()

	testutils.TearDownDb(db)
	t.Run("can send one message to one person without preexisting conversation", oneOneWithoutPreexistingConversation(db))
	testutils.TearDownDb(db)
	t.Run("can send one message to one person with preexisting conversation", oneOneWithPreexistingConversation(db))
}

func oneOneWithoutPreexistingConversation(db *sqlx.DB) func(t *testing.T) {
	return func(t *testing.T) {
		userA := testutils.AddMockUser(db, t)
		userB := testutils.AddMockUser(db, t)

		msg := MessageInput{Handles: []string{userB.Handle}, Text: "hey"}
		body, _ := json.Marshal(msg)
		req := testutils.MakeRequest(userA, db, http.MethodPost, "/message", bytes.NewBuffer(body))
		w := &httptest.ResponseRecorder{}
		Message(w, req)

		buf := new(strings.Builder)
		io.Copy(buf, w.Result().Body)
		require.Equal(t, http.StatusOK, w.Result().StatusCode, "msg: '%s'", buf.String())

		cnt, err := models.Groups().Count(req.Context(), db)
		require.Equal(t, int64(1), cnt, "err: %+v", err)
	}
}

func oneOneWithPreexistingConversation(db *sqlx.DB) func(t *testing.T) {
	return func(t *testing.T) {
		userA := testutils.AddMockUser(db, t)
		userB := testutils.AddMockUser(db, t)
		testutils.AddMockGroup(db, t, userA, userB)

		msg := MessageInput{Handles: []string{userB.Handle}, Text: "hey"}
		body, _ := json.Marshal(msg)
		req := testutils.MakeRequest(userA, db, http.MethodPost, "/message", bytes.NewBuffer(body))
		w := &httptest.ResponseRecorder{}
		Message(w, req)

		buf := new(strings.Builder)
		io.Copy(buf, w.Result().Body)
		require.Equal(t, http.StatusOK, w.Result().StatusCode, "msg: '%s'", buf.String())

		cnt, err := models.Groups().Count(req.Context(), db)
		require.Equal(t, int64(1), cnt, "err: %+v", err)
	}
}
