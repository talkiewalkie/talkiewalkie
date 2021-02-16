package authenticated

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/repository"
)

func Setup(r *mux.Router, c *common.Components) {
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.Use(
		jwtauth.Verifier(c.JwtAuth),
		jwtauth.Authenticator)

	authRouter.HandleFunc("/user/{handle}", mountHandler(UserListHandler)).Methods(http.MethodGet)
	authRouter.HandleFunc("/me", mountHandler(MeHandler)).Methods(http.MethodGet)
	authRouter.HandleFunc("/upload", mountHandler(UploadHandler(c))).Methods(http.MethodPost)
	authRouter.HandleFunc("/walk/create", mountHandler(CreateWalkHandler)).Methods(http.MethodPost)
}

type authenticatedContext struct {
	Db              *sqlx.DB
	User            *models.User
	UserRepository  repository.UserRepository
	AssetRepository repository.AssetRepository
}

type AuthHandler func(*http.Request, *authenticatedContext) (interface{}, *common.HttpError)

func mountHandler(handler AuthHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		c, err := buildAuthContext(r)
		if err != nil {
			log.Printf("failed to build context: %v", err)
			http.Error(w, fmt.Sprintf("failed to build context: %v", err), http.StatusInternalServerError)
			return
		}

		response, httpError := handler(r, c)
		if httpError != nil {
			log.Println(httpError.Msg)
			http.Error(w, httpError.Msg, httpError.Code)
			return
		}

		w.WriteHeader(http.StatusOK)
		err = common.JsonOut(w, response)
		if err != nil {
			log.Printf("failed to marshal response: %v", err)
			http.Error(w, fmt.Sprintf("failed to marshal response: %v", err), http.StatusInternalServerError)
		}
	}
}

func buildAuthContext(r *http.Request) (*authenticatedContext, error) {
	db, ok := r.Context().Value("db").(*sqlx.DB)
	if !ok {
		return nil, fmt.Errorf("no 'db' value in context")
	}

	userRepo := repository.PgUserRepository{Db: db}
	assetRepo := repository.PgAssetRepository{Db: db}

	_, claims, _ := jwtauth.FromContext(r.Context())
	userUuid, ok := claims["userUuid"].(string)
	if !ok {
		return nil, fmt.Errorf("no 'userUid' value in context")
	}

	user, err := userRepo.GetUserByUuid(userUuid)
	if err != nil {
		return nil, err
	}

	return &authenticatedContext{Db: db, User: user, UserRepository: userRepo, AssetRepository: assetRepo}, nil
}
