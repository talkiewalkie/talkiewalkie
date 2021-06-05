package unauthenticated

import (
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"

	"github.com/talkiewalkie/talkiewalkie/common"
)

type walkOutput struct {
	listedWalkOutput
	AudioUrl  string `json:"audioUrl"`
	LikeCount int64  `json:"likeCount"`
}

func WalkHandler(w http.ResponseWriter, r *http.Request, c *unauthenticatedContext) (interface{}, *common.HttpError) {
	uidRaw, ok := mux.Vars(r)["uuid"]
	if !ok {
		return nil, common.ServerError("bad route")
	}
	uid, err := uuid.FromString(uidRaw)
	if err != nil {
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

	return walkOutput{
		listedWalkOutput: listedWalkOutput{
			Uuid:     walk.UUID,
			Title:    walk.Title,
			Author:   authorOutput{Uuid: u.UUID, Handle: u.Handle},
			CoverUrl: coverUrl,
		},
		AudioUrl:  audioUrl,
		LikeCount: cnt,
	}, nil
}
