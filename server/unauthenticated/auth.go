package unauthenticated

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/talkiewalkie/talkiewalkie/common"
)

type signInPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUserHandler(components *common.Components) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		withUnauthContext(w, r, func(c unauthenticatedContext) {
			var p signInPayload
			if err := common.JsonIn(w, r, &p); err != nil {
				return
			}

			key := make([]byte, 64)
			if _, err := rand.Read(key); err != nil {
				http.Error(w, fmt.Sprintf("could not generate random key: %v", err), http.StatusInternalServerError)
				return
			}

			emailContent := fmt.Sprintf("bienvue sur takliewalkie, ton code de verif est %x", key)
			if err := components.EmailClient.SendEmail(p.Email, []byte(emailContent)); err != nil {
				http.Error(w, fmt.Sprintf("failed to send verification email: %v", err), http.StatusInternalServerError)
				return
			}

			hashed, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, fmt.Sprintf("could not hash the password; %v", err), http.StatusInternalServerError)
				return
			}

			dbU, err := c.UserRepository.CreateUser(p.Email, hex.EncodeToString(hashed), hex.EncodeToString(key))
			if err != nil {
				http.Error(w, fmt.Sprintf("could not create user in db: %v", err), http.StatusInternalServerError)
				return
			}

			common.JsonOut(w, dbU)
		})
	}
}

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
				http.Error(w, "did not find the user in db: %v", http.StatusUnauthorized)
				return
			}

			hashed, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, fmt.Sprintf("could not hash the password; %v", err), http.StatusInternalServerError)
				return
			}

			if u.Password != hex.EncodeToString(hashed) {
				http.Error(w, "passwords don't match", http.StatusUnauthorized)
				return
			}

			_, signed, err := components.JwtAuth.Encode(map[string]interface{}{"userUuid": u.Uuid})
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to build jwt: %v", err), http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:    "jwt",
				Value:   signed,
				Path:    "/",
				Expires: time.Now().Add(time.Hour),
			})

			common.JsonOut(w, u)
		})
	}
}
