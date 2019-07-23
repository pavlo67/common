package server_http

import (
	"net/http"

	"github.com/pavlo67/punctum/auth"
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/server"
	"github.com/pavlo67/punctum/server/controller"
	"github.com/pavlo67/punctum/starter/joiner"
	"github.com/pkg/errors"
)

const InterfaceKey joiner.InterfaceKey = "server_http"

type Operator interface {
	Start()

	Handle(endpoint controller.Endpoint, worker controller.Worker)
	HandleHTTP(endpoint controller.Endpoint, workerHTTP WorkerHTTP)
	HandleFiles(serverPath, localPath string, mimeType *string)
}

type WorkerHTTP func(*auth.User, basis.Params, *http.Request) (server.Response, error)

func InitEndpoints(op Operator, endpoints map[string]controller.Endpoint, handlers map[string]WorkerHTTP) basis.Errors {
	var errs basis.Errors

	for key, ep := range endpoints {
		if binaryHandler, ok := handlers[key]; ok {
			op.HandleHTTP(ep, binaryHandler)
		} else {
			errs = append(errs, errors.New("no handler for endpoint: "+key))
		}
	}

	return errs
}
