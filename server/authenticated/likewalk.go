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
	uid, ok := mux.Vars(r)["uuid"]
	if !ok {
		return nil, common.ServerError("could not parse request")
	}

	if _, err := uuid.FromString(uid); err != nil {
		return nil, common.ServerError("could not parse request: %+v", err)
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
