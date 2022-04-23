package main

import (
	"log"

	uuid2 "github.com/satori/go.uuid"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/talkiewalkie/talkiewalkie/pb"
)

func main() {
	u := pb.User{
		DisplayName:   "toto",
		Uuid:          uuid2.NewV4().String(),
		Conversations: nil,
		Phone:         "+3368594093239",
	}

	str, err := protojson.Marshal(&u)
	if err != nil {
		panic(err)
	}

	log.Println(string(str))
	umu := pb.User{}
	if err := protojson.Unmarshal([]byte(str), &umu); err != nil {
		panic(err)
	}

	log.Println("success!")
}
