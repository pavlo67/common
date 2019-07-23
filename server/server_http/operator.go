package server_http

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/auth"
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/server"
	"github.com/pavlo67/punctum/server/controller"
	"github.com/pavlo67/punctum/starter/joiner"
)

const InterfaceKey joiner.InterfaceKey = "server_http"

type Operator interface {
	controller.Operator

	HandleHTTP(endpoint controller.Endpoint, workerHTTP WorkerHTTP)
	HandleFiles(serverPath, localPath string, mimeType *string)
}

type WorkerHTTP func(*auth.User, basis.Params, *http.Request) (server.Response, error)

func InitEndpoints(op Operator, endpoints map[string]controller.Endpoint, workersHTTP map[string]WorkerHTTP) basis.Errors {
	var errs basis.Errors

	for key, ep := range endpoints {
		if workerHTTP, ok := workersHTTP[key]; ok {
			op.HandleHTTP(ep, workerHTTP)
		} else {
			errs = append(errs, errors.New("no handler for endpoint: "+key))
		}
	}

	return errs
}
