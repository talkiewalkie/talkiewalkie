package authenticated

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
	Uuid   string
	Handle string
	Email  string
}

func UserListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		withAuthContext(w, r, func(c authenticatedContext) {
			handle := mux.Vars(r)["handle"]
			u, err := c.UserRepository.GetUserByHandle(handle)
			if u == nil || err != nil {
				log.Printf("did not get any user for handle '%s': %v", handle, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write([]byte(u.Handle))
		})
	}
}
