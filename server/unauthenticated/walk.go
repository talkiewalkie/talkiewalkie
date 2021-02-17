package unauthenticated

import (
	"net/http"

	"github.com/docker/distribution/uuid"
	"github.com/gorilla/mux"

	"github.com/talkiewalkie/talkiewalkie/common"
)

type walkOutput struct {
	listedWalkOutput
	AudioUrl string `json:"audioUrl"`
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
	audio, err := walk.Audio().One(r.Context(), c.Db)
	if err != nil {
		return nil, common.ServerError(err.Error())
	}
	audioUrl, err := c.Storage.Url(audio.UUID)
	if err != nil {
		return nil, common.ServerError(err.Error())
	}

	return walkOutput{
		listedWalkOutput: listedWalkOutput{
			Uuid:     walk.UUID,
			Title:    walk.Title,
			Author:   authorOutput{Uuid: u.UUID, Handle: u.Handle},
			CoverUrl: coverUrl,
		},
		AudioUrl: audioUrl,
	}, nil
}
