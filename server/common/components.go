package common

import (
	"context"
	"log"
	"os"

	"github.com/go-chi/jwtauth"
)

type Components struct {
	EmailClient EmailClient
	JwtAuth     *jwtauth.JWTAuth
	Storage     StorageClient
}

func InitComponents() Components {
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)
	emailClient := initEmailClient()

	storageClient, err := initStorageClient(context.Background())
	if err != nil {
		log.Panicf("could not init the storage ")
	}

	return Components{
		EmailClient: emailClient,
		JwtAuth:     tokenAuth,
		Storage:     storageClient,
	}
}
