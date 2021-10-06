package entities

import (
	"errors"
	"github.com/talkiewalkie/talkiewalkie/models"
	"strings"
)

func ConversationDisplay(c models.Conversation) string {
	display := c.Name.String
	if c.Name.Valid {
		return display
	}

	handles := []string{}
	for _, ug := range c.R.UserConversations {
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

	return strings.Join(handles, ", ")
}

func CanAccessConversation(c *models.Conversation, u *models.User) (bool, error) {
	if c.R.UserConversations == nil {
		return false, errors.New("Need user conversations to be eager loaded.")
	}

	for _, uc := range c.R.UserConversations {
		if uc.UserID == u.ID {
			return true, nil
		}
	}
	return false, nil
}
