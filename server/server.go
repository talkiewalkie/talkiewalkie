package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	errors2 "github.com/friendsofgo/errors"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/soheilhy/cmux"
	coco "github.com/talkiewalkie/talkiewalkie/grpc"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/talkiewalkie/talkiewalkie/common"
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
	flag.Parse()

	var host string
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
		os.Getenv("POSTGRES_HOST"), "5432",
		false)
	common.RunMigrations("./migrations", dbUrl)

	components, err := common.InitComponents()
	if err != nil {
		panic(err)
	}

	boil.DebugMode = true

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.Use(
		// No logging for root
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.String() == "/" {
					next.ServeHTTP(w, r)
				} else {
					handlers.CombinedLoggingHandler(os.Stdout, next).ServeHTTP(w, r)
				}
			})
		},
		common.WithContextMiddleWare(components),
		common.RecoverMiddleWare)

	server := grpc.NewServer(
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
			grpc_auth.StreamServerInterceptor(myAuth(components)),
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
			grpc_auth.UnaryServerInterceptor(myAuth(components)),
		)),
	)

	us := coco.NewUserService(components)
	pb.RegisterUserServiceServer(server, us)
	cs := coco.NewConversationService(components)
	pb.RegisterConversationServiceServer(server, cs)
	ms := coco.NewMessageService(components)
	pb.RegisterMessageServiceServer(server, ms)

	lis, err := net.Listen("tcp", *port)
	if err != nil {
		panic(err)
	}

	mux := cmux.New(lis)
	grpcLis := mux.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpLis := mux.Match(cmux.HTTP1Fast())

	go func() {
		corsWrapper := handlers.CORS(
			handlers.AllowCredentials(),
			// TODO: The sentry-trace header is sent by the web client on initial calls for some reason. It's a bit strange
			//    and should be investigated - but I wasted enough time on strange CORS errors for the day.
			handlers.AllowedHeaders([]string{"Authorization", "Content-Type", "sentry-trace", "X-TalkieWalkie-Auth"}),
			handlers.AllowedOrigins([]string{host}),
			handlers.AllowedMethods([]string{http.MethodPost, http.MethodGet, http.MethodOptions, http.MethodHead}))

		s := &http.Server{Handler: corsWrapper(router)}
		log.Printf("http server listening to [%s]", *port)
		if err := s.Serve(httpLis); err != nil {
			panic(err)
		}
	}()

	log.Printf("grpc server listening to [%s]", *port)
	if err = server.Serve(grpcLis); err != nil {
		panic(err)
	}
}

func myAuth(c *common.Components) func(ctx context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("failed to get call metadata")
		}
		jwts := md.Get("Authorization")
		if len(jwts) != 1 {
			return nil, status.Error(codes.PermissionDenied, "missing authorization metadata key")
		}

		tok, err := c.FbAuth.VerifyIDTokenAndCheckRevoked(ctx, strings.Replace(jwts[0], "Bearer ", "", 1))
		if err != nil {
			return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("auth header provided couldn't be verified: %+v", err))
		}

		u, err := models.Users(models.UserWhere.FirebaseUID.EQ(null.StringFrom(tok.UID))).One(ctx, c.Db)
		if err != nil && errors2.Cause(err) == sql.ErrNoRows {
			var handle, picture string
			if name, ok := tok.Claims["name"]; ok {
				handle = slug.Make(name.(string))
			}
			if email, ok := tok.Claims["email"]; ok && handle == "" {
				handle = slug.Make(email.(string))
			}
			if url, ok := tok.Claims["picture"]; ok {
				picture = url.(string)
			}

			fmt.Printf("%s %s", handle, picture)
			u = &models.User{
				Handle:         handle,
				FirebaseUID:    null.NewString(tok.UID, true),
				ProfilePicture: null.NewInt(0, false), // TODO reupload picture
			}
			if err = u.Insert(ctx, c.Db, boil.Infer()); err != nil {
				return nil, status.Error(codes.Internal, fmt.Sprintf("could not create matching db user for new firebase user: %+v", err))
			}
		} else if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to query for user uid: %+v", err))
		}

		newCtx := context.WithValue(ctx, "user", u)
		return newCtx, nil
	}
}
