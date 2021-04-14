package authenticated

import (
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"

	"github.com/talkiewalkie/talkiewalkie/common"
)

type userByUuidWalkOutput struct {
	Uuid  uuid.UUID `json:"uuid"`
	Title string    `json:"title"`
}

type userByUuidOutput struct {
	Uuid         uuid.UUID              `json:"uuid"`
	Handle       string                 `json:"handle"`
	WalksCreated []userByUuidWalkOutput `json:"walks"`
}

func UserByUuidHandler(r *http.Request, c *authenticatedContext) (interface{}, *common.HttpError) {
	uidRaw, ok := mux.Vars(r)["uuid"]
	if !ok {
		return nil, common.ServerError("bad request, could not get uuid")
	}
	uid, err := uuid.FromString(uidRaw)
	if err != nil {
		return nil, common.ServerError("bad uuid '%s': %+v", uidRaw, err)
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
