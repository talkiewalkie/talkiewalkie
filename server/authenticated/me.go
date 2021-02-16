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

func MeHandler(r *http.Request, c *authenticatedContext) (interface{}, *common.HttpError) {
	return meResponse{User: c.User, Verified: !c.User.EmailToken.Valid}, nil

}
