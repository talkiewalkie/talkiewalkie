package authenticated

import (
	"net/http"

	"github.com/talkiewalkie/talkiewalkie/common"
)

func UserListHandler(r *http.Request, c *authenticatedContext) (interface{}, *common.HttpError) {
	return nil, &common.HttpError{Code: http.StatusInternalServerError, Msg: "not implemented"}

}
