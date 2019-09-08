package ep_flow

import (
	"net/http"

	"github.com/pavlo67/workshop/apps/rest/flow/flow_routes"
	"github.com/pavlo67/workshop/common/libs/filelib"
	"github.com/pavlo67/workshop/common/server"
	"github.com/pavlo67/workshop/common/server/server_http"
	"github.com/pavlo67/workshop/components/auth"
)

var ToInit = server_http.InitEndpoint(&flow_routes.Endpoints, "GET", filelib.RelativePath(flow_routes.Prefix, filelib.CurrentFile(true)), nil, workerList, "")

var _ server_http.WorkerHTTP = workerList

func workerList(user *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {

	briefs, err := flow_routes.DataOp.List(nil, nil, nil)
	if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, err)
	}

	return server.ResponseRESTOk(briefs)
}
