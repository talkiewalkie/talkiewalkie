package unauthenticated

import (
	"fmt"
	"net/http"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
)

type authorOutput struct {
	Uuid   string `json:"uuid"`
	Handle string `json:"handle"`
}

type listedWalkOutput struct {
	Uuid     string       `json:"uuid"`
	Title    string       `json:"title"`
	Author   authorOutput `json:"author"`
	CoverUrl string       `json:"coverUrl"`
}

func WalksHandler(w http.ResponseWriter, r *http.Request, c *unauthenticatedContext) (interface{}, *common.HttpError) {
	walks, err := models.Walks(
		qm.Load(models.WalkRels.Author),
		qm.Load(models.WalkRels.Cover),
		qm.Load(models.WalkRels.Audio),
		qm.OrderBy(fmt.Sprintf("%s DESC", models.WalkColumns.CreatedAt)),
		qm.Limit(20)).
		All(r.Context(), c.Db)
	if err != nil {
		return nil, common.ServerError("failed to fetch walks: %v", err)
	}

	responseWalks := []listedWalkOutput{}
	for _, walk := range walks {
		coverUrl, err := c.Storage.Url(walk.R.Cover.UUID)
		if err != nil {
			return nil, common.ServerError("could not make a signed url: %v", err)
		}

		responseWalks = append(responseWalks, listedWalkOutput{
			Uuid:     walk.UUID,
			Title:    walk.Title,
			Author:   authorOutput{walk.R.Author.UUID, walk.R.Author.Handle},
			CoverUrl: coverUrl,
		})
	}
	return responseWalks, nil
}
