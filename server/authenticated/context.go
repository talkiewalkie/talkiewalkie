package authenticated

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/jmoiron/sqlx"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/repository"
)

type authenticatedContext struct {
	Db              *sqlx.DB
	User            *repository.User
	UserRepository  repository.UserRepository
	AssetRepository repository.AssetRepository
}

func withAuthContext(w http.ResponseWriter, r *http.Request, block func(c authenticatedContext)) {
	db, ok := r.Context().Value("db").(*sqlx.DB)
	if !ok {
		log.Print("no 'db' value in context")
		common.Error(w, "", http.StatusInternalServerError)
		return
	}
	userRepo := repository.PgUserRepository{Db: db}
	assetRepo := repository.PgAssetRepository{Db: db}

	_, claims, _ := jwtauth.FromContext(r.Context())
	userUuid, ok := claims["userUuid"].(string)
	if !ok {
		log.Printf("no 'userUid' value in context")
		common.Error(w, "", http.StatusInternalServerError)
		return
	}
	user, err := userRepo.GetUserByUuid(userUuid)
	if err != nil {
		common.Error(w, fmt.Sprintf("failed to retrieve user '%v': %v", userUuid, err), http.StatusUnauthorized)
		return
	}

	block(authenticatedContext{Db: db, User: user, UserRepository: userRepo, AssetRepository: assetRepo})
}
