package common

import (
	"context"
	"log"
	"os"

	"google.golang.org/grpc"

	"github.com/talkiewalkie/talkiewalkie/pb"
)

func grpcLogger(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log.Printf("grpc request: %s", method)
	return nil
}

func NewAudioClient() (pb.CompressionClient, error) {
	conn, err := grpc.Dial(os.Getenv("AUDIO_SERVICE_URL"), grpc.WithInsecure(), grpc.WithChainUnaryInterceptor(grpcLogger))
	if err != nil {
		return nil, err
	}
	client := pb.NewCompressionClient(conn)
	return client, nil
}
