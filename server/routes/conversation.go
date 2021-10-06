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
)

// ---------------

type ConversationsOutput struct {
	Conversations []ConversationOutput `json:"conversations"`
}

type ConversationOutput struct {
	Uuid    string   `json:"uuid"`
	Display string   `json:"display"`
	Handles []string `json:"handles"`
}

func Conversations(w http.ResponseWriter, r *http.Request) {
	ctx := common.WithAuthedContext(r)

	pageSz := 20

	var offset int
	offsetVal := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetVal)
	if offsetVal != "" && err != nil {
		http.Error(w, "offset param is not an integer", http.StatusBadRequest)
		return
	}

	conversations, err := entities.UserConversations(ctx, offset, pageSz)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not fetch user's conversations: %+v", err), http.StatusInternalServerError)
		return
	}

	conversationOutput := []ConversationOutput{}
	for _, g := range conversations {
		handles := []string{}
		for _, ug := range g.R.UserConversations {
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

		conversationOutput = append(conversationOutput, ConversationOutput{
			Uuid:    g.UUID.String(),
			Display: display,
			Handles: handles,
		})
	}

	out := ConversationsOutput{Conversations: conversationOutput}
	common.JsonOut(w, out)
}

// -------------

type ConversationByUuidOutput struct {
	Uuid     string                `json:"uuid"`
	Display  string                `json:"display"`
	Handles  []string              `json:"handles"`
	Messages []ConversationMessage `json:"messages"`
}

type ConversationMessage struct {
	AuthorHandle string `json:"authorHandle"`
	Text         string `json:"text"`
	CreatedAt    string `json:"createdAt"`
}

func ConversationByUuid(w http.ResponseWriter, r *http.Request) {
	ctx := common.WithAuthedContext(r)

	vars := mux.Vars(r)
	uuidraw, ok := vars["uuid"]
	if !ok {
		http.Error(w, "expect a conversation uuid", http.StatusBadRequest)
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

	conversation, err := models.Conversations(
		models.ConversationWhere.UUID.EQ(uuid),
		qm.Load(qm.Rels(models.ConversationRels.UserConversations, models.UserConversationRels.User)),
		qm.Load(
			qm.Rels(models.ConversationRels.Messages, models.MessageRels.Author),
			qm.Limit(pageSz), qm.Offset(offset), qm.OrderBy(fmt.Sprintf("%s DESC", models.MessageColumns.CreatedAt))),
	).One(ctx.Context, ctx.Components.Db)
	if errors.Cause(err) == sql.ErrNoRows {
		http.Error(w, fmt.Sprintf("no conversation for '%s'", uuid.String()), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("could not find conversation: %+v", err), http.StatusInternalServerError)
		return
	}
	ok, err = entities.CanAccessConversation(conversation, ctx.User)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "can't access this conversation, not listed in its participants", http.StatusForbidden)
		return
	}

	messages := []ConversationMessage{}
	//for _, msg := range conversation.R.Messages {
	//messages = append(messages, ConversationMessage{
	//	AuthorHandle: msg.R.Author.Handle,
	//	Text:         msg.Text,
	//	CreatedAt:    msg.CreatedAt.Format(time.RFC3339),
	//})
	//}

	handles := []string{}
	for _, userConversation := range conversation.R.UserConversations {
		redundant := false
		for _, h := range handles {
			if h == userConversation.R.User.Handle {
				redundant = true
			}
		}
		if !redundant {
			handles = append(handles, userConversation.R.User.Handle)
		}
	}

	display := conversation.Name.String
	if !conversation.Name.Valid {
		display = strings.Join(handles, ", ")
	}

	out := ConversationByUuidOutput{
		Uuid:     conversation.UUID.String(),
		Display:  display,
		Handles:  handles,
		Messages: messages,
	}
	common.JsonOut(w, out)
}
