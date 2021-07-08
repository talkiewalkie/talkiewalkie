package common

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/models"
)

type Context struct {
	Components *Components
	User       *models.User
}

func WithContextMiddleWare(comps *Components) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var u *models.User

			_, claims, _ := jwtauth.FromContext(r.Context())
			userUuid, ok := claims["userUuid"].(string)
			if ok {
				uid, err := uuid.FromString(userUuid)
				if err != nil {
					http.Error(w, "bad uuid", http.StatusInternalServerError)
					return
				}

				u, err = models.Users(models.UserWhere.UUID.EQ(uid)).One(r.Context(), comps.Db)
				if err != nil {
					http.Error(w, "no user for uuid", http.StatusInternalServerError)
					return
				}
			}

			myCtx := Context{
				Components: comps,
				User:       u,
			}
			ctx := context.WithValue(r.Context(), "context", myCtx)
			newR := r.WithContext(ctx)
			next.ServeHTTP(w, newR)
		})
	}
}

func WithContext(r *http.Request) Context {
	ctx := r.Context()

	services, ok := ctx.Value("context").(Context)
	if !ok {
		panic("failed to get services from context")
	}

	return services
}

func WithAuthedContext(r *http.Request) Context {
	ctx := r.Context()

	services, ok := ctx.Value("context").(Context)
	if !ok {
		panic("failed to get services from context")
	}

	if services.User == nil {
		panic(errors.New("auth"))
	}

	return services
}
