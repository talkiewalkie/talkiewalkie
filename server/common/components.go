package common

import (
	"os"

	"github.com/go-chi/jwtauth"
)

type Components struct {
	EmailClient *EmailClient
	JwtAuth     *jwtauth.JWTAuth
}

func InitComponents() Components {
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)
	emailClient := initEmailClient()
	return Components{EmailClient: &emailClient, JwtAuth: tokenAuth}
}
