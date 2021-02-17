package unauthenticated

import (
	"net/http"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
)

type walkItemResponse struct {
	*models.Walk
	Author   *models.User `db:"author" json:"author"`
	CoverUrl string       `json:"coverUrl"`
}

func WalksHandler(w http.ResponseWriter, r *http.Request, c *unauthenticatedContext) (interface{}, *common.HttpError) {
	walks, err := models.Walks(
		qm.Load(models.WalkRels.Author),
		qm.Load(models.WalkRels.Cover),
		qm.Load(models.WalkRels.Audio),
		qm.Limit(20)).
		All(r.Context(), c.Db)
	if err != nil {
		return nil, common.ServerError("failed to fetch walks: %v", err)
	}

	responseWalks := []walkItemResponse{}
	for _, walk := range walks {
		coverUrl, err := c.Storage.Url(walk.R.Cover.UUID)
		if err != nil {
			return nil, common.ServerError("could not make a signed url: %v", err)
		}

		responseWalks = append(responseWalks, walkItemResponse{
			Walk:     walk,
			Author:   walk.R.Author,
			CoverUrl: coverUrl,
		})
	}
	return responseWalks, nil
}
