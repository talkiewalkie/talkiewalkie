package authenticated

import (
	"net/http"

	uuid "github.com/satori/go.uuid"

	"github.com/talkiewalkie/talkiewalkie/common"
)

type meOutput struct {
	UUID     uuid.UUID `json:"uuid"`
	Handle   string    `json:"handle"`
	Email    string    `json:"email"`
	Verified bool      `json:"verified"`
}

func MeHandler(r *http.Request, c *authenticatedContext) (interface{}, *common.HttpError) {
	return meOutput{c.User.UUID, c.User.Handle, c.User.Email, !c.User.EmailToken.Valid}, nil
}
