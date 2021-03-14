package server_http

import (
	"net/http"

	"github.com/pavlo67/common/common/auth"

	"github.com/pavlo67/common/common/joiner"

	"github.com/pavlo67/common/common/server"
)

const OnRequestMiddlewareInterfaceKey joiner.InterfaceKey = "server_http_on_request_middleware"
const InterfaceKey joiner.InterfaceKey = "server_http"

type PathParams map[string]string
type WorkerHTTP func(Operator, *http.Request, PathParams, *auth.Identity) (server.Response, error)

type OnRequestMiddleware interface {
	Identity(r *http.Request) (*auth.Identity, error)
}

type StaticPath struct {
	LocalPath string
	MIMEType  *string
}

type Operator interface {
	ResponseRESTOk(status int, data interface{}, req *http.Request) (server.Response, error)
	ResponseRESTError(status int, err error, req *http.Request) (server.Response, error)
	HandleEndpoint(key joiner.InterfaceKey, serverPath string, endpoint Endpoint) error
	HandleFiles(key joiner.InterfaceKey, serverPath string, staticPath StaticPath) error

	Start() error
	Addr() (port int, https bool)
}
