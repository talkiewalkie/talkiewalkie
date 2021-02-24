package common

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

	"github.com/talkiewalkie/talkiewalkie/pb"
)

func grpcLogger(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log.Printf("grpc request: %s", method)
	err := invoker(ctx, method, req, reply, cc, opts...)
	return err
}

func NewAudioClient() (pb.CompressionClient, error) {
	conn, err := grpc.Dial(os.Getenv("AUDIO_SERVICE_URL"), grpc.WithInsecure(), grpc.WithChainUnaryInterceptor(grpcLogger))
	if err != nil {
		return nil, err
	}

	// poll status, if connection was not established after a second we consider the service unavailable
	tick := time.NewTicker(100 * time.Millisecond)
	timeout := time.After(1 * time.Second)
	func() {
		for {
			select {
			case <-tick.C:
				if conn.GetState() == connectivity.Ready {
					return
				}
			case <-timeout:
				return
			}
		}
	}()

	if conn.GetState() != connectivity.Ready {
		return nil, fmt.Errorf("could not reach audio service at '%s', status is '%v'",
			os.Getenv("AUDIO_SERVICE_URL"),
			conn.GetState())
	}
	client := pb.NewCompressionClient(conn)
	return client, nil
}
