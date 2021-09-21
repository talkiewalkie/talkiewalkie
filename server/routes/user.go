package routes

import (
	"database/sql"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/gorilla/mux"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/entities"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"strconv"
	"strings"
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

// --------

type FriendsOutput struct {
	Friends []FriendsOutputGroup `json:"friends"`
	Randoms []string             `json:"randoms"`
}

type FriendsOutputGroup struct {
	Uuid    string   `json:"uuid"`
	Display string   `json:"display"`
	Handles []string `json:"handles"`
}

func Friends(w http.ResponseWriter, r *http.Request) {
	ctx := common.WithAuthedContext(r)

	pageSz := 20

	var offset int
	offsetVal := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetVal)
	if offsetVal != "" && err != nil {
		http.Error(w, "offset param is not an integer", http.StatusBadRequest)
		return
	}

	groups, err := entities.UserConversations(ctx, offset, pageSz)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not fetch user's groups: %+v", err), http.StatusInternalServerError)
		return
	}

	groupOutput := []FriendsOutputGroup{}
	for _, g := range groups {
		handles := []string{}
		for _, ug := range g.R.UserGroups {
			redundant := false
			for _, h := range handles {
				if h == ug.R.User.Handle {
					redundant = true
				}
			}
			if !redundant {
				handles = append(handles, ug.R.User.Handle)
			}
		}

		display := g.Name.String
		if !g.Name.Valid {
			display = strings.Join(handles, ", ")
		}

		groupOutput = append(groupOutput, FriendsOutputGroup{
			Uuid:    g.UUID.String(),
			Display: display,
			Handles: handles,
		})
	}

	// TODO: not really random, fine for now, for postgres randomness see https://stackoverflow.com/a/34990197
	randomUsers, err := models.Users(qm.Limit(10)).All(r.Context(), ctx.Components.Db)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not fetch set of random users: %+v", err), http.StatusInternalServerError)
		return
	}
	randoms := []string{}
	for _, user := range randomUsers {
		randoms = append(randoms, user.Handle)
	}

	out := FriendsOutput{Friends: groupOutput, Randoms: randoms}
	common.JsonOut(w, out)
}
