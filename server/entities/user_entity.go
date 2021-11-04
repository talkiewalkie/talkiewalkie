package entities

import (
	"fmt"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/talkiewalkie/talkiewalkie/pb"
	"strings"
)

type User struct {
	Record     *models.User
	Components *common.Components
}

func (user User) DisplayName() string {
	if user.Record.DisplayName.Valid {
		return user.Record.DisplayName.String
	} else {
		return user.Record.PhoneNumber
	}
}

func (user User) PubSubTopic() string {
	return strings.Replace(fmt.Sprintf("user-conn-%s", user.Record.UUID), "-", "_", -1)
}

func (user User) ToPb() *pb.User {
	return &pb.User{
		DisplayName:   user.DisplayName(),
		Uuid:          user.Record.UUID.String(),
		Conversations: nil,
		Phone:         user.Record.PhoneNumber,
	}
}

func (user User) HasAccess(u *User) (bool, error) {
	ucs, err := user.Components.UserConversationStore.ByUserIds(user.Record.ID, u.Record.ID)
	if err != nil {
		return false, err
	}

	for _, uc := range ucs[0] {
		found := false
		for _, otherUc := range ucs[1] {
			if uc.Record.ConversationID == otherUc.Record.ConversationID {
				found = true
				break
			}
		}

		if found {
			return true, nil
		}
	}

	return false, nil
}
