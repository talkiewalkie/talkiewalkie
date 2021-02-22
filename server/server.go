package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/talkiewalkie/talkiewalkie/authenticated"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/unauthenticated"
)

var (
	port = flag.String("port", ":8080", "port to run")
	env  = flag.String("env", "dev", "dev|prod")
)

//go:generate sqlboiler psql

func main() {
	flag.Parse()

	// todo: should remove prod env file, kubernetes secrets do the job or service accounts will
	err := godotenv.Load(fmt.Sprintf(".env.%s", *env))
	if err != nil {
		log.Panicf("could not load env: %v", err)
	}

	var host string
	switch *env {
	case "dev":
		host = "http://localhost:3000"
	case "prod":
		host = "https://talkiewalkie.app"
	}

	checkMigrations()
	components := common.InitComponents()

	router := mux.NewRouter()
	router.Use(
		//mux.CORSMethodMiddleware(router),
		func(next http.Handler) http.Handler {
			return handlers.CombinedLoggingHandler(os.Stdout, next)
		},
		WithDbMiddleWare)

	unauthenticated.Setup(router, &components)
	authenticated.Setup(router, &components)

	corsWrapper := handlers.CORS(
		handlers.AllowCredentials(),
		handlers.AllowedHeaders([]string{"Authorization"}),
		handlers.AllowedOrigins([]string{host}))

	log.Printf("listening on port %s", *port)
	if err := http.ListenAndServe(*port, corsWrapper(router)); err != nil {
		log.Printf("could not serve: %v", err)
	}
}

func WithDbMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dsName := fmt.Sprintf(
			"postgres://%s:%s@%s/talkiewalkie?sslmode=disable",
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_HOST"),
		)
		db, err := sqlx.Connect("postgres", dsName)

		if err != nil {
			http.Error(w, fmt.Sprintf("could not connect to db: %v", err), http.StatusInternalServerError)
			return
		}

		defer db.Close()

		ctx := context.WithValue(r.Context(), "db", db)
		newR := r.WithContext(ctx)
		next.ServeHTTP(w, newR)
	})
}
