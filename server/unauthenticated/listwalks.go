package unauthenticated

import (
	"net/http"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
)

type walkResponse struct {
	*models.Walk
	Author   *models.User `db:"author" json:"author"`
	CoverUrl string       `json:"coverUrl"`
	AudioUrl string       `json:"audioUrl"`
}

func WalksHandler(c *common.Components) UnauthHandler {
	return func(w http.ResponseWriter, r *http.Request, ctx *unauthenticatedContext) (interface{}, *common.HttpError) {
		walks, err := models.Walks(
			qm.Load(models.WalkRels.Author),
			qm.Load(models.WalkRels.Cover),
			qm.Load(models.WalkRels.Audio),
			qm.Limit(20)).
			All(r.Context(), ctx.Db)
		if err != nil {
			return nil, common.ServerError("failed to fetch walks: %v", err)
		}

		responseWalks := []walkResponse{}
		for _, walk := range walks {
			coverUrl, err := c.Storage.Url(walk.R.Cover.UUID)
			if err != nil {
				return nil, common.ServerError("could not make a signed url: %v", err)
			}
			audioUrl, err := c.Storage.Url(walk.R.Audio.UUID)
			if err != nil {
				return nil, common.ServerError("could not make a signed url: %v", err)
			}

			responseWalks = append(responseWalks, walkResponse{
				Walk:     walk,
				Author:   walk.R.Author,
				CoverUrl: coverUrl,
				AudioUrl: audioUrl,
			})
		}
		return responseWalks, nil
	}
}
