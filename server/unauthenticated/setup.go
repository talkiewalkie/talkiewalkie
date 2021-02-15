package unauthenticated

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/talkiewalkie/talkiewalkie/common"
)

func Setup(r *mux.Router, c *common.Components) {
	unauthRouter := r.PathPrefix("/unauth").Subrouter()
	unauthRouter.HandleFunc("/walks", WalksHandler()).Methods(http.MethodGet)
	unauthRouter.HandleFunc("/create", CreateUserHandler(c)).Methods(http.MethodPost)
	unauthRouter.HandleFunc("/login", LoginHandler(c)).Methods(http.MethodPost)
}
