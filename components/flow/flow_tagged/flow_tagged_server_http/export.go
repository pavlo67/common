package flow_tagged_server_http

import (
	"net/http"

	"github.com/pkg/errors"

	"strconv"
	"strings"

	"time"

	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/selectors"
	"github.com/pavlo67/workshop/common/server"
	"github.com/pavlo67/workshop/common/server/server_http"
	"github.com/pavlo67/workshop/common/transport"
)

// TODO!!! add parameters info into responces

// ListFlow --------------------------------------------------------------------------------------

const AfterIDParam = "after_id"

var ExportFlowEndpoint = server_http.Endpoint{Method: "GET", PathParams: nil, QueryParams: []string{AfterIDParam}, WorkerHTTP: ExportFlow}

func ExportFlow(user *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {
	queryParams := req.URL.Query()
	afterIDStr := strings.TrimSpace(queryParams.Get(AfterIDParam))

	var afterID int
	if afterIDStr != "" {
		var err error
		afterID, err = strconv.Atoi(afterIDStr)
		if err != nil {
			return server.ResponseRESTError(http.StatusBadRequest, errors.Errorf("ERROR on GET workspace/...ExportFlow: ",
				errors.Errorf("can't strconv.Atoi(%s) for after_id parameter", afterIDStr, err)))
		}
	}

	// TODO!!! selector with item.CreatedAt / UpdatedAt if original .ID isn't autoincrement (for mongoDB, for example)

	selector := selectors.Binary(selectors.Gt, "id", selectors.Value{afterID})
	items, err := flowTaggedOp.List(selector, &crud.GetOptions{OrderBy: []string{"id"}})

	l.Infof("%#v / %#v", items, selector)

	if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on GET workspace/...ExportFlow: ", err))
	}

	// TODO!!! MaxID with item.CreatedAt / UpdatedAt if original .ID isn't autoincrement (for mongoDB, for example)

	var maxID uint64
	if len(items) > 0 {
		maxIDParsed, err := strconv.ParseUint(string(items[len(items)-1].ID), 10, 64)
		if err != nil {
			l.Errorf("can't strconv.ParseUint(%s) with .ID: %s", string(items[len(items)-1].ID), err)
		} else {
			maxID = maxIDParsed
		}
	}

	return server.ResponseRESTOk(transport.Packet{
		// SourceURL: "",
		CreatedAt: time.Now(),
		Type:      transport.DataItemsDataType,
		Data:      items,
		MaxID:     maxID,
	})
}
