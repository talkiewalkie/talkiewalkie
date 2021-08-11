package routes

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// ------------

type CreateWalkInput struct {
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	CoverArtUuid uuid.UUID `json:"coverArtUuid"`
	AudioUuid    uuid.UUID `json:"audioUuid"`
}

type CreateWalkOutput struct {
	Uuid  uuid.UUID `json:"uuid"`
	Title string    `json:"title"`
}

func CreateWalk(w http.ResponseWriter, r *http.Request) {
	ctx := common.WithAuthedContext(r)

	var p CreateWalkInput
	if err := common.JsonIn(r, &p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	assets, err := models.Assets(
		models.AssetWhere.UUID.IN([]uuid.UUID{p.AudioUuid, p.CoverArtUuid}),
	).All(r.Context(), ctx.Components.Db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var audio, cover *models.Asset
	if strings.HasPrefix(assets[0].MimeType, "image") && strings.HasPrefix(assets[1].MimeType, "video") {
		cover = assets[0]
		audio = assets[1]
	} else if strings.HasPrefix(assets[1].MimeType, "image") && strings.HasPrefix(assets[0].MimeType, "video") {
		cover = assets[1]
		audio = assets[0]
	} else {
		http.Error(w, "bad asset references", http.StatusInternalServerError)
		return
	}

	walk := &models.Walk{
		Title:    p.Title,
		CoverID:  null.NewInt(cover.ID, true),
		AudioID:  null.NewInt(audio.ID, true),
		AuthorID: ctx.User.ID,
	}
	if err := walk.Insert(r.Context(), ctx.Components.Db, boil.Infer()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.JsonOut(w, CreateWalkOutput{walk.UUID, walk.Title})
}

// ------------

type AuthorWalksItemOutput struct {
	Uuid   uuid.UUID `json:"uuid"`
	Handle string    `json:"handle"`
}

type WalksItemOutput struct {
	Uuid        uuid.UUID             `json:"uuid"`
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Author      AuthorWalksItemOutput `json:"author"`
	CoverUrl    string                `json:"coverUrl"`
}

func Walks(w http.ResponseWriter, r *http.Request) {
	ctx := common.WithContext(r)

	walks, err := models.Walks(
		qm.Limit(20),
		qm.Offset(0),
		qm.OrderBy(models.WalkColumns.CreatedAt),
		qm.Load(models.WalkRels.Author),
		qm.Load(models.WalkRels.Cover)).All(r.Context(), ctx.Components.Db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseWalks := []WalksItemOutput{}
	for _, walk := range walks {
		var coverUrl string
		if walk.R.Cover != nil {
			coverUrl, err = ctx.Components.Storage.AssetUrl(walk.R.Cover)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		responseWalks = append(responseWalks, WalksItemOutput{
			Uuid:        walk.UUID,
			Title:       walk.Title,
			Description: walk.Description.String,
			Author:      AuthorWalksItemOutput{walk.R.Author.UUID, walk.R.Author.Handle},
			CoverUrl:    coverUrl,
		})
	}
	common.JsonOut(w, responseWalks)
}

// ------------

type WalkByUuidOutput struct {
	Uuid        uuid.UUID             `json:"uuid"`
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Author      AuthorWalksItemOutput `json:"author"`
	CoverUrl    string                `json:"coverUrl"`
	AudioUrl    string                `json:"audioUrl"`
}

func WalkByUuid(w http.ResponseWriter, r *http.Request) {
	ctx := common.WithContext(r)

	vars := mux.Vars(r)
	uuid_url, ok := vars["uuid"]
	if !ok {
		http.Error(w, "expects a well formed uuid", http.StatusBadRequest)
		return
	}

	uid, err := uuid.FromString(uuid_url)
	if err != nil {
		http.Error(w, "expects a well formed uuid", http.StatusBadRequest)
		return
	}

	walk, err := models.Walks(
		qm.Where("uuid = ?", uid),
		qm.Load(models.WalkRels.Author),
		qm.Load(models.WalkRels.Cover),
		qm.Load(models.WalkRels.Audio)).One(r.Context(), ctx.Components.Db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if walk == nil {
		http.Error(w, "no walk found", http.StatusBadRequest)
		return
	}

	var coverUrl string
	if walk.R.Cover != nil {
		coverUrl, err = ctx.Components.Storage.AssetUrl(walk.R.Cover)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	audioUrl, err := ctx.Components.Storage.AssetUrl(walk.R.Audio)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out := WalkByUuidOutput{
		Uuid:        walk.UUID,
		Title:       walk.Title,
		Description: walk.Description.String,
		Author:      AuthorWalksItemOutput{Uuid: walk.R.Author.UUID, Handle: walk.R.Author.Handle},
		CoverUrl:    coverUrl,
		AudioUrl:    audioUrl,
	}
	common.JsonOut(w, out)
}
