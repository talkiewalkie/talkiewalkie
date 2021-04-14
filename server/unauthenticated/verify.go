package unauthenticated

import (
	"net/http"

	uuid "github.com/satori/go.uuid"

	"github.com/talkiewalkie/talkiewalkie/common"
)

func VerifyHandler(w http.ResponseWriter, r *http.Request, c *unauthenticatedContext) (interface{}, *common.HttpError) {
	params := r.URL.Query()
	badQueryErr := &common.HttpError{Code: http.StatusUnauthorized, Msg: "received malformed query"}

	tokens := params["token"]
	if len(tokens) != 1 || tokens[0] == "" {
		return nil, badQueryErr
	}

	users := params["user"]
	if len(users) != 1 || users[0] == "" {
		return nil, badQueryErr
	}

	uid, err := uuid.FromString(users[0])
	if err != nil {
		return nil, common.ServerError("failed to parse uuid: '%s': %+v", users[0], err)
	}
	u, err := c.UserRepository.GetUserByUuid(uid)
	if err != nil {
		return nil, common.ServerError(err.Error())
	}

	if u.EmailToken.String != tokens[0] {
		return nil, &common.HttpError{Code: http.StatusUnauthorized, Msg: "bad email token"}
	}

	if _, err = c.Db.Exec(`
				UPDATE "user"
				SET email_token = null
				WHERE uuid = $1;
			`, users[0]); err != nil {
		return nil, common.ServerError("failed to update user row: %v", err)
	}

	// todo: redirect & set cookie?
	return nil, nil
}
