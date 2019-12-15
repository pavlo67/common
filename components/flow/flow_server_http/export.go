package flow_server_http

import (
	"net/http"

	"github.com/pkg/errors"

	"time"

	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/server"
	"github.com/pavlo67/workshop/common/server/server_http"
	"github.com/pavlo67/workshop/common/transport"
)

// TODO!!! add parameters info into responces

// Flow --------------------------------------------------------------------------------------

const AfterIDParam = "after_id"

var ExportFlowEndpoint = server_http.Endpoint{Method: "GET", PathParams: nil, QueryParams: []string{AfterIDParam}, WorkerHTTP: ExportFlow}

func ExportFlow(user *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {
	afterIDStr := req.URL.Query().Get(AfterIDParam)

	items, err := flowTaggedOp.Export(afterIDStr, nil)
	if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on GET storage/...ExportFlow: ", err))
	}

	return server.ResponseRESTOk(transport.Packet{
		// SourceURL: "",
		CreatedAt: time.Now(),
		Type:      transport.DataItemsDataType,
		Data:      items,
	})
}
