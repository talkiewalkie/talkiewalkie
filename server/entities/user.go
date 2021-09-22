package entities

import (
	"fmt"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"strings"
)

func UserConversations(ctx common.Context, offset, pageSz int) (models.GroupSlice, error) {
	myGroups, err := models.UserGroups(
		models.UserGroupWhere.UserID.EQ(ctx.User.ID),
		qm.OrderBy(fmt.Sprintf("%s DESC", models.UserGroupColumns.CreatedAt)),
		qm.Offset(offset), qm.Limit(pageSz),
	).All(ctx.Context, ctx.Components.Db)
	if err != nil {
		return nil, err
	}

	groupIds := []int{}
	for _, ugs := range myGroups {
		groupIds = append(groupIds, ugs.GroupID)
	}
	groups, err := models.Groups(
		models.GroupWhere.ID.IN(groupIds),
		qm.Load(qm.Rels(models.GroupRels.UserGroups, models.UserGroupRels.User)),
	).All(ctx.Context, ctx.Components.Db)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func UserPubSubTopic(u *models.User) string {
	return strings.Replace(fmt.Sprintf("user-conn-%s", u.UUID), "-", "_", -1)
}
