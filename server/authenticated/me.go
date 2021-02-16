package authenticated

import (
	"net/http"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
)

type meResponse struct {
	*models.User
	Verified bool `json:"verified"`
}

func MeHandler(r *http.Request, c *authenticatedContext) (interface{}, *common.HttpError) {
	return meResponse{User: c.User, Verified: !c.User.EmailToken.Valid}, nil

}
