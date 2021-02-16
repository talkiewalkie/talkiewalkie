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

func CreateUserHandler(components *common.Components) UnauthHandler {
	return func(w http.ResponseWriter, r *http.Request, c *unauthenticatedContext) (interface{}, *common.HttpError) {
		var p signInPayload
		if err := common.JsonIn(r, &p); err != nil {
			return nil, common.ServerError(err.Error())
		}

		key := make([]byte, 64)
		if _, err := rand.Read(key); err != nil {
			return nil, common.ServerError("could not generate random key: %v", err)
		}

		emailContent := fmt.Sprintf("bienvue sur takliewalkie, ton code de verif est %x", key)
		if err := components.EmailClient.SendEmail([]byte(emailContent), []string{p.Email}); err != nil {
			// TODO: fix email sending!
			//common.Error(w, fmt.Sprintf("failed to send verification email: %v", err), http.StatusInternalServerError)
			//return
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, common.ServerError("could not hash the password; %v", err)
		}

		dbU, err := c.UserRepository.CreateUser(p.Email, hashed, hex.EncodeToString(key))
		if err != nil {
			return nil, common.ServerError("could not create user in db: %v", err)
		}

		return dbU, nil
	}
}
