package main

import (
	"context"
	"flag"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/talkiewalkie/talkiewalkie/common"
	coco "github.com/talkiewalkie/talkiewalkie/grpc"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

var (
	port = flag.String("port", ":8080", "port to run")
	env  = flag.String("env", "dev", "dev|prod")
)

//go:generate sqlboiler psql
//go:generate protoc -I=protos/ --go_out=pb --go-grpc_out=pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false protos/audio_proc.proto
//go:generate protoc -I=protos/ --go_out=pb --go-grpc_out=pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false protos/app.proto

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
	log.SetOutput(os.Stdout)
	flag.Parse()

	var host string
	println(host)
	switch *env {
	case "dev":
		host = "http://localhost:3000"
		if err := godotenv.Load(fmt.Sprintf(".env.%s", *env)); err != nil {
			log.Panicf("could not load env: %v", err)
		}
	case "prod":
		host = "https://web.talkiewalkie.app"
	default:
		log.Panicf("bad env: %s", *env)
	}

	dbUrl := common.DbUrl(
		"talkiewalkie",
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"),
		false)
	common.RunMigrations("./migrations", dbUrl)

	components, err := common.InitComponents()
	if err != nil {
		panic(err)
	}

	boil.DebugMode = true

	server := grpc.NewServer(
		//grpc.KeepaliveParams(keepalive.ServerParameters{Timeout: time.Hour}),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(),
			grpc_ctxtags.StreamServerInterceptor(),
			func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
				log.Printf("[grpc,stream] %q", info.FullMethod)

				err := handler(srv, ss)
				if err != nil {
					log.Printf("[grpc,stream] %q finished with error: %+v", info.FullMethod, err)
				} else {
					log.Printf("[grpc,stream] %q finished", info.FullMethod)
				}

				return err
			},
			//grpc_opentracing.StreamServerInterceptor(),
			//grpc_prometheus.StreamServerInterceptor,
			func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
				if strings.HasPrefix(info.FullMethod, "/grpc.reflection") {
					return handler(srv, ss)
				} else {
					return grpc_auth.StreamServerInterceptor(common.AuthInterceptor(components))(srv, ss, info, handler)
				}
			},
		)),

		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(),
			func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
				log.Printf("[grpc,unary] %q", info.FullMethod)

				res, err := handler(ctx, req)
				if err != nil {
					log.Printf("[grpc,unary] %q finished with error: %+v", info.FullMethod, err)
				} else {
					log.Printf("[grpc,unary] %q finished", info.FullMethod)
				}

				return res, err
			},
			//grpc_opentracing.UnaryServerInterceptor(),
			//grpc_prometheus.UnaryServerInterceptor,
			//grpc_zap.UnaryServerInterceptor(zapLogger),
			func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
				if strings.HasPrefix(info.FullMethod, "/grpc.health.v1.Health/") ||
					strings.HasPrefix(info.FullMethod, "/grpc.reflection") {
					return handler(ctx, req)
				} else {
					return grpc_auth.UnaryServerInterceptor(common.AuthInterceptor(components))(ctx, req, info, handler)
				}
			},
			//grpc_auth.UnaryServerInterceptor(myAuth(components)),
		)),
	)

	us := coco.NewUserService(components)
	pb.RegisterUserServiceServer(server, us)
	cs := coco.NewConversationService(components)
	pb.RegisterConversationServiceServer(server, cs)
	ms := coco.NewMessageService(components)
	pb.RegisterMessageServiceServer(server, ms)
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	if *env != "prod" {
		reflection.Register(server)
	}

	go func() {
		lis, err := net.Listen("tcp", ":8081")
		if err != nil {
			panic(err)
		}

		err = http.Serve(lis, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		if err != nil {
			panic(err)
		}
	}()

	lis, err := net.Listen("tcp", *port)
	if err != nil {
		panic(err)
	}

	log.Printf("grpc server listening to [%s]", *port)
	if err = server.Serve(lis); err != nil {
		panic(err)
	}

}
