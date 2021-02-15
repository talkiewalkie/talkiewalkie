package authenticated

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/talkiewalkie/talkiewalkie/common"
)

func Setup(r *mux.Router, c *common.Components) {
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.Use(NewFirebaseAuth(c.FirebaseAuth))

	authRouter.HandleFunc("/user/{handle}", UserListHandler()).Methods(http.MethodGet)
}
