package unauthenticated

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"firebase.google.com/go/auth"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/repository"
)

type signInPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignInHandler(components *common.Components) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		withUnauthContext(w, r, func(c unauthenticatedContext) {
			var p signInPayload
			if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
				log.Printf("could not decode post: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			u, err := components.FirebaseAuth.CreateUser(r.Context(), (&auth.UserToCreate{}).Email(p.Email).Password(p.Password))
			if err != nil {
				log.Printf("could not create a firebase user: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if err = components.EmailClient.SendEmail(p.Email, []byte("verifie ton compte gros")); err != nil {
				log.Printf("failed to send verification email: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			code, err := components.FirebaseAuth.EmailSignInLink(r.Context(), p.Email, &auth.ActionCodeSettings{URL: "http://localhost:8080/verify", HandleCodeInApp: true})
			if err != nil {
				log.Printf("could not send email: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			dbU, err := c.UserRepository.CreateUser(p.Email, u.UID, code)
			if err != nil {
				log.Printf("could not create user in db: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			response, err := json.Marshal(dbU)
			_, err = w.Write(response)
			if err != nil {
				log.Printf("could not write out: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		})
	}
}

type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(components *common.Components) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		withUnauthContext(w, r, func(c unauthenticatedContext) {
			var p loginPayload
			if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
				log.Printf("could not decode post: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			fbU, err := components.FirebaseAuth.GetUserByEmail(r.Context(), p.Email)
			if err != nil {
				log.Printf("could not find user: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			emptyChan := make(chan int, 1)
			userChan := make(chan *repository.User, 1)
			go func(w http.ResponseWriter) {
				token, err := components.FirebaseAuth.CustomTokenWithClaims(r.Context(), fbU.UID, map[string]interface{}{})
				if err != nil {
					log.Printf("could not create auth token: %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				http.SetCookie(w, &http.Cookie{Name: "fbToken", Value: token, Path: "/", Expires: time.Now().Add(365 * 24 * time.Hour)})
				emptyChan <- 1
			}(w)

			go func() {
				u, err := c.UserRepository.GetUserByUid(fbU.UID)
				if err != nil {
					log.Printf("did not find the user in db: %v", err)
				}
				userChan <- u
			}()

			response, err := json.Marshal(<-userChan)
			<-emptyChan
			if _, err = w.Write(response); err != nil {
				log.Printf("error marshalling response: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		})
	}
}
