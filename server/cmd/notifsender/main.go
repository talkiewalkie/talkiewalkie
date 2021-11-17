package main

import (
	"context"
	"flag"
	"log"

	"github.com/joho/godotenv"

	"github.com/talkiewalkie/talkiewalkie/clients"
	"github.com/talkiewalkie/talkiewalkie/common"
)

var (
	fbUid = flag.String("fbUid", "", "firebase uid")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	if *fbUid == "" {
		log.Panic("use with '-fbUid <some uid>'")
	}

	if err := godotenv.Load(".env.dev"); err != nil {
		log.Panicf("could not load env: %+v", err)
	}

	components := common.InitComponents()

	res, err := components.MessagingClient.Send(ctx, clients.MessageInput{
		Topic: *fbUid,
		Data:  map[string]string{"uuid": "this is not a uuid haha you fell in my trap"},
		Body:  "hey",
		Title: "theo",
	})

	if err != nil {
		log.Printf("err: %+v", err)
	}
	log.Printf("sucess: '%s'", res)
}
