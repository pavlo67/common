package receiver_server_http

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/server"
	"github.com/pavlo67/workshop/common/server/server_http"

	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/packs"
)

// TODO!!! add parameters info into responces

// Flow --------------------------------------------------------------------------------------

var ExportFlowEndpoint = server_http.Endpoint{Method: "POST", WorkerHTTP: ExportFlow}

func ExportFlow(user *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {

	items, err := receiverTaggedOp.Export(afterIDStr, nil)
	if err != nil {
		return server.ResponseRESTError(http.StatusInternalServerError, errors.Errorf("ERROR on GET storage/...ExportFlow: ", err))
	}

	return server.ResponseRESTOk(packs.Pack{
		// SourceURL: "",
		// SentAt:  time.Now(),
		TypeKey: data.TypesKeyDataItems,
		Content: items,
	})
}
