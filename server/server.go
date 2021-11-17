package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/talkiewalkie/talkiewalkie/api"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.String("port", ":8080", "port to run")
	env  = flag.String("env", "dev", "dev|prod")
)

//go:generate sqlboiler psql
//go:generate protoc -I=protos/ --go_out=pb --go-grpc_out=pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false protos/audio_proc.proto
//go:generate protoc -I=protos/ --go_out=pb --go-grpc_out=pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false protos/app.proto

//go:generate genny -in=pkg/generics/cache.go -out=repositories/caches/asset.go -pkg caches gen "CacheKey=int,uuid:uuid2.UUID CacheValue=Asset:models.Asset"
//go:generate genny -in=pkg/generics/cache.go -out=repositories/caches/conversation.go -pkg caches gen "CacheKey=int,uuid:uuid2.UUID CacheValue=Conversation:models.Conversation"
//go:generate genny -in=pkg/generics/cache.go -out=repositories/caches/message.go -pkg caches gen "CacheKey=int,uuid:uuid2.UUID CacheValue=Message:models.Message"
//go:generate genny -in=pkg/generics/cache.go -out=repositories/caches/user.go -pkg caches gen "CacheKey=int,uuid:uuid2.UUID,string CacheValue=User:models.User"
//go:generate genny -in=pkg/generics/multicache.go -out=repositories/caches/userconversation.go -pkg caches gen "CacheKey=int CacheValue=UserConversation:models.UserConversation"

//go:generate genny -in=pkg/generics/slice.go -out=pkg/slices/builtins.go -pkg slices gen "ItemType=BUILTINS,uuid:uuid2.UUID"

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
	log.SetOutput(os.Stdout)
	flag.Parse()

	switch *env {
	case "dev":
		boil.DebugMode = true
		if err := godotenv.Load(fmt.Sprintf(".env.%s", *env)); err != nil {
			log.Panicf("could not load env: %v", err)
		}
	case "prod":
		break
	default:
		log.Panicf("bad env: %s", *env)
	}

	common.RunMigrations("./migrations", common.DbUri(
		"talkiewalkie",
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"),
		false))
	components := common.InitComponents()

	server := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			api.ServerStreamAuthInterceptor(components),
			api.ServerStreamLoggerInterceptor,
			api.ServerStreamRequestComponentsInterceptor(components),
		)),

		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			api.UnaryAuthInterceptor(components),
			api.UnaryLoggerInterceptor,
			api.UnaryRequestComponentsInterceptor(components),
		)),
	)

	us := api.NewUserService()
	pb.RegisterUserServiceServer(server, us)
	cs := api.NewConversationService()
	pb.RegisterConversationServiceServer(server, cs)
	es := api.EventService{}
	pb.RegisterEventServiceServer(server, es)

	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	if *env != "prod" {
		reflection.Register(server)
	}

	// starting an http listener for cloud health checks (GKE requires it.)
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
