package authenticated

import (
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
)

type likeWalkOutput struct {
}

func LikeWalkByUuid(r *http.Request, c *authenticatedContext) (interface{}, *common.HttpError) {
	uidRaw, ok := mux.Vars(r)["uuid"]
	uid, err := uuid.FromString(uidRaw)
	if err != nil {
		return nil, common.ServerError("failed to parse uuid: '%s': %+v", uidRaw, err)
	}

	if !ok {
		return nil, common.ServerError("could not parse request")
	}

	w, err := models.Walks(models.WalkWhere.UUID.EQ(uid)).One(r.Context(), c.Db)
	if err != nil {
		return nil, common.ServerError("could not find walk: %+v", err)
	}

	err = (&models.UserWalk{UserID: c.User.ID, WalkID: w.ID}).Insert(r.Context(), c.Db, boil.Infer())
	if err != nil {
		return nil, common.ServerError(err.Error())
	}

	return likeWalkOutput{}, nil
}
