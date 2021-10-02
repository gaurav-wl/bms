package server

import (
	"errors"
	bmsError "github.com/gauravcoco/bms/errors"
	"github.com/gauravcoco/bms/utils"
	"net/http"
)

func (srv *Server) userInfo(resp http.ResponseWriter, req *http.Request) {
	uc := srv.GetUserContext(req)
	if uc == nil {
		bmsError.RespondClientErr(resp, req, errors.New("error getting user context"), http.StatusUnauthorized, "Error getting user")
		return
	}

	user, err := srv.DBHelper.GetUserByID(req.Context(), uc.ID)
	if err != nil {
		bmsError.RespondGenericServerErr(resp, req, err, "Error getting user")
		return
	}

	utils.EncodeJSONBody(resp, http.StatusOK, srv.Converter.ToUser(user))
}
