package common

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func JsonIn(w http.ResponseWriter, r *http.Request, p interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(p); err != nil {
		Error(w, fmt.Sprintf("could not decode post: %v", err), http.StatusInternalServerError)
		return err
	}
	return nil
}

func JsonOut(w http.ResponseWriter, load interface{}) {
	response, err := json.Marshal(load)
	if err != nil {
		Error(w, fmt.Sprintf("could not marshal payload: %v", err), http.StatusInternalServerError)
	}
	_, err = w.Write(response)
	if err != nil {
		Error(w, fmt.Sprintf("could not write body: %v", err), http.StatusInternalServerError)
	}
}

func Error(w http.ResponseWriter, msg string, code int) {
	log.Println(msg)
	http.Error(w, msg, code)
}
