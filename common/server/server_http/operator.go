package server_http

import (
	"net/http"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/joiner"

	"github.com/pavlo67/workshop/common/server"
	"github.com/pavlo67/workshop/components/auth"
)

const InterfaceKey joiner.InterfaceKey = "server_http"

type WorkerHTTP func(*auth.User, Params, *http.Request) (server.Response, error)

type Operator interface {
	HandleEndpoint(endpoint Endpoint) error
	HandleFiles(serverPath, localPath string, mimeType *string) error

	Start() error
}

func InitEndpoints(op Operator, endpoints []Endpoint) common.Errors {
	var errs common.Errors

	for _, ep := range endpoints {
		errs = errs.Append(op.HandleEndpoint(ep))
	}

	return errs
}
