package authenticated

import (
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"

	"github.com/talkiewalkie/talkiewalkie/common"
)

type userByUuidWalkOutput struct {
	Uuid  string `json:"uuid"`
	Title string `json:"title"`
}

type userByUuidOutput struct {
	Uuid         string                 `json:"uuid"`
	Handle       string                 `json:"handle"`
	WalksCreated []userByUuidWalkOutput `json:"walks"`
}

func UserByUuidHandler(r *http.Request, c *authenticatedContext) (interface{}, *common.HttpError) {
	uid, ok := mux.Vars(r)["uuid"]
	if !ok {
		return nil, common.ServerError("bad request, could not get uuid")
	}
	if _, err := uuid.FromString(uid); err != nil {
		return nil, common.ServerError("bad request, could not parse uuid: %+v", err)
	}

	u, err := c.UserRepository.GetUserByUuid(uid)
	if err != nil {
		return nil, common.ServerError(err.Error())
	}

	walks, err := u.AuthorWalks().All(r.Context(), c.Db)
	wOutput := []userByUuidWalkOutput{}
	for _, w := range walks {
		wOutput = append(wOutput, userByUuidWalkOutput{Uuid: w.UUID, Title: w.Title})
	}

	return userByUuidOutput{
		Uuid:         u.UUID,
		Handle:       u.Handle,
		WalksCreated: wOutput,
	}, nil
}
