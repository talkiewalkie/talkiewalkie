package authenticated

import (
	"net/http"
	"strings"

	"github.com/satori/go.uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
)

type createWalkInput struct {
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	CoverArtUuid uuid.UUID `json:"coverArtUuid"`
	AudioUuid    uuid.UUID `json:"audioUuid"`
}

func CreateWalkHandler(r *http.Request, ctx *authenticatedContext) (interface{}, *common.HttpError) {
	var p createWalkInput
	if err := common.JsonIn(r, &p); err != nil {
		return nil, common.ServerError(err.Error())
	}

	assets, err := ctx.AssetRepository.GetAllByUuid([]string{p.AudioUuid.String(), p.CoverArtUuid.String()})
	if err != nil {
		return nil, common.ServerError("could not find assets in db: %v", err)
	}

	if assets == nil || len(assets) == 0 {
		return nil, common.ServerError("did not find assets")
	}

	var audio, cover *models.Asset
	if strings.HasPrefix(assets[0].MimeType, "image") && strings.HasPrefix(assets[1].MimeType, "video") {
		cover = assets[0]
		audio = assets[1]
	} else if strings.HasPrefix(assets[1].MimeType, "image") && strings.HasPrefix(assets[0].MimeType, "video") {
		cover = assets[1]
		audio = assets[0]
	} else {
		return nil, common.ServerError("bad asset references")
	}

	w := &models.Walk{Title: p.Title, CoverID: null.NewInt(cover.ID, true), AudioID: null.NewInt(audio.ID, true), AuthorID: ctx.User.ID}
	if err := w.Insert(r.Context(), ctx.Db, boil.Infer()); err != nil {
		return nil, common.ServerError(err.Error())
	}
	return w, nil
}
