package routes

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/testutils"
)

func TestAssets(t *testing.T) {
	db := testutils.SetupDb()

	t.Run("can create assets (image)", createImageAssetUploadTest(db))
	t.Run("can create assets (audio)", createAudioAssetUploadTest(db))
	testutils.TearDownDb(db)
}

var (
	pngMagicBytes  = "\x89PNG\r\n\x1a\n"
	webmMagicBytes = "\x1A\x45\xDF\xA3"
)

func createImageAssetUploadTest(db *sqlx.DB) func(t *testing.T) {
	return func(t *testing.T) {
		user := testutils.AddMockUser(db, t)

		body := new(bytes.Buffer)
		formWriter := multipart.NewWriter(body)
		mimeHeader := textproto.MIMEHeader{}
		mimeHeader.Set("content-type", "image/png")
		formWriter.CreatePart(mimeHeader)
		fileWriter, _ := formWriter.CreateFormFile("main", "test.png")
		fileWriter.Write([]byte(pngMagicBytes))
		formWriter.Close()

		req := testutils.MakeRequest(user, db, http.MethodPost, "/asset", body)
		req.Header.Set("Content-Type", fmt.Sprintf(formWriter.FormDataContentType()))
		w := &httptest.ResponseRecorder{}

		UploadHandler(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		assets, _ := models.Assets().All(req.Context(), db)
		assert.NotEmptyf(t, assets, "no asset in db")
	}
}

func createAudioAssetUploadTest(db *sqlx.DB) func(t *testing.T) {
	return func(t *testing.T) {
		user := testutils.AddMockUser(db, t)

		body := new(bytes.Buffer)
		formWriter := multipart.NewWriter(body)
		mimeHeader := textproto.MIMEHeader{}
		mimeHeader.Set("content-type", "image/png")
		formWriter.CreatePart(mimeHeader)
		fileWriter, _ := formWriter.CreateFormFile("main", "test.png")
		fileWriter.Write([]byte(webmMagicBytes))
		formWriter.Close()

		req := testutils.MakeRequest(user, db, http.MethodPost, "/asset", body)
		req.Header.Set("Content-Type", fmt.Sprintf(formWriter.FormDataContentType()))
		w := &httptest.ResponseRecorder{}

		UploadHandler(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		assets, _ := models.Assets().All(req.Context(), db)
		assert.NotEmptyf(t, assets, "no asset in db")
	}
}
