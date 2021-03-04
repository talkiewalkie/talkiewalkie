package unauthenticated

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/repository"
)

func Setup(r *mux.Router, c *common.Components) {
	unauthRouter := r.PathPrefix("/unauth").Subrouter()
	unauthRouter.HandleFunc("/verify", mountHandler(c, VerifyHandler)).Methods(http.MethodGet)
	unauthRouter.HandleFunc("/walk/{uuid}", mountHandler(c, WalkHandler)).Methods(http.MethodGet)
	unauthRouter.HandleFunc("/walks", mountHandler(c, WalksHandler)).Methods(http.MethodGet)

	unauthRouter.HandleFunc("/login", mountHandler(c, LoginHandler)).Methods(http.MethodPost)
	unauthRouter.HandleFunc("/user/create", mountHandler(c, CreateUserHandler)).Methods(http.MethodPost)
}

type unauthenticatedContext struct {
	*common.Components
	Db             *sqlx.DB
	WalkRepository repository.WalkRepository
	UserRepository repository.UserRepository
}

type UnauthHandler func(http.ResponseWriter, *http.Request, *unauthenticatedContext) (interface{}, *common.HttpError)

func mountHandler(components *common.Components, handler UnauthHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := buildUnauthContext(components, r)
		if err != nil {
			log.Printf("failed to build context: %v", err)
			http.Error(w, fmt.Sprintf("failed to build context: %v", err), http.StatusInternalServerError)
			return
		}

		response, httpError := handler(w, r, c)
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

func buildUnauthContext(components *common.Components, r *http.Request) (*unauthenticatedContext, error) {
	db, ok := r.Context().Value("db").(*sqlx.DB)
	if !ok {
		return nil, fmt.Errorf("no 'db' value in context")
	}

	walkRepo := repository.PgWalkRepository{Components: components, Db: db, Ctx: r.Context()}
	userRepo := repository.PgUserRepository{Components: components, Db: db, Ctx: r.Context()}

	return &unauthenticatedContext{components, db, walkRepo, userRepo}, nil
}
