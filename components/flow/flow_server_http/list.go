package flow_server_http

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/server"
	"github.com/pavlo67/workshop/common/server/server_http"
)

// TODO!!! add parameters info into responces

// Flow --------------------------------------------------------------------------------------

var FlowEndpoint = server_http.Endpoint{Method: "GET", PathParams: nil, WorkerHTTP: Flow}

func Flow(user *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {
	items, err := flowTaggedOp.List(nil, &crud.GetOptions{Limit1: 200})

	l.Debugf("%#v", items)

	if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on GET ...Flow: ", err))
	}

	return server.ResponseRESTOk(items)
}

var FlowReadEndpoint = server_http.Endpoint{Method: "GET", PathParams: []string{"id"}, WorkerHTTP: FlowRead}

func FlowRead(user *auth.User, params server_http.Params, req *http.Request) (server.Response, error) {
	id := common.ID(params["id"])

	item, err := flowTaggedOp.Read(id, nil)
	if err == common.ErrNotFound {
		return server.ResponseRESTError(http.StatusNotFound, errors.Errorf("ERROR on GET ...FlowRead: not found item with id = %s", id))
	} else if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on GET ...FlowRead: ", err))
	}

	err = flowTaggedOp.SetDetails(item)
	if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on GET ...FlowRead: ", err))
	}

	return server.ResponseRESTOk(item)
}
