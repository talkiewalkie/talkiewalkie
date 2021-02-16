package unauthenticated

import (
	"net/http"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/repository"
)

type walkAuthor struct {
	*repository.Walk
	Author *repository.User  `db:"author" json:"author"`
	Cover  *repository.Asset `db:"cover" json:"cover"`
	Audio  *repository.Asset `db:"audio" json:"audio"`
}

type walkResponse struct {
	*walkAuthor
	CoverUrl string `json:"coverUrl"`
	AudioUrl string `json:"audioUrl"`
}

func WalksHandler(c *common.Components) UnauthHandler {
	return func(w http.ResponseWriter, r *http.Request, ctx *unauthenticatedContext) (interface{}, *common.HttpError) {
		walks := []walkAuthor{}
		ctx.Db.Unsafe()

		if err := ctx.Db.Select(&walks, `
		SELECT walk.*, 
			   author.id AS "author.id", 
			   author.uuid AS "author.uuid", 
			   author.handle AS "author.handle", 
			   author.email AS "author.email", 
			   author.password AS "author.password", 
			   author.email_token AS "author.email_token",
		       cover.id AS "cover.id",
		       cover.uuid AS "cover.uuid",
		       cover.file_name AS "cover.file_name",
		       cover.mime_type AS "cover.mime_type",
		       cover.uploaded_at AS "cover.uploaded_at",
		       audio.id AS "audio.id",
		       audio.uuid AS "audio.uuid",
		       audio.file_name AS "audio.file_name",
		       audio.mime_type AS "audio.mime_type",
		       audio.uploaded_at AS "audio.uploaded_at"
		FROM "walk" 
		    INNER JOIN "user" AS author ON author.id = author_id
		    INNER JOIN "asset" AS cover ON cover.id = walk.cover_id
		    INNER JOIN "asset" AS audio ON audio.id = walk.audio_id;`,
		); err != nil {
			return nil, common.ServerError("could not fetch walks: %v", err)
		}

		responseWalks := []walkResponse{}
		for _, walk := range walks {
			coverUrl, err := c.Storage.Url(walk.Cover.Uuid)
			if err != nil {
				return nil, common.ServerError("could not make a signed url: %v", err)
			}
			audioUrl, err := c.Storage.Url(walk.Audio.Uuid)
			if err != nil {
				return nil, common.ServerError("could not make a signed url: %v", err)
			}

			responseWalks = append(responseWalks, walkResponse{
				walkAuthor: &walk,
				CoverUrl:   coverUrl,
				AudioUrl:   audioUrl,
			})
		}

		return responseWalks, nil
	}
}
