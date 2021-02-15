package authenticated

import (
	"context"
	"log"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/gorilla/mux"
)

func NewFirebaseAuth(c *auth.Client) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bearerToken := r.Header.Get("Authorization")
			idToken := strings.TrimSpace(strings.Replace(bearerToken, "Bearer", "", 1))
			ctx := r.Context()
			token, err := c.VerifyIDToken(ctx, idToken)
			if err != nil {
				log.Printf("could not verify token: %v", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			newR := r.WithContext(context.WithValue(ctx, "userUid", token.UID))
			next.ServeHTTP(w, newR)
		})
	}
}
