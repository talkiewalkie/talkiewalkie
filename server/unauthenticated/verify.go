package unauthenticated

import (
	"fmt"
	"net/http"

	"github.com/talkiewalkie/talkiewalkie/common"
)

func VerifyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		withUnauthContext(w, r, func(c unauthenticatedContext) {
			params := r.URL.Query()

			tokens := params["token"]
			if len(tokens) != 1 || tokens[0] == "" {
				common.Error(w, "received malformed query", http.StatusUnauthorized)
				return
			}

			users := params["user"]
			if len(users) != 1 || users[0] == "" {
				common.Error(w, "bad query", http.StatusUnauthorized)
				return
			}

			u, err := c.UserRepository.GetUserByUuid(users[0])
			if err != nil {
				common.Error(w, fmt.Sprintf("error retrieving user: %v", err), http.StatusInternalServerError)
				return
			}

			if u.EmailToken.String != tokens[0] {
				common.Error(w, "bad email token", http.StatusUnauthorized)
				return
			}

			if _, err = c.Db.Exec(`
				UPDATE "user"
				SET email_token = null
				WHERE uuid = $1;
			`, users[0]); err != nil {
				common.Error(w, fmt.Sprintf("failed to update user row: %v", err), http.StatusInternalServerError)
			}
		})
	}
}
