package authenticated

import (
	"net/http"
)

type User struct {
	Uuid   string
	Handle string
	Email  string
}

func UserListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		withAuthContext(w, r, func(c authenticatedContext) {

		})
	}
}
