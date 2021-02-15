package unauthenticated

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/talkiewalkie/talkiewalkie/repository"
)

type unauthenticatedContext struct {
	Db             *sqlx.DB
	WalkRepository repository.WalkRepository
	UserRepository repository.UserRepository
}

func withUnauthContext(w http.ResponseWriter, r *http.Request, block func(c unauthenticatedContext)) {
	db, ok := r.Context().Value("db").(*sqlx.DB)
	if !ok {
		log.Print("no 'db' value in context")
		http.Error(w, "no db value in context", http.StatusInternalServerError)
		return
	}
	walkRepo := repository.PgWalkRepository{Db: db}
	userRepo := repository.PgUserRepository{Db: db}
	block(unauthenticatedContext{Db: db, WalkRepository: walkRepo, UserRepository: userRepo})
}
