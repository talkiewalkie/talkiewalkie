package routes

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
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
	Uuid              uuid.UUID             `json:"uuid"`
	Title             string                `json:"title"`
	Description       string                `json:"description"`
	Author            AuthorWalksItemOutput `json:"author"`
	CoverUrl          string                `json:"coverUrl"`
	DistanceFromPoint float64               `json:"distanceFromPoint"`
}

func Walks(w http.ResponseWriter, r *http.Request) {
	ctx := common.WithContext(r)

	params := r.URL.Query()
	lngs, foundLng := params["lng"]
	lats, foundLat := params["lat"]
	if foundLng != foundLat {
		http.Error(w, "provide none or both coords params (lng,lat)", http.StatusBadRequest)
		return
	}
	var lng, lat float64
	if foundLat {
		val, err := strconv.ParseFloat(lngs[0], 64)
		if err != nil {
			http.Error(w, "bad longitude", http.StatusBadRequest)
			return
		}
		lng = val
		val, err = strconv.ParseFloat(lats[0], 64)
		if err != nil {
			http.Error(w, "bad latitude", http.StatusBadRequest)
			return
		}
		lat = val
	}

	var offset int
	if vals, ok := params["offset"]; ok && len(vals) > 0 {
		v, err := strconv.Atoi(vals[0])
		if err != nil {
			http.Error(w, "bad offset", http.StatusBadRequest)
			return
		}
		offset = v
	}

	walks, err := models.Walks(
		qm.Limit(20),
		qm.Offset(offset),
		qm.OrderBy(fmt.Sprintf(
			"earth_distance(ll_to_earth(%s[0], %s[1]),  ll_to_earth(%f, %f))",
			models.WalkColumns.StartPoint,
			models.WalkColumns.StartPoint,
			lat,
			lng,
		)),
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
			// XXX: it's a shame to recompute what postgres has done for the sorting, but it's simpler atm
			DistanceFromPoint: Distance(lat, lng, walk.StartPoint.X, walk.StartPoint.Y),
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

// HELPERS

// https://gist.github.com/cdipaolo/d3f8db3848278b49db68
// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// Distance function returns the distance (in meters) between two points of
//     a given longitude and latitude relatively accurately (using a spherical
//     approximation of the Earth) through the Haversin Distance Formula for
//     great arc distance on a sphere with accuracy for small distances
//
// point coordinates are supplied in degrees and converted into rad. in the func
//
// distance returned is METERS!!!!!!
// http://en.wikipedia.org/wiki/Haversine_formula
func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}
