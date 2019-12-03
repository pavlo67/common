package server_http

import (
	"net/http"

	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/server"
)

const InterfaceKey joiner.InterfaceKey = "server_http"

type WorkerHTTP func(*auth.User, Params, *http.Request) (server.Response, error)

type Operator interface {
	HandleEndpoint(serverPath string, endpoint Endpoint) error
	HandleFiles(serverPath string, staticPath StaticPath) error

	Start() error
}

type StaticPath struct {
	LocalPath string
	MIMEType  *string
}
