package rector_server

import (
	"net/http"

	"github.com/pavlo67/workshop/basis/auth"
	"github.com/pavlo67/workshop/basis/common"
	"github.com/pavlo67/workshop/basis/server"
	"github.com/pavlo67/workshop/basis/server/server_http"
)

var _ = server_http.InitEndpoint(&endpoints, "GET", "/", nil, addRights, "")

var _ server_http.WorkerHTTP = addRights

func addRights(*auth.User, common.Params, *http.Request) (server.Response, error) {
	return server.ResponseREST(http.StatusOK, nil)
}
