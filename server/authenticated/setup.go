package authenticated

import (
	"fmt"
	"log"
	"net/http"
	"time"

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

	authRouter.HandleFunc("/user/{uuid}", mountHandler(c, UserByUuidHandler)).Methods(http.MethodGet)
	authRouter.HandleFunc("/walk/{uuid}", mountHandler(c, WalkHandler)).Methods(http.MethodGet)
	authRouter.HandleFunc("/me", mountHandler(c, MeHandler)).Methods(http.MethodGet)

	authRouter.HandleFunc("/upload", mountHandler(c, UploadHandler)).Methods(http.MethodPost)
	authRouter.HandleFunc("/walk/create", mountHandler(c, CreateWalkHandler)).Methods(http.MethodPost)
	authRouter.HandleFunc("/walk/like/{uuid}", mountHandler(c, LikeWalkByUuid)).Methods(http.MethodPost)
}

type authenticatedContext struct {
	*common.Components
	Db              *sqlx.DB
	User            *models.User
	UserRepository  repository.UserRepository
	AssetRepository repository.AssetRepository
}

type AuthHandler func(*http.Request, *authenticatedContext) (interface{}, *common.HttpError)

func mountHandler(components *common.Components, handler AuthHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		c, err := buildAuthContext(components, r)
		if err != nil {
			log.Printf("failed to build context: %+v", err)
			http.Error(w, fmt.Sprintf("failed to build context: %+v", err), http.StatusInternalServerError)
			return
		}

		response, httpError := handler(r, c)
		if httpError != nil {
			log.Println(httpError.Msg)
			http.Error(w, httpError.Msg, httpError.Code)
			return
		}

		_, signed, err := c.JwtAuth.Encode(map[string]interface{}{"userUuid": c.User.UUID})
		if err != nil {
			log.Printf("failed to build jwt: %+v", err)
			http.Error(w, fmt.Sprintf("failed to build jwt: %+v", err), http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    signed,
			Path:     "/",
			Expires:  time.Now().Add(time.Hour),
			Secure:   true,
			HttpOnly: true,
		})

		w.WriteHeader(http.StatusOK)
		err = common.JsonOut(w, response)
		if err != nil {
			log.Printf("failed to marshal response: %v", err)
			http.Error(w, fmt.Sprintf("failed to marshal response: %v", err), http.StatusInternalServerError)
		}
	}
}

func buildAuthContext(c *common.Components, r *http.Request) (*authenticatedContext, error) {
	db, ok := r.Context().Value("db").(*sqlx.DB)
	if !ok {
		return nil, fmt.Errorf("no 'db' value in context")
	}

	userRepo := repository.PgUserRepository{Components: c, Db: db, Ctx: r.Context()}
	assetRepo := repository.PgAssetRepository{Components: c, Db: db, Ctx: r.Context()}

	_, claims, _ := jwtauth.FromContext(r.Context())
	userUuid, ok := claims["userUuid"].(string)
	if !ok {
		return nil, fmt.Errorf("no 'userUid' value in context")
	}

	user, err := userRepo.GetUserByUuid(userUuid)
	if err != nil {
		return nil, err
	}

	return &authenticatedContext{
		Components:      c,
		Db:              db,
		User:            user,
		UserRepository:  userRepo,
		AssetRepository: assetRepo,
	}, nil
}
