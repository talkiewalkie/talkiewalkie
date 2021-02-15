package unauthenticated

import (
	"fmt"
	"net/http"

	"github.com/talkiewalkie/talkiewalkie/common"
)

func WalksHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		withUnauthContext(w, r, func(c unauthenticatedContext) {
			walks, err := c.WalkRepository.GetAll()
			if err != nil {
				common.Error(w, fmt.Sprintf("could not fetch walks: %v", err), http.StatusInternalServerError)
				return
			}
			common.JsonOut(w, walks)
		})
	}
}
