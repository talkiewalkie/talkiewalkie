package main

import (
	"context"
	"crypto/tls"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"time"
	"toto/grpc/hwpb"
)

//go:generate protoc -I=./ --go_out=./hwpb --go-grpc_out=./hwpb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false helloworld.proto

var addr = flag.String("addr", "grpc.002fa7.net:443", "server address with port")

func main() {
	flag.Parse()

	connCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(connCtx,
		*addr,
		// with hosted containers we can enable security like so - unclear how it's handling TLS but
		// we're getting through somehow.
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
		// with ngrok we should do insecure loads - let's see how that translate in swift-grpc
		//grpc.WithInsecure(),
		grpc.WithBlock(),
	)
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
