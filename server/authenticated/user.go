package authenticated

import (
	"net/http"
)

func UserListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		withAuthContext(w, r, func(c authenticatedContext) {

		})
	}
}
