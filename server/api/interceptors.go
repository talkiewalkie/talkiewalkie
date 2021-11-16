package api

import (
	"context"
	"log"
	"strings"

	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ServerStreamAuthInterceptor(components *common.Components) func(interface{}, grpc.ServerStream, *grpc.StreamServerInfo, grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if strings.HasPrefix(info.FullMethod, "/grpc.reflection") {
			return handler(srv, ss)
		} else {
			return grpcauth.StreamServerInterceptor(common.AuthInterceptor(components))(srv, ss, info, handler)
		}
	}
}

func ServerStreamLoggerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	panicked := true
	defer func() {
		if recErr := recover(); recErr != nil || panicked {
			err = status.Errorf(codes.Internal, "%+v", recErr)
			log.Printf("[grpc,stream] %q panicked: %+v", info.FullMethod, err)
		}
	}()

	log.Printf("[grpc,stream] %q", info.FullMethod)

	err = handler(srv, ss)
	if err != nil {
		log.Printf("[grpc,stream] %q finished with error: %+v", info.FullMethod, err)
	} else {
		log.Printf("[grpc,stream] %q finished", info.FullMethod)
	}

	panicked = false
	return err
}

func UnaryLoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	panicked := true
	defer func() {
		if recErr := recover(); recErr != nil || panicked {
			err = status.Errorf(codes.Internal, "%+v", recErr)
			log.Printf("[grpc,unary] %q panicked: %+v", info.FullMethod, err)
		}
	}()

	log.Printf("[grpc,unary] %q", info.FullMethod)

	res, err := handler(ctx, req)
	if err != nil {
		log.Printf("[grpc,unary] %q finished with error: %+v", info.FullMethod, err)
	} else {
		log.Printf("[grpc,unary] %q finished", info.FullMethod)
	}

	panicked = false
	return res, err
}

func UnaryRequestComponentsInterceptor(components *common.Components) func(context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler) (resp interface{}, err error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		newCtx := context.WithValue(ctx, "components", components)
		components.ResetEntityStores(newCtx)

		return handler(newCtx, req)
	}
}

func ServerStreamRequestComponentsInterceptor(components *common.Components) func(interface{}, grpc.ServerStream, *grpc.StreamServerInfo, grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		newCtx := context.WithValue(ss.Context(), "components", components)
		components.ResetEntityStores(newCtx)

		return handler(newCtx, ss)
	}
}

func UnaryAuthInterceptor(components *common.Components) func(context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler) (resp interface{}, err error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if strings.HasPrefix(info.FullMethod, "/grpc.health.v1.Health/") ||
			strings.HasPrefix(info.FullMethod, "/grpc.reflection") {
			return handler(ctx, req)
		} else {
			return grpcauth.UnaryServerInterceptor(common.AuthInterceptor(components))(ctx, req, info, handler)
		}
	}
}

func WithAuthedContext(ctx context.Context) (*common.Components, *models.User, error) {
	components, ok := ctx.Value("components").(*common.Components)
	if !ok {
		return nil, nil, status.Error(codes.Internal, "[components] not found in context")
	}

	me, ok := ctx.Value("me").(*models.User)
	if !ok {
		return nil, nil, status.Error(codes.Internal, "[me] not found in context")
	}

	return components, me, nil
}
