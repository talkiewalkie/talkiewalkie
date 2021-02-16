package unauthenticated

import (
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/talkiewalkie/talkiewalkie/common"
)

type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(components *common.Components) UnauthHandler {
	return func(w http.ResponseWriter, r *http.Request, c *unauthenticatedContext) (interface{}, *common.HttpError) {
		var p loginPayload
		if err := common.JsonIn(r, &p); err != nil {
			return nil, common.ServerError(err.Error())
		}

		unauthErr := &common.HttpError{Code: http.StatusUnauthorized}
		u, err := c.UserRepository.GetUserByEmail(p.Email)
		if err != nil {
			return nil, unauthErr
		}

		if bcrypt.CompareHashAndPassword(u.Password, []byte(p.Password)) != nil {
			return nil, unauthErr
		}

		_, signed, err := components.JwtAuth.Encode(map[string]interface{}{"userUuid": u.UUID})
		if err != nil {
			return nil, common.ServerError("failed to build jwt: %v", err)
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "jwt",
			Value:   signed,
			Path:    "/",
			Expires: time.Now().Add(time.Hour),
		})

		return nil, nil
	}
}
