package server

import (
	"errors"
	"github.com/bms/dbModels"
	"net/http"
)

func (srv *Server) GetUserContext(req *http.Request) (*dbModels.UserContext, error) {
	user, ok := req.Context().Value(dbModels.UserContextKey).(*dbModels.UserContext)
	if !ok {
		return nil, errors.New("error getting userIDContext")
	}
	return user, nil
}

