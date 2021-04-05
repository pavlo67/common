package server_http

import (
	"net/http"

	"github.com/pavlo67/common/common/server"

	"github.com/pavlo67/common/common/auth"

	"github.com/pavlo67/common/common/joiner"
)

const OnRequestMiddlewareInterfaceKey joiner.InterfaceKey = "server_http_on_request_middleware"
const InterfaceKey joiner.InterfaceKey = "server_http"

type PathParams map[string]string

type OnRequestMiddleware interface {
	Identity(r *http.Request) (*auth.Identity, error)
}

type StaticPath struct {
	LocalPath string
	MIMEType  *string
}

type WorkerHTTP func(Operator, *http.Request, PathParams, *auth.Identity) (server.Response, error)

type Operator interface {
	HandleEndpoint(key EndpointKey, serverPath string, endpoint Endpoint) error
	HandleFiles(key EndpointKey, serverPath string, staticPath StaticPath) error

	Start() error
	Addr() (port int, https bool)
}
