package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/talkiewalkie/talkiewalkie/routes"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/talkiewalkie/talkiewalkie/common"
)

var (
	port = flag.String("port", ":8080", "port to run")
	env  = flag.String("env", "dev", "dev|prod")
)

//go:generate sqlboiler psql
//go:generate protoc -I=protos/ --go_out=pb --go-grpc_out=pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false protos/audio_proc.proto

func main() {
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
		func(next http.Handler) http.Handler {
			return handlers.CombinedLoggingHandler(os.Stdout, next)
		},
		common.WithContextMiddleWare(components),
		common.RecoverMiddleWare)

	router.HandleFunc("/walks", routes.Walks).Methods(http.MethodGet)
	router.HandleFunc("/walk/{uuid}", routes.WalkByUuid).Methods(http.MethodGet)
	router.HandleFunc("/walk", routes.CreateWalk).Methods(http.MethodPost)

	router.HandleFunc("/user/{handle}", routes.UserByHandle).Methods(http.MethodGet)

	router.HandleFunc("/asset", routes.UploadHandler).Methods(http.MethodPost)

	router.HandleFunc("/ws", ws)

	corsWrapper := handlers.CORS(
		handlers.AllowCredentials(),
		handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
		handlers.AllowedOrigins([]string{host}),
		handlers.AllowedMethods([]string{http.MethodPost, http.MethodGet, http.MethodOptions, http.MethodHead}))

	log.Printf("listening on port %s", *port)
	if err := http.ListenAndServe(*port, corsWrapper(router)); err != nil {
		log.Printf("could not serve: %v", err)
	}
}
