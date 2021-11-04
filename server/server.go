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
	"github.com/talkiewalkie/talkiewalkie/common"
	coco "github.com/talkiewalkie/talkiewalkie/grpc"
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

//go:generate genny -in=entities/genericcache/cache.go -out=entities/asset_cache.go -pkg entities gen "CacheKey=int,uuid2.UUID CacheValue=Asset"
//go:generate genny -in=entities/genericcache/cache.go -out=entities/conversation_cache.go -pkg entities gen "CacheKey=int,uuid2.UUID CacheValue=Conversation"
//go:generate genny -in=entities/genericcache/cache.go -out=entities/message_cache.go -pkg entities gen "CacheKey=int,uuid2.UUID CacheValue=Message"
//go:generate genny -in=entities/genericcache/cache.go -out=entities/user_cache.go -pkg entities gen "CacheKey=int,uuid2.UUID,string CacheValue=User"
//go:generate genny -in=entities/genericmulticache/cache.go -out=entities/userconversation_cache.go -pkg entities gen "CacheKey=int CacheValue=UserConversation"

//go:generate genny -in=pkg/slices/generic_slice.go -out=entities/asset_slice.go -pkg entities gen "ItemType=Asset"
//go:generate genny -in=pkg/slices/generic_slice.go -out=entities/conversation_slice.go -pkg entities gen "ItemType=Conversation"
//go:generate genny -in=pkg/slices/generic_slice.go -out=entities/message_slice.go -pkg entities gen "ItemType=Message"
//go:generate genny -in=pkg/slices/generic_slice.go -out=entities/user_slice.go -pkg entities gen "ItemType=User"
//go:generate genny -in=pkg/slices/generic_slice.go -out=entities/userconversation_slice.go -pkg entities gen "ItemType=UserConversation"

//go:generate genny -in=pkg/slices/generic_slicemap.go -out=entities/asset_slicemap.go -pkg entities gen "ItemType=Asset MapTarget=BUILTINS"
//go:generate genny -in=pkg/slices/generic_slicemap.go -out=entities/conversation_slicemap.go -pkg entities gen "ItemType=Conversation MapTarget=BUILTINS"
//go:generate genny -in=pkg/slices/generic_slicemap.go -out=entities/message_slicemap.go -pkg entities gen "ItemType=Message MapTarget=BUILTINS"
//go:generate genny -in=pkg/slices/generic_slicemap.go -out=entities/user_slicemap.go -pkg entities gen "ItemType=User MapTarget=BUILTINS"
//go:generate genny -in=pkg/slices/generic_slicemap.go -out=entities/userconversation_slicemap.go -pkg entities gen "ItemType=UserConversation MapTarget=BUILTINS"

//go:generate genny -in=pkg/slices/generic_slice.go -out=pkg/slices/builtins.go -pkg slices gen "ItemType=BUILTINS,uuid2.UUID"
//go:generate genny -in=pkg/slices/generic_slicemap.go -out=pkg/slices/builtins_map.go -pkg slices gen "ItemType=BUILTINS MapTarget=BUILTINS"

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

	dbUrl := common.DbUri(
		"talkiewalkie",
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"),
		false)
	common.RunMigrations("./migrations", dbUrl)

	components, err := common.InitComponents()
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			coco.ServerStreamAuthInterceptor(components),
			coco.ServerStreamLoggerInterceptor,
			coco.ServerStreamRequestComponentsInterceptor(components),
		)),

		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			coco.UnaryAuthInterceptor(components),
			coco.UnaryLoggerInterceptor,
			coco.UnaryRequestComponentsInterceptor(components),
		)),
	)

	us := coco.NewUserService()
	pb.RegisterUserServiceServer(server, us)
	cs := coco.NewConversationService()
	pb.RegisterConversationServiceServer(server, cs)
	ms := coco.NewMessageService()
	pb.RegisterMessageServiceServer(server, ms)

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
