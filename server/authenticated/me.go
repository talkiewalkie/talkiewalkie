package authenticated

import (
	"net/http"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/repository"
)

type meResponse struct {
	*repository.User
	Verified bool `json:"verified"`
}

func MeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		withAuthContext(w, r, func(c authenticatedContext) {
			common.JsonOut(w, meResponse{User: c.User, Verified: !c.User.EmailToken.Valid})
		})
	}
}
