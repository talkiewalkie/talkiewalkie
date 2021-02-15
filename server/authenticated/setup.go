package authenticated

import (
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/gorilla/mux"

	"github.com/talkiewalkie/talkiewalkie/common"
)

func Setup(r *mux.Router, c *common.Components) {
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.Use(
		jwtauth.Verifier(c.JwtAuth),
		jwtauth.Authenticator)

	authRouter.HandleFunc("/user/{handle}", UserListHandler()).Methods(http.MethodGet)
	authRouter.HandleFunc("/me", MeHandler()).Methods(http.MethodGet)
	authRouter.HandleFunc("/upload", UploadHandler(c)).Methods(http.MethodPost)
}
