package unauthenticated

import (
	"encoding/json"
	"log"
	"net/http"
)

func WalksHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		withUnauthContext(w, r, func(c unauthenticatedContext) {
			walks, err := c.WalkRepository.GetAll()
			if err != nil {
				log.Printf("could not fetch walks: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			b, err := json.Marshal(walks)
			w.Write(b)
		})
	}
}
