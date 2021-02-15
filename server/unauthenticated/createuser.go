package unauthenticated

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

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
				common.Error(w, fmt.Sprintf("could not generate random key: %v", err), http.StatusInternalServerError)
				return
			}

			emailContent := fmt.Sprintf("bienvue sur takliewalkie, ton code de verif est %x", key)
			if err := components.EmailClient.SendEmail(p.Email, []byte(emailContent)); err != nil {
				// TODO: fix email sending!
				//common.Error(w, fmt.Sprintf("failed to send verification email: %v", err), http.StatusInternalServerError)
				//return
			}

			hashed, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
			if err != nil {
				common.Error(w, fmt.Sprintf("could not hash the password; %v", err), http.StatusInternalServerError)
				return
			}

			dbU, err := c.UserRepository.CreateUser(p.Email, hashed, hex.EncodeToString(key))
			if err != nil {
				common.Error(w, fmt.Sprintf("could not create user in db: %v", err), http.StatusInternalServerError)
				return
			}

			common.JsonOut(w, dbU)
		})
	}
}
