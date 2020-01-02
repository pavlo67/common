package flowserver_http

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

// FlowList --------------------------------------------------------------------------------------

var listEndpoint = server_http.Endpoint{Method: "GET", PathParams: nil, WorkerHTTP: FlowList}

func FlowList(user *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {
	items, err := dataTaggedOp.List(nil, &crud.GetOptions{Limit1: 200})

	l.Debugf("%#v", items)

	if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on GET ...FlowList: ", err))
	}

	return server.ResponseRESTOk(items)
}

var readEndpoint = server_http.Endpoint{Method: "GET", PathParams: []string{"id"}, WorkerHTTP: FlowRead}

func FlowRead(user *auth.User, params server_http.Params, req *http.Request) (server.Response, error) {
	id := common.ID(params["id"])

	item, err := dataTaggedOp.Read(id, nil)
	if err == common.ErrNotFound {
		return server.ResponseRESTError(http.StatusNotFound, errors.Errorf("ERROR on GET ...FlowRead: not found item with id = %s", id))
	} else if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on GET ...FlowRead: ", err))
	}

	err = dataTaggedOp.SetDetails(item)
	if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on GET ...FlowRead: ", err))
	}

	return server.ResponseRESTOk(item)
}
