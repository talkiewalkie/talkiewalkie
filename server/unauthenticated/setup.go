package unauthenticated

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/repository"
)

func Setup(r *mux.Router, c *common.Components) {
	unauthRouter := r.PathPrefix("/unauth").Subrouter()
	unauthRouter.HandleFunc("/walks", WalksHandler()).Methods(http.MethodGet)
	unauthRouter.HandleFunc("/create", CreateUserHandler(c)).Methods(http.MethodPost)
	unauthRouter.HandleFunc("/login", LoginHandler(c)).Methods(http.MethodPost)
	unauthRouter.HandleFunc("/verify", VerifyHandler()).Methods(http.MethodGet)
}

type unauthenticatedContext struct {
	Db             *sqlx.DB
	WalkRepository repository.WalkRepository
	UserRepository repository.UserRepository
}

func withUnauthContext(w http.ResponseWriter, r *http.Request, block func(c unauthenticatedContext)) {
	db, ok := r.Context().Value("db").(*sqlx.DB)
	if !ok {
		common.Error(w, "no db value in context", http.StatusInternalServerError)
		return
	}
	walkRepo := repository.PgWalkRepository{Db: db}
	userRepo := repository.PgUserRepository{Db: db}
	block(unauthenticatedContext{Db: db, WalkRepository: walkRepo, UserRepository: userRepo})
}
