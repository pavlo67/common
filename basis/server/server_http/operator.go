package server_http

import (
	"net/http"

	"github.com/pavlo67/workshop/basis/auth"
	"github.com/pavlo67/workshop/basis/common"
	"github.com/pavlo67/workshop/basis/joiner"
	"github.com/pavlo67/workshop/basis/server"
)

const InterfaceKey joiner.InterfaceKey = "server_http"

type WorkerHTTP func(*auth.User, Params, *http.Request) (server.Response, error)

type Operator interface {
	HandleEndpoint(endpoint Endpoint) error
	HandleFiles(serverPath, localPath string, mimeType *string) error

	Start()
}

func InitEndpoints(op Operator, endpoints []Endpoint) common.Errors {
	var errs common.Errors

	for _, ep := range endpoints {
		errs = errs.Append(op.HandleEndpoint(ep))
	}

	return errs
}
