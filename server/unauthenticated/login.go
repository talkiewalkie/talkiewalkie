package unauthenticated

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/talkiewalkie/talkiewalkie/common"
)

type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(components *common.Components) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		withUnauthContext(w, r, func(c unauthenticatedContext) {
			var p loginPayload
			if err := common.JsonIn(w, r, &p); err != nil {
				return
			}

			u, err := c.UserRepository.GetUserByEmail(p.Email)
			if err != nil {
				common.Error(w, "did not find the user in db: %v", http.StatusUnauthorized)
				return
			}

			if bcrypt.CompareHashAndPassword(u.Password, []byte(p.Password)) != nil {
				common.Error(w, "passwords don't match", http.StatusUnauthorized)
				return
			}

			_, signed, err := components.JwtAuth.Encode(map[string]interface{}{"userUuid": u.Uuid})
			if err != nil {
				common.Error(w, fmt.Sprintf("failed to build jwt: %v", err), http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:    "jwt",
				Value:   signed,
				Path:    "/",
				Expires: time.Now().Add(time.Hour),
			})

			w.WriteHeader(http.StatusOK)
		})
	}
}
