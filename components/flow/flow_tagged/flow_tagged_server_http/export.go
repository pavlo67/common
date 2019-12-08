package flow_tagged_server_http

import (
	"net/http"

	"github.com/pkg/errors"

	"strconv"
	"strings"

	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/selectors"
	"github.com/pavlo67/workshop/common/server"
	"github.com/pavlo67/workshop/common/server/server_http"
)

// TODO!!! add parameters info into responces

// ListFlow --------------------------------------------------------------------------------------

const afterIDParam = "after_id"

var ExportFlowEndpoint = server_http.Endpoint{Method: "GET", PathParams: nil, QueryParams: []string{afterIDParam}, WorkerHTTP: ExportFlow}

func ExportFlow(user *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {
	queryParams := req.URL.Query()
	afterIDStr := strings.TrimSpace(queryParams.Get(afterIDParam))

	var afterID int
	if afterIDStr != "" {
		var err error
		afterID, err = strconv.Atoi(afterIDStr)
		if err != nil {
			return server.ResponseRESTError(http.StatusBadRequest, errors.Errorf("ERROR on GET workspace/...ExportFlow: ",
				errors.Errorf("can't strconv.Atoi(%s) for after_id parameter", afterIDStr, err)))
		}
	}

	selector := selectors.Binary(selectors.Gt, "id", selectors.Value{afterID})

	// l.Infof("111111111111 %#v", selector)

	items, err := flowTaggedOp.List(selector, &crud.GetOptions{OrderBy: []string{"id"}})

	l.Infof("%#v", items)

	if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on GET workspace/...ExportFlow: ", err))
	}

	return server.ResponseRESTOk(items)
}
