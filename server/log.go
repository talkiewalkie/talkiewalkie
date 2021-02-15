package main

import (
	"log"
	"net/http"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("GET '%s'", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
