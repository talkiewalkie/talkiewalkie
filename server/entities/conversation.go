package entities

import (
	"errors"
	"github.com/talkiewalkie/talkiewalkie/models"
	"strings"
)

func ConversationDisplay(c *models.Conversation) (string, error) {
	if c.R.UserConversations == nil {
		return "", errors.New("need to load user conversations")
	} else if len(c.R.UserConversations) > 0 && c.R.UserConversations[0].R.User == nil {
		return "", errors.New("need to load user conversations users")
	}

	display := c.Name.String
	if c.Name.Valid {
		return display, nil
	}

	handles := []string{}
	for _, ug := range c.R.UserConversations {
		redundant := false
		displayName := UserDisplayName(ug.R.User)
		for _, h := range handles {
			if h == displayName {
				redundant = true
			}
		}
		if !redundant {
			handles = append(handles, displayName)
		}
	}

	return strings.Join(handles, ", "), nil
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
