package entities

import (
	"fmt"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"strings"
)

func UserConversations(ctx common.Context, offset, pageSz int) (models.ConversationSlice, error) {
	myConversations, err := models.UserConversations(
		models.UserConversationWhere.UserID.EQ(ctx.User.ID),
		qm.OrderBy(fmt.Sprintf("%s DESC", models.UserConversationColumns.CreatedAt)),
		qm.Offset(offset), qm.Limit(pageSz),
	).All(ctx.Context, ctx.Components.Db)
	if err != nil {
		return nil, err
	}

	conversationIds := []int{}
	for _, ugs := range myConversations {
		conversationIds = append(conversationIds, ugs.ConversationID)
	}
	conversations, err := models.Conversations(
		models.ConversationWhere.ID.IN(conversationIds),
		qm.Load(qm.Rels(models.ConversationRels.UserConversations, models.UserConversationRels.User)),
	).All(ctx.Context, ctx.Components.Db)
	if err != nil {
		return nil, err
	}

	return conversations, nil
}

func UserPubSubTopic(u *models.User) string {
	return strings.Replace(fmt.Sprintf("user-conn-%s", u.UUID), "-", "_", -1)
}
