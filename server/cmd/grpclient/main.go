package main

import (
	"context"
	"flag"
	"github.com/talkiewalkie/talkiewalkie/cmd/grpclient/hwpb"
	"google.golang.org/grpc"
	"log"
)

//go:generate protoc -I=./ --go_out=./hwpb --go-grpc_out=./hwpb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false helloworld.proto

var addr = flag.String("addr", "grpc.002fa7.net:443", "server address with port")

func main() {
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Printf("could not connect to server: %+v", err)
		return
	}
	log.Printf("connected to '%s'", *addr)

	client := hwpb.NewGreeterClient(conn)
	out, err := client.SayHello(context.Background(), &hwpb.HelloRequest{Name: "HELLO"})
	if err != nil {
		log.Printf("error in request: %+v", err)
		return
	}

	log.Printf("success: '%s'", out.Message)
}
