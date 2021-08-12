package common

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"net/http"
	"strconv"
)

type Context struct {
	Components *Components
	User       *models.User
}

type FirebaseJWT struct {
	IdToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
}

type FirebaseClaim struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func WithContextMiddleWare(comps *Components) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var u *models.User

			tokenb64, err := r.Cookie("TalkieWalkie.AuthUserTokens")
			if err == nil {
				quotedStr, err := base64.StdEncoding.DecodeString(tokenb64.Value)
				if err != nil {
					http.Error(w, fmt.Sprintf("could not b64 decode cookie: %+v", err), http.StatusBadRequest)
					return
				}

				jwtStr, err := strconv.Unquote(string(quotedStr))
				if err != nil {
					http.Error(w, "failed to unquote json string", http.StatusBadRequest)
					return
				}

				var fbToken FirebaseJWT
				if err = json.Unmarshal([]byte(jwtStr), &fbToken); err != nil {
					http.Error(w, fmt.Sprintf("could not deserialize jwt: %+v", err), http.StatusBadRequest)
					return
				}

				tok, err := comps.FbAuth.VerifyIDTokenAndCheckRevoked(r.Context(), fbToken.IdToken)
				if err != nil {
					http.Error(w, fmt.Sprintf("auth cookie provided but couldn't be verified: %+v", err), http.StatusBadRequest)
					return
				}
				u, err = models.Users(models.UserWhere.FirebaseUID.EQ(null.NewString(tok.UID, true))).One(r.Context(), comps.Db)
				if errors.Cause(err) == sql.ErrNoRows {
					var handle, picture string
					if name, ok := tok.Claims["name"]; ok {
						handle = slug.Make(name.(string))
					}
					if email, ok := tok.Claims["email"]; ok && handle == "" {
						handle = slug.Make(email.(string))
					}
					if url, ok := tok.Claims["picture"]; ok {
						picture = url.(string)
					}

					fmt.Printf("%s %s", handle, picture)
					u = &models.User{
						Handle:         handle,
						FirebaseUID:    null.NewString(tok.UID, true),
						ProfilePicture: null.NewInt(0, false), // TODO reupload picture
					}
					if err = u.Insert(r.Context(), comps.Db, boil.Infer()); err != nil {
						http.Error(w, fmt.Sprintf("could not create matching db user for new firebase user: %+v", err), http.StatusInternalServerError)
						return
					}
				} else if err != nil {
					http.Error(w, fmt.Sprintf("failed to query for user uid: %+v", err), http.StatusInternalServerError)
					return
				}
			}

			myCtx := Context{
				Components: comps,
				User:       u,
			}
			ctx := context.WithValue(r.Context(), "context", myCtx)
			newR := r.WithContext(ctx)
			next.ServeHTTP(w, newR)
		})
	}
}

func WithContext(r *http.Request) Context {
	ctx := r.Context()

	services, ok := ctx.Value("context").(Context)
	if !ok {
		panic("failed to get services from context")
	}

	return services
}

func WithAuthedContext(r *http.Request) Context {
	ctx := r.Context()

	services, ok := ctx.Value("context").(Context)
	if !ok {
		panic("failed to get services from context")
	}

	if services.User == nil {
		panic(errors.New("auth"))
	}

	return services
}
