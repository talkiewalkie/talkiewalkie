package routes

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/friendsofgo/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type UploadOutput struct {
	Uuid uuid.UUID `json:"uuid"`
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	ctx := common.WithAuthedContext(r)

	f, h, err := r.FormFile("main")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// todo : extend to better detection - https://stackoverflow.com/a/52266455
	p := make([]byte, 200)
	if _, err = f.ReadAt(p, 0); err != nil {
		panic(errors.New(fmt.Sprintf("%s: %+v", "could not read file: %v", err)))
	}
	contentType := http.DetectContentType(p)

	var uploadedF io.Reader
	if strings.HasPrefix(contentType, "image/") {
		fsF, err := ioutil.TempFile("", h.Filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s: %+v", "could not create temporary file: %+v", err), http.StatusInternalServerError)
		}
		_, _ = f.Seek(0, io.SeekStart)
		if _, err = io.Copy(fsF, f); err != nil {
			http.Error(w, fmt.Sprintf("%s: %+v", "could not create temporary file: %+v", err), http.StatusInternalServerError)
		}

		compressed, err := ctx.Components.CompressImg(fsF.Name(), 600)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s: %+v", "failed to compress image: %+v", err), http.StatusInternalServerError)
		}

		uf, err := os.Open(compressed)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s: %+v", "could not open compressed image: %+v", err), http.StatusInternalServerError)
		}

		uploadedF = uf
	} else if ctx.Components.Audio != nil && (strings.HasPrefix(contentType, "video/") || strings.HasPrefix(contentType, "audio/")) {
		_, _ = f.Seek(0, io.SeekStart)
		content, err := ioutil.ReadAll(f)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s: %+v", "could not read file: %+v", err), http.StatusInternalServerError)
		}

		output, err := (*ctx.Components.Audio).FormatAndCompress(r.Context(), &pb.FormatAndCompressInput{
			Content:  content,
			FileName: h.Filename,
			MimeType: contentType,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("%s: %+v", "failed to compress audio: %+v", err), http.StatusInternalServerError)
		}

		if len(output.Content) != len(content) {
			http.Error(w, fmt.Sprintf("%s: %+v", "audio service error: sent %v bytes received %v", err), http.StatusInternalServerError)
		}
		uploadedF = bytes.NewReader(output.Content)
	} else {
		uploadedF = f
	}

	uid, err := ctx.Components.Storage.Upload(r.Context(), uploadedF)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s: %+v", "could not upload file: %v", err), http.StatusInternalServerError)
	}

	a := &models.Asset{UUID: *uid, FileName: h.Filename, MimeType: contentType}
	if err = a.Insert(r.Context(), ctx.Components.Db, boil.Infer()); err != nil {
		http.Error(w, fmt.Sprintf("%s: %+v", "could not register asset in db: %v", err), http.StatusInternalServerError)
	}

	common.JsonOut(w, UploadOutput{a.UUID})
}
