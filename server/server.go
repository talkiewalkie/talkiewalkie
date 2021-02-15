package main

import (
	"context"
	"flag"
	"log"
	"net/http"

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
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panicf("could not load env: %v", err)
	}

	checkMigrations()
	components := common.InitComponents()

	router := mux.NewRouter()
	router.Use(
		//mux.CORSMethodMiddleware(router),
		Logger,
		WithDbMiddleWare)

	unauthenticated.Setup(router, &components)
	authenticated.Setup(router, &components)

	corsWrapper := handlers.CORS(
		handlers.AllowCredentials(),
		handlers.AllowedHeaders([]string{"Authorization"}),
		handlers.AllowedOrigins([]string{"http://localhost:3000"}))

	log.Printf("listening on port %s", *port)
	if err := http.ListenAndServe(*port, corsWrapper(router)); err != nil {
		log.Printf("could not serve: %v", err)
	}
}

func WithDbMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db, err := sqlx.Connect("postgres", "user=theo dbname=talkiewalkie sslmode=disable")

		if err != nil {
			log.Printf("could not connect to db: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer db.Close()

		ctx := context.WithValue(r.Context(), "db", db)
		newR := r.WithContext(ctx)
		next.ServeHTTP(w, newR)
	})
}
