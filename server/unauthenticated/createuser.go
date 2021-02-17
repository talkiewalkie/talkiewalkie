package unauthenticated

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"golang.org/x/crypto/bcrypt"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
)

type createUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request, c *unauthenticatedContext) (interface{}, *common.HttpError) {
	var p createUserInput
	if err := common.JsonIn(r, &p); err != nil {
		return nil, common.ServerError(err.Error())
	}

	key := make([]byte, 64)
	if _, err := rand.Read(key); err != nil {
		return nil, common.ServerError("could not generate random key: %v", err)
	}

	emailContent := fmt.Sprintf("bienvue sur takliewalkie, ton code de verif est %x", key)
	if err := c.EmailClient.SendEmail([]byte(emailContent), []string{p.Email}); err != nil {
		// TODO: fix email sending!
		//common.Error(w, fmt.Sprintf("failed to send verification email: %v", err), http.StatusInternalServerError)
		//return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, common.ServerError("could not hash the password; %v", err)
	}

	u := &models.User{Handle: p.Email, Email: p.Email, Password: hashed, EmailToken: null.NewString(hex.EncodeToString(key), true)}
	if err = u.Insert(r.Context(), c.Db, boil.Infer()); err != nil {
		return nil, common.ServerError("could not create user in db: %v", err)
	}

	return u, nil
}
