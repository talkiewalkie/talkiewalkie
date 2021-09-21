package routes

import (
	"fmt"
	uuid2 "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"reflect"
	"sort"
	"sync"
)

// ---------------

type MessageInput struct {
	// Clients should only specify one the options below - if groupUuid is set it will prevail over the list of handles
	GroupUuid string   `json:"groupUuid"`
	Handles   []string `json:"handles"`

	Text string `json:"text"`
}

func Message(w http.ResponseWriter, r *http.Request) {
	ctx := common.WithAuthedContext(r)

	var msg MessageInput
	if err := common.JsonIn(r, &msg); err != nil {
		http.Error(w, "could not parse input", http.StatusBadRequest)
		return
	}

	if msg.GroupUuid != "" {
		uuid, err := uuid2.FromString(msg.GroupUuid)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not parse uuid: %+v", err), http.StatusInternalServerError)
			return
		}
		group, err := models.Groups(models.GroupWhere.UUID.EQ(uuid)).One(r.Context(), ctx.Components.Db)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not find group: %+v", err), http.StatusInternalServerError)
			return
		}
		message := models.Message{
			Text:     msg.Text,
			AuthorID: null.IntFrom(ctx.User.ID),
			GroupID:  null.IntFrom(group.ID),
		}

		if err = message.Insert(r.Context(), ctx.Components.Db, boil.Infer()); err != nil {
			http.Error(w, fmt.Sprintf("could not insert message: %+v", err), http.StatusInternalServerError)
			return
		}
		return
	}

	if len(msg.Handles) == 0 {
		http.Error(w, "message needs a recipient", http.StatusBadRequest)
		return
	}

	handles := make(map[string]int)
	for _, handle := range msg.Handles {
		handles[handle] += 1
	}

	uniqueHandles := []string{}
	for handle, _ := range handles {
		uniqueHandles = append(uniqueHandles, handle)
	}

	recipients, err := models.Users(models.UserWhere.Handle.IN(uniqueHandles)).All(r.Context(), ctx.Components.Db)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not find recipients: %+v", err), http.StatusBadRequest)
		return
	}
	if len(recipients) != len(msg.Handles) {
		http.Error(w, "some users where not found", http.StatusBadRequest)
		return
	}

	ids := []int{ctx.User.ID}
	for _, recipient := range recipients {
		ids = append(ids, recipient.ID)
	}

	ugs, err := models.UserGroups(
		models.UserGroupWhere.UserID.EQ(ctx.User.ID),
		qm.Load(qm.Rels(models.UserGroupRels.Group, models.GroupRels.UserGroups)),
	).All(r.Context(), ctx.Components.Db)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not find recipients groups: %+v", err), http.StatusInternalServerError)
		return
	}

	groupId := null.NewInt(0, false)
	sort.Ints(ids)
	for _, ug := range ugs {
		groupIds := []int{}
		for _, ug := range ug.R.Group.R.UserGroups {
			// TODO: somehow traversing the dependencies brings redundant rows, e.g. the list we're iterating on can
			// 		yield [115, 115, 116] as user ids.
			redundant := false
			for _, id := range groupIds {
				if ug.UserID == id {
					redundant = true
				}
			}
			if !redundant {
				groupIds = append(groupIds, ug.UserID)
			}
		}
		sort.Ints(groupIds)
		if reflect.DeepEqual(groupIds, ids) {
			groupId = null.IntFrom(ug.GroupID)
			break
		}
	}

	// TODO: use a batch insert method like COPY which would make things faster
	tx, err := ctx.Components.Db.BeginTx(r.Context(), nil)
	if err != nil {
		tx.Rollback()
		http.Error(w, fmt.Sprintf("could not start transaction: %+v", err), http.StatusInternalServerError)
		return
	}
	if !groupId.Valid {
		newGroup := models.Group{}
		if err = newGroup.Insert(r.Context(), tx, boil.Infer()); err != nil {
			tx.Rollback()
			http.Error(w, fmt.Sprintf("could not create new group: %+v", err), http.StatusInternalServerError)
			return
		}

		errs := make(chan error, 1)
		var wg sync.WaitGroup
		for _, id := range ids {
			wg.Add(1)
			uid := id
			go func() {
				ug := models.UserGroup{
					UserID:  uid,
					GroupID: newGroup.ID,
				}
				if err := ug.Insert(r.Context(), tx, boil.Infer()); err != nil {
					errs <- err
				}
				wg.Done()
			}()
		}
		wg.Wait()
		close(errs)
		for err := range errs {
			tx.Rollback()
			http.Error(w, fmt.Sprintf("could not add recipient to new group: %+v", err), http.StatusInternalServerError)
			return
		}
		groupId = null.IntFrom(newGroup.ID)
	}

	message := models.Message{
		Text:     msg.Text,
		AuthorID: null.IntFrom(ctx.User.ID),
		GroupID:  groupId,
	}
	if err = message.Insert(r.Context(), tx, boil.Infer()); err != nil {
		tx.Rollback()
		http.Error(w, fmt.Sprintf("could not insert message: %+v", err), http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		http.Error(w, fmt.Sprintf("could not commit transaction: %+v", err), http.StatusInternalServerError)
		return
	}
}
