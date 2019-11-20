package rights_server

import (
	"net/http"

	"github.com/pavlo67/workshop/common/server"
	"github.com/pavlo67/workshop/common/server/server_http"
	"github.com/pavlo67/workshop/components/auth"
)

var _ = server_http.InitEndpoint(&endpoints, "GET", "/", nil, addRights, "")

var _ server_http.WorkerHTTP = addRights

func addRights(*auth.User, libs.Params, *http.Request) (server.Response, error) {
	return server.ResponseREST(http.StatusOK, nil)
}
