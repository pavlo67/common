package server_http

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/pavlo67/constructor/auth"
	"github.com/pavlo67/constructor/basis"
	"github.com/pavlo67/constructor/server"
	"github.com/pavlo67/constructor/starter/joiner"
)

const InterfaceKey joiner.InterfaceKey = "server_http"

type Operator interface {
	HandleHTTP(endpoint Endpoint, workerHTTP WorkerHTTP)
	HandleFiles(serverPath, localPath string, mimeType *string)

	Start()
}

type WorkerHTTP func(*auth.User, basis.Params, *http.Request) (server.Response, error)

func InitEndpoints(op Operator, endpoints map[string]Endpoint, workersHTTP map[string]WorkerHTTP) basis.Errors {
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
