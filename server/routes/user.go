package routes

import (
	"database/sql"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/gorilla/mux"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
)

// ------------

type UserByHandleOutput struct {
	Handle  string                   `json:"handle"`
	Bio     string                   `json:"bio"`
	Profile string                   `json:"profile"`
	Walks   []UserByHandleWalkOutput `json:"walks"`
	Likes   int                      `json:"likes"`
}

type UserByHandleWalkOutput struct {
	Uuid  string `json:"uuid"`
	Title string `json:"title"`
}

func UserByHandle(w http.ResponseWriter, r *http.Request) {
	ctx := common.WithContext(r)

	vars := mux.Vars(r)
	handle, ok := vars["handle"]
	if !ok {
		http.Error(w, "expect a user handle", http.StatusBadRequest)
		return
	}

	u, err := models.Users(
		models.UserWhere.Handle.EQ(handle),
		qm.Load(models.UserRels.ProfilePictureAsset),
		qm.Load(models.UserRels.AuthorWalks),
		qm.Load(qm.Rels(models.UserRels.UserWalks, models.UserWalkRels.Walk))).One(r.Context(), ctx.Components.Db)
	if errors.Cause(err) == sql.ErrNoRows {
		http.Error(w, fmt.Sprintf("no user for '%s'", handle), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("could not find user: %+v", err), http.StatusInternalServerError)
		return
	}

	var profile string
	if u.R.ProfilePictureAsset != nil {
		profile, err = ctx.Components.Storage.AssetUrl(u.R.ProfilePictureAsset)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not make url for profile picture: %+v", err), http.StatusInternalServerError)
			return
		}
	}

	var walks []UserByHandleWalkOutput
	for _, userWalk := range u.R.AuthorWalks {
		walks = append(walks, UserByHandleWalkOutput{Uuid: userWalk.UUID.String(), Title: userWalk.Title})
	}

	out := UserByHandleOutput{
		Handle:  u.Handle,
		Bio:     u.Bio.String,
		Profile: profile,
		Walks:   walks,
		Likes:   len(u.R.UserWalks),
	}
	common.JsonOut(w, out)
}
