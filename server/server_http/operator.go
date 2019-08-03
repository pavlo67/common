package server_http

import (
	"net/http"

	"github.com/pavlo67/constructor/auth"
	"github.com/pavlo67/constructor/basis"
	"github.com/pavlo67/constructor/server"
	"github.com/pavlo67/constructor/starter/joiner"
)

const InterfaceKey joiner.InterfaceKey = "server_http"

type WorkerHTTP func(*auth.User, basis.Params, *http.Request) (server.Response, error)

type Operator interface {
	HandleEndpoint(endpoint Endpoint) error
	HandleFiles(serverPath, localPath string, mimeType *string) error

	Start()
}

func InitEndpoints(op Operator, endpoints []Endpoint) basis.Errors {
	var errs basis.Errors

	for _, ep := range endpoints {
		errs = errs.Append(op.HandleEndpoint(ep))
	}

	return errs
}
