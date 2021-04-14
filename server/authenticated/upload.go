package authenticated

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	uuid "github.com/satori/go.uuid"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/pb"
)

type uploadOutput struct {
	Uuid uuid.UUID `json:"uuid"`
}

func UploadHandler(r *http.Request, ctx *authenticatedContext) (interface{}, *common.HttpError) {
	f, h, err := r.FormFile("main")
	if err != nil {
		return nil, common.ServerError("could not find a file to upload: %v", err)
	}
	defer f.Close()

	// todo : extend to better detection - https://stackoverflow.com/a/52266455
	p := make([]byte, 200)
	if _, err = f.ReadAt(p, 0); err != nil {
		return nil, common.ServerError("could not read file: %v", err)
	}
	contentType := http.DetectContentType(p)

	var uploadedF io.Reader
	if strings.HasPrefix(contentType, "image/") {
		fsF, err := ioutil.TempFile("", h.Filename)
		if err != nil {
			return nil, common.ServerError("could not create temporary file: %+v", err)
		}
		_, _ = f.Seek(0, io.SeekStart)
		if _, err = io.Copy(fsF, f); err != nil {
			return nil, common.ServerError("could not create temporary file: %+v", err)
		}

		compressed, err := ctx.CompressImg(fsF.Name(), 600)
		if err != nil {
			return nil, common.ServerError("failed to compress image: %+v", err)
		}

		uf, err := os.Open(compressed)
		if err != nil {
			return nil, common.ServerError("could not open compressed image: %+v", err)
		}

		uploadedF = uf
	} else if ctx.Audio != nil && (strings.HasPrefix(contentType, "video/") || strings.HasPrefix(contentType, "audio/")) {
		_, _ = f.Seek(0, io.SeekStart)
		content, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, common.ServerError("could not read file: %+v", err)
		}

		output, err := ctx.Audio.FormatAndCompress(r.Context(), &pb.FormatAndCompressInput{
			Content:  content,
			FileName: h.Filename,
			MimeType: contentType,
		})
		if err != nil {
			return nil, common.ServerError("failed to compress audio: %+v", err)
		}

		if len(output.Content) != len(content) {
			return nil, common.ServerError("audio service error: sent %v bytes received %v", len(content), len(output.Content))
		}
		uploadedF = bytes.NewReader(output.Content)
	} else {
		uploadedF = f
	}

	uid, err := ctx.Storage.Upload(r.Context(), uploadedF)
	if err != nil {
		return nil, common.ServerError("could not upload file: %v", err)
	}

	a, err := ctx.AssetRepository.Create(*uid, h.Filename, contentType)
	if err != nil {
		return nil, common.ServerError("could not register asset in db: %v", err)
	}

	return uploadOutput{a.UUID}, nil
}
