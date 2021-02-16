package common

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HttpError struct {
	Code int
	Msg  string
}

func JsonIn(r *http.Request, p interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(p); err != nil {
		return fmt.Errorf("could not unmarshal body: %v", err)
	}
	return nil
}

func JsonOut(w http.ResponseWriter, load interface{}) error {
	response, err := json.Marshal(load)
	if err != nil {
		return nil
	}

	_, err = w.Write(response)
	if err != nil {
		return err
	}
	return nil
}

func ServerError(msg string, a ...interface{}) *HttpError {
	return &HttpError{
		Code: http.StatusInternalServerError,
		Msg:  fmt.Sprintf(msg, a...),
	}
}
