package routes

import (
	"database/sql"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/gorilla/mux"
	uuid2 "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/entities"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ---------------

type GroupsOutput struct {
	Groups []GroupOutput `json:"groups"`
}

type GroupOutput struct {
	Uuid    string   `json:"uuid"`
	Display string   `json:"display"`
	Handles []string `json:"handles"`
}

func Groups(w http.ResponseWriter, r *http.Request) {
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

	groupOutput := []GroupOutput{}
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

		groupOutput = append(groupOutput, GroupOutput{
			Uuid:    g.UUID.String(),
			Display: display,
			Handles: handles,
		})
	}

	out := GroupsOutput{Groups: groupOutput}
	common.JsonOut(w, out)
}

// -------------

type GroupByUuidOutput struct {
	Uuid     string         `json:"uuid"`
	Display  string         `json:"display"`
	Handles  []string       `json:"handles"`
	Messages []GroupMessage `json:"messages"`
}

type GroupMessage struct {
	AuthorHandle string `json:"authorHandle"`
	Text         string `json:"text"`
	CreatedAt    string `json:"createdAt"`
}

func GroupByUuid(w http.ResponseWriter, r *http.Request) {
	ctx := common.WithAuthedContext(r)

	vars := mux.Vars(r)
	uuidraw, ok := vars["uuid"]
	if !ok {
		http.Error(w, "expect a group uuid", http.StatusBadRequest)
		return
	}

	pageSz := 20

	var offset int
	offsetVal := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetVal)
	if offsetVal != "" && err != nil {
		http.Error(w, "offset param is not an integer", http.StatusBadRequest)
		return
	}

	uuid, err := uuid2.FromString(uuidraw)
	if err != nil {
		http.Error(w, fmt.Sprintf("uuid malformed: %v", err), http.StatusBadRequest)
		return
	}

	group, err := models.Groups(
		models.GroupWhere.UUID.EQ(uuid),
		qm.Load(qm.Rels(models.GroupRels.UserGroups, models.UserGroupRels.User)),
		qm.Load(
			qm.Rels(models.GroupRels.Messages, models.MessageRels.Author),
			qm.Limit(pageSz), qm.Offset(offset), qm.OrderBy(fmt.Sprintf("%s DESC", models.MessageColumns.CreatedAt))),
	).One(ctx.Context, ctx.Components.Db)
	if errors.Cause(err) == sql.ErrNoRows {
		http.Error(w, fmt.Sprintf("no group for '%s'", uuid.String()), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("could not find group: %+v", err), http.StatusInternalServerError)
		return
	}

	messages := []GroupMessage{}
	for _, msg := range group.R.Messages {
		messages = append(messages, GroupMessage{
			AuthorHandle: msg.R.Author.Handle,
			Text:         msg.Text,
			CreatedAt:    msg.CreatedAt.Format(time.RFC3339),
		})
	}

	handles := []string{}
	for _, userGroup := range group.R.UserGroups {
		redundant := false
		for _, h := range handles {
			if h == userGroup.R.User.Handle {
				redundant = true
			}
		}
		if !redundant {
			handles = append(handles, userGroup.R.User.Handle)
		}
	}

	display := group.Name.String
	if !group.Name.Valid {
		display = strings.Join(handles, ", ")
	}

	out := GroupByUuidOutput{
		Uuid:     group.UUID.String(),
		Display:  display,
		Handles:  handles,
		Messages: messages,
	}
	common.JsonOut(w, out)
}
