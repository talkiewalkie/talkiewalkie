package unauthenticated

import (
	"net/http"

	"github.com/docker/distribution/uuid"
	"github.com/gorilla/mux"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
)

type walkResponse struct {
	*models.Walk
	Author   *models.User `db:"author" json:"author"`
	CoverUrl string       `json:"coverUrl"`
	AudioUrl string       `json:"audioUrl"`
}

func WalkHandler(w http.ResponseWriter, r *http.Request, c *unauthenticatedContext) (interface{}, *common.HttpError) {
	uid, ok := mux.Vars(r)["uuid"]
	if !ok {
		return nil, common.ServerError("bad route")
	}
	if _, err := uuid.Parse(uid); err != nil {
		return nil, common.ServerError("bad route")
	}

	walk, err := c.WalkRepository.GetByUuid(uid)
	if err != nil {
		return nil, common.ServerError(err.Error())
	}
	u, err := walk.Author().One(r.Context(), c.Db)
	if err != nil {
		return nil, common.ServerError("failed to load attached: %v", err)
	}
	cover, err := walk.Cover().One(r.Context(), c.Db)
	if err != nil {
		return nil, common.ServerError("failed to load cover", err)
	}
	coverUrl, err := c.Storage.Url(cover.UUID)
	if err != nil {
		return nil, common.ServerError(err.Error())
	}

	return walkResponse{Walk: walk, Author: u, CoverUrl: coverUrl, AudioUrl: ""}, nil
}
