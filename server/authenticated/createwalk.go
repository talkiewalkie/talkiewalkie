package authenticated

import (
	"net/http"
	"strings"

	"github.com/satori/go.uuid"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/repository"
)

type createWalkPayload struct {
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	CoverArtUuid uuid.UUID `json:"coverArtUuid"`
	AudioUuid    uuid.UUID `json:"audioUuid"`
}

func CreateWalkHandler(r *http.Request, ctx *authenticatedContext) (interface{}, *common.HttpError) {
	var p createWalkPayload
	if err := common.JsonIn(r, &p); err != nil {
		return nil, common.ServerError(err.Error())
	}

	assets, err := ctx.AssetRepository.GetAllByUuid([]uuid.UUID{p.AudioUuid, p.CoverArtUuid})
	if err != nil {
		return nil, common.ServerError("could not find assets in db: %v", err)
	}

	var audio, cover repository.Asset
	if strings.HasPrefix(assets[0].MimeType, "image") && strings.HasPrefix(assets[1].MimeType, "video") {
		cover = assets[0]
		audio = assets[1]
	} else if strings.HasPrefix(assets[1].MimeType, "image") && strings.HasPrefix(assets[0].MimeType, "video") {
		cover = assets[1]
		audio = assets[0]
	} else {
		return nil, common.ServerError("bad asset references: %v", err)
	}

	ctx.Db.MustExec(`INSERT INTO "walk" (title, cover_id, audio_id,  author_id) VALUES ($1, $2, $3, $4)`, p.Title, cover.Id, audio.Id, ctx.User.Id)
	return nil, nil
}
