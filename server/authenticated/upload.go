package authenticated

import (
	"net/http"

	"github.com/talkiewalkie/talkiewalkie/common"
)

func UploadHandler(c *common.Components) AuthHandler {
	return func(r *http.Request, ctx *authenticatedContext) (interface{}, *common.HttpError) {
		_, h, err := r.FormFile("main")
		if err != nil {
			return nil, common.ServerError("could not find a file to upload: %v", err)
		}

		f, err := h.Open()
		if err != nil {
			return nil, common.ServerError("could not open form file: %v", err)
		}
		defer f.Close()

		uid, err := c.Storage.Upload(r.Context(), f)
		if err != nil {
			return nil, common.ServerError("could not upload file: %v", err)
		}

		p := make([]byte, 200)
		if _, err = f.ReadAt(p, 0); err != nil {
			return nil, common.ServerError("could not read file: %v", err)
		}

		// todo : extend to better detection - https://stackoverflow.com/a/52266455
		contentType := http.DetectContentType(p)
		a, err := ctx.AssetRepository.Create(*uid, h.Filename, contentType)
		if err != nil {
			return nil, common.ServerError("could not register asset in db: %v", err)
		}

		return a, nil
	}
}
