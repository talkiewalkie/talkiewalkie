package authenticated

import (
	"fmt"
	"net/http"

	"github.com/talkiewalkie/talkiewalkie/common"
)

func UploadHandler(c *common.Components) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		withAuthContext(w, r, func(ctx authenticatedContext) {
			_, h, err := r.FormFile("main")
			if err != nil {
				common.Error(w, fmt.Sprintf("could not find a file to upload: %v", err), http.StatusInternalServerError)
				return
			}

			f, err := h.Open()
			if err != nil {
				common.Error(w, fmt.Sprintf("could not open form file: %v", err), http.StatusInternalServerError)
				return
			}
			defer f.Close()

			uid, err := c.Storage.Upload(r.Context(), f)
			if err != nil {
				common.Error(w, fmt.Sprintf("could not upload file: %v", err), http.StatusInternalServerError)
				return
			}

			p := make([]byte, 200)
			if _, err = f.ReadAt(p, 0); err != nil {
				common.Error(w, fmt.Sprintf("could not read file: %v", err), http.StatusInternalServerError)
				return
			}

			// https://stackoverflow.com/a/52266455
			contentType := http.DetectContentType(p)
			a, err := ctx.AssetRepository.Create(*uid, h.Filename, contentType)
			if err != nil {
				common.Error(w, fmt.Sprintf("could not register asset in db: %v", err), http.StatusInternalServerError)
				return
			}

			common.JsonOut(w, a)
		})
	}
}
