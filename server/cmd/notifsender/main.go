package main

import (
	"context"
	"firebase.google.com/go/v4/messaging"
	"flag"
	"github.com/joho/godotenv"
	"github.com/talkiewalkie/talkiewalkie/common"
	"log"
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

	components, err := common.InitComponents()
	if err != nil {
		log.Panicf("could not initiate components: %+v", err)
	}

	res, err := components.FbMssg.Send(ctx, &messaging.Message{
		Topic: *fbUid,
		Notification: &messaging.Notification{
			Body:  "hey",
			Title: "theo",
		},
	})

	if err != nil {
		log.Printf("err: %+v", err)
	}
	log.Printf("sucess: '%s'", res)
}
