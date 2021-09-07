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
	"strconv"
)

// ------------

type UserByHandleOutput struct {
	Handle     string                   `json:"handle"`
	Bio        string                   `json:"bio"`
	Profile    string                   `json:"profile"`
	TotalWalks int64                    `json:"totalWalks"`
	Walks      []UserByHandleWalkOutput `json:"walks"`
	Likes      int                      `json:"likes"`
}

type UserByHandleWalkOutput struct {
	Uuid  string `json:"uuid"`
	Title string `json:"title"`
}

func UserByHandle(w http.ResponseWriter, r *http.Request) {
	ctx := common.WithContext(r)

	var offset int
	params := r.URL.Query()
	if vals, ok := params["offset"]; ok && len(vals) > 0 {
		value, err := strconv.Atoi(vals[0])
		if err != nil {
			http.Error(w, "bad offset", http.StatusBadRequest)
			return
		}
		offset = value
	}

	vars := mux.Vars(r)
	handle, ok := vars["handle"]
	if !ok {
		http.Error(w, "expect a user handle", http.StatusBadRequest)
		return
	}

	u, err := models.Users(
		models.UserWhere.Handle.EQ(handle),
		qm.Load(models.UserRels.ProfilePictureAsset),
		qm.Load(models.UserRels.AuthorWalks, qm.Limit(20), qm.Offset(offset)),
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

	count, err := models.Walks(models.WalkWhere.AuthorID.EQ(u.ID)).Count(r.Context(), ctx.Components.Db)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not count the total number of walks: %+v", err), http.StatusInternalServerError)
		return
	}

	out := UserByHandleOutput{
		Handle:     u.Handle,
		Bio:        u.Bio.String,
		Profile:    profile,
		Walks:      walks,
		TotalWalks: count,
		Likes:      len(u.R.UserWalks),
	}
	common.JsonOut(w, out)
}

// --------

type MeOutput struct {
	Handle string `json:"handle"`
}

func Me(w http.ResponseWriter, r *http.Request) {
	ctx := common.WithAuthedContext(r)

	out := MeOutput{Handle: ctx.User.Handle}
	common.JsonOut(w, out)
}
