package server

import (
	"github.com/gauravcoco/bms/dbModels"
	"net/http"
)

func (srv *Server) GetUserContext(req *http.Request) *dbModels.UserContext {
	return req.Context().Value(dbModels.UserContextKey).(*dbModels.UserContext)

}
