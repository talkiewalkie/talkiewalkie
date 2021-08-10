package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
	"log"
)

func main() {
	if err := godotenv.Load(".env.dev"); err != nil {
		log.Panicf("could not load env: %v", err)
	}

	components, err := common.InitComponents()
	if err != nil {
		panic(err)
	}

	walks, err := models.Walks().All(context.Background(), components.Db)
	if err != nil {
		panic(err)
	}

	for _, walk := range walks {
		print(walk.Title)
	}
	print("hey")
}
