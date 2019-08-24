package ep_flow

import (
	"net/http"

	"github.com/pavlo67/workshop/apps/rest/flow/flow_starter"
	"github.com/pavlo67/workshop/basis/auth"
	"github.com/pavlo67/workshop/basis/server"
	"github.com/pavlo67/workshop/basis/server/server_http"
)

var ToInit = server_http.InitEndpoint(&flow_starter.Endpoints, "GET", "/flow/v1/list", nil, workerList, "")
var _ server_http.WorkerHTTP = workerList

func workerList(user *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {

	briefs, err := flow_starter.DataOp.List(nil, nil, nil)
	if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, err)
	}

	return server.ResponseRESTOk(briefs)
}
