package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"toto/grpc/hwpb"
)

//go:generate protoc -I=./ --go_out=./hwpb --go-grpc_out=./hwpb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false helloworld.proto

var port = flag.String("port", ":50051", "port")

type greeterService struct {
}

func (g greeterService) SayHello(ctx context.Context, request *hwpb.HelloRequest) (*hwpb.HelloReply, error) {
	return &hwpb.HelloReply{Message: fmt.Sprintf("hello %s", request.Name)}, nil
}

var _ hwpb.GreeterServer = greeterService{}

func main() {
	flag.Parse()

	s := grpc.NewServer(grpc.UnaryInterceptor(
		func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			res, err := handler(ctx, req)
			if err != nil {
				log.Printf("%s ended with error: %+v", info.FullMethod, err)
			} else {
				log.Printf("%s succeeded", info.FullMethod)
			}

			return res, err
		}))
	hwpb.RegisterGreeterServer(s, &greeterService{})

	lis, err := net.Listen("tcp", *port)
	if err != nil {
		panic(err)
	}

	log.Printf("grpc server listening to [%s]", *port)
	if err = s.Serve(lis); err != nil {
		panic(err)
	}
}
