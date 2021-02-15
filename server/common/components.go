package common

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

type Components struct {
	FirebaseApp  *firebase.App
	FirebaseAuth *auth.Client
	EmailClient  *EmailClient
}

func InitComponents() Components {
	fbClient, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Panicf("could not load firebase client: %v", err)
	}

	c, err := fbClient.Auth(context.Background())
	if err != nil {
		log.Panicf("could not load firebase auth: %v", err)
	}

	emailClient := initEmailClient()
	return Components{FirebaseApp: fbClient, FirebaseAuth: c, EmailClient: &emailClient}
}
