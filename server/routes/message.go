package routes

import (
	"fmt"
	uuid2 "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/entities"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"log"
	"net/http"
	"reflect"
	"sort"
	"sync"
	"time"
)

// ---------------

type MessageInput struct {
	// Clients should only specify one the options below - if conversationUuid is set it will prevail over the list of handles
	ConversationUuid string   `json:"conversationUuid"`
	Handles          []string `json:"handles"`

	Text string `json:"text"`
}

func Message(w http.ResponseWriter, r *http.Request) {
	ctx := common.WithAuthedContext(r)

	var msg MessageInput
	if err := common.JsonIn(r, &msg); err != nil {
		http.Error(w, "could not parse input", http.StatusBadRequest)
		return
	}

	if msg.ConversationUuid != "" {
		uuid, err := uuid2.FromString(msg.ConversationUuid)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not parse uuid: %+v", err), http.StatusInternalServerError)
			return
		}
		conversation, err := models.Conversations(models.ConversationWhere.UUID.EQ(uuid)).One(r.Context(), ctx.Components.Db)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not find conversation: %+v", err), http.StatusInternalServerError)
			return
		}
		message := models.Message{
			//Text:           msg.Text,
			AuthorID:       null.IntFrom(ctx.User.ID),
			ConversationID: conversation.ID,
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

	uniqueHandles := []string{ctx.User.Handle}
	for handle, _ := range handles {
		if handle != ctx.User.Handle {
			uniqueHandles = append(uniqueHandles, handle)
		}
	}

	recipients, err := models.Users(models.UserWhere.Handle.IN(uniqueHandles)).All(r.Context(), ctx.Components.Db)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not find recipients: %+v", err), http.StatusBadRequest)
		return
	}
	if len(recipients) != len(uniqueHandles) {
		http.Error(w, "some users where not found", http.StatusBadRequest)
		return
	}

	ids := []int{ctx.User.ID}
	for _, recipient := range recipients {
		if recipient.ID != ctx.User.ID {
			ids = append(ids, recipient.ID)
		}
	}

	ugs, err := models.UserConversations(
		models.UserConversationWhere.UserID.EQ(ctx.User.ID),
		qm.Load(qm.Rels(models.UserConversationRels.Conversation, models.ConversationRels.UserConversations)),
	).All(r.Context(), ctx.Components.Db)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not find recipients conversations: %+v", err), http.StatusInternalServerError)
		return
	}

	var conversation *models.Conversation
	sort.Ints(ids)
	for _, ug := range ugs {
		conversationIds := []int{}
		for _, ug := range ug.R.Conversation.R.UserConversations {
			// TODO: somehow traversing the dependencies brings redundant rows, e.g. the list we're iterating on can
			// 		yield [115, 115, 116] as user ids.
			redundant := false
			for _, id := range conversationIds {
				if ug.UserID == id {
					redundant = true
				}
			}
			if !redundant {
				conversationIds = append(conversationIds, ug.UserID)
			}
		}
		sort.Ints(conversationIds)
		if reflect.DeepEqual(conversationIds, ids) {
			conversation = ug.R.Conversation
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
	if conversation == nil {
		newConversation := models.Conversation{}
		if err = newConversation.Insert(r.Context(), tx, boil.Infer()); err != nil {
			tx.Rollback()
			http.Error(w, fmt.Sprintf("could not create new conversation: %+v", err), http.StatusInternalServerError)
			return
		}

		errs := make(chan error, 1)
		var wg sync.WaitGroup
		for _, id := range ids {
			wg.Add(1)
			uid := id
			go func() {
				ug := models.UserConversation{
					UserID:         uid,
					ConversationID: newConversation.ID,
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
			http.Error(w, fmt.Sprintf("could not add recipient to new conversation: %+v", err), http.StatusInternalServerError)
			return
		}
		conversation = &newConversation
	}

	message := models.Message{
		//Text:           msg.Text,
		AuthorID:       null.IntFrom(ctx.User.ID),
		ConversationID: conversation.ID,
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

	for _, rec := range recipients {
		if rec.ID == ctx.User.ID {
			continue
		}

		if err = ctx.Components.PgPubSub.Publish(entities.UserPubSubTopic(rec), common.NewMessageEvent{
			PubSubEvent:      common.PubSubEvent{Type: "newmessage", Timestamp: time.Now()},
			Text:             msg.Text,
			AuthorUuid:       ctx.User.UUID.String(),
			AuthorHandle:     ctx.User.Handle,
			ConversationUuid: msg.ConversationUuid,
		}); err != nil {
			log.Printf("failed to notify recipients of message: %+v", err)
		}
	}
}
