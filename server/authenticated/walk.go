package authenticated

import (
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
)

type authorOutput struct {
	Uuid   uuid.UUID `json:"uuid"`
	Handle string    `json:"handle"`
}

type walkOutput struct {
	Uuid      uuid.UUID    `json:"uuid"`
	Title     string       `json:"title"`
	Author    authorOutput `json:"author"`
	CoverUrl  string       `json:"coverUrl"`
	AudioUrl  string       `json:"audioUrl"`
	LikeCount int64        `json:"likeCount"`
	IsLiked   bool         `json:"isLiked"`
}

func WalkHandler(r *http.Request, c *authenticatedContext) (interface{}, *common.HttpError) {
	uidRaw, ok := mux.Vars(r)["uuid"]
	if !ok {
		return nil, common.ServerError("bad route")
	}
	uid, err := uuid.FromString(uidRaw)
	if err != nil {
		return nil, common.ServerError("bad route")
	}

	walk, err := models.Walks(models.WalkWhere.UUID.EQ(uid)).One(r.Context(), c.Db)
	if err != nil {
		return nil, common.ServerError(err.Error())
	}
	u, err := walk.Author().One(r.Context(), c.Db)
	if err != nil {
		return nil, common.ServerError("failed to load attached: %v", err)
	}

	exists, err := walk.Cover().Exists(r.Context(), c.Db)
	if err != nil {
		return nil, common.ServerError("failed to check existence of cover", err)
	}

	var coverUrl string
	if exists {

		cover, err := walk.Cover().One(r.Context(), c.Db)
		if err != nil {
			return nil, common.ServerError("failed to load cover", err)
		}

		coverUrl, err = c.Storage.Url(cover.UUID.String())
		if err != nil {
			return nil, common.ServerError(err.Error())
		}
	}

	audio, err := walk.Audio().One(r.Context(), c.Db)
	if err != nil {
		return nil, common.ServerError(err.Error())
	}
	audioUrl, err := c.Storage.Url(audio.UUID.String())
	if err != nil {
		return nil, common.ServerError(err.Error())
	}

	cnt, err := walk.UserWalks().Count(r.Context(), c.Db)
	if err != nil {
		return nil, common.ServerError("could not count likes: %+v", err)
	}

	uw, err := models.UserWalks(models.UserWalkWhere.UserID.EQ(c.User.ID), models.UserWalkWhere.WalkID.EQ(walk.ID)).Exists(r.Context(), c.Db)
	if err != nil {
		return nil, common.ServerError("could not retrieve user walk: %+v", err)
	}

	return walkOutput{
		Uuid:      walk.UUID,
		Title:     walk.Title,
		Author:    authorOutput{Uuid: u.UUID, Handle: u.Handle},
		CoverUrl:  coverUrl,
		AudioUrl:  audioUrl,
		LikeCount: cnt,
		IsLiked:   uw,
	}, nil
}
