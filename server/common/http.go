package common

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Json(w http.ResponseWriter, load interface{}) {
	response, err := json.Marshal(load)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not marshal payload: %v", err), http.StatusInternalServerError)
	}
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not write body: %v", err), http.StatusInternalServerError)
	}
}
