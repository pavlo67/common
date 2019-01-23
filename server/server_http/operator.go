package server_http

import (
	"net/http"

	"github.com/pavlo67/punctum/auth"
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/server"
	"github.com/pavlo67/punctum/server/router"
	"github.com/pavlo67/punctum/starter/joiner"
	"github.com/pkg/errors"
)

const InterfaceKey joiner.InterfaceKey = "server_http"

type Operator interface {
	Start()

	HandleGetFile(serverPath, localPath string, mimeType *string) error
	HandleGetString(serverPath, str string, mimeType *string)
	HandleRaw(endpoint router.Endpoint, rawHandler RawHandler, allowedIDs []auth.ID)
	HandleHTML(endpoint router.Endpoint, htmlHandler HTMLHandler, allowedIDs []auth.ID)
	HandleTemplatorHTML(templatorHTML Templator)
	HandleREST(endpoint router.Endpoint, restHandler RESTHandler, allowedIDs []auth.ID)
	HandleBinary(endpoint router.Endpoint, binaryHandler BinaryHandler, allowedIDs []auth.ID)
	HandleWorker(endpoint router.Endpoint, worker router.Worker, allowedIDs []auth.ID)
}

// !!! requires internal variables (so it can't be a simple function only)
type Templator interface {
	Context(*auth.User, *http.Request, basis.Params) map[string]string
}

type RawHandler func(*auth.User, *http.Request, basis.Params, http.ResponseWriter) error
type BinaryHandler func(*auth.User, *http.Request, basis.Params) (server.BinaryResponse, error)
type RESTHandler func(*auth.User, *http.Request, basis.Params) (server.DataResponse, error)
type HTMLHandler func(*auth.User, *http.Request, basis.Params) (HTMLResponse, error)

func InitEndpoints(op Operator, endpoints map[string]router.Endpoint, htmlHandlers map[string]HTMLHandler, restHandlers map[string]RESTHandler,
	binaryHandlers map[string]BinaryHandler, allowedIDs []auth.ID) basis.Errors {
	var errs basis.Errors

	for key, ep := range endpoints {
		if htmlHandler, ok := htmlHandlers[key]; ok {
			op.HandleHTML(ep, htmlHandler, allowedIDs)
		} else if restHandler, ok := restHandlers[key]; ok {
			op.HandleREST(ep, restHandler, allowedIDs)
		} else if binaryHandler, ok := binaryHandlers[key]; ok {
			op.HandleBinary(ep, binaryHandler, allowedIDs)
		} else {
			errs = append(errs, errors.New("no handler for endpoint: "+key))
		}
	}

	return errs
}
