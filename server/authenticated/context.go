package authenticated

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/talkiewalkie/talkiewalkie/repository"
)

type authenticatedContext struct {
	Db             *sqlx.DB
	User           *repository.User
	UserRepository repository.UserRepository
}

func withAuthContext(w http.ResponseWriter, r *http.Request, block func(c authenticatedContext)) {
	db, ok := r.Context().Value("db").(*sqlx.DB)
	if !ok {
		log.Print("no 'db' value in context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userRepo := repository.PgUserRepository{Db: db}

	userUid, ok := r.Context().Value("userUid").(string)
	if !ok {
		log.Printf("no 'userUid' value in context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user, err := userRepo.GetUserByHandle(userUid)
	if err != nil {
		log.Printf("failed to retrieve user '%v'", userUid)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	block(authenticatedContext{Db: db, User: user, UserRepository: userRepo})
}
