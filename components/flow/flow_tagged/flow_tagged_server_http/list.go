package flow_tagged_server_http

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/server"
	"github.com/pavlo67/workshop/common/server/server_http"
)

// TODO!!! add parameters info into responces

// ListFlow --------------------------------------------------------------------------------------

var ListFlowEndpoint = server_http.Endpoint{Method: "GET", PathParams: nil, WorkerHTTP: ListFlow}

func ListFlow(user *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {
	items, err := flowTaggedOp.List(nil, &crud.GetOptions{Limit1: 200})

	l.Infof("%#v", items)

	if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on GET workspace/...ListFlow: ", err))
	}

	return server.ResponseRESTOk(items)
}
