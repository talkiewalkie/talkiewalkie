package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

// RecoverMiddleWare from https://stackoverflow.com/a/28746725
func RecoverMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			r := recover()
			if r != nil {
				var err error
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}
				log.Println(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
