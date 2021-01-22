package server_http

import (
	"net/http"

	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/server"
)

const InterfaceKey joiner.InterfaceKey = "server_http"
const PortInterfaceKey joiner.InterfaceKey = "server_http_port"
const NoHTTPSInterfaceKey joiner.InterfaceKey = "server_http_no_https"

type Params map[string]string
type RequestOptions func(r *http.Request) (*crud.Options, error)
type WorkerHTTP func(Operator, *http.Request, Params, *crud.Options) (server.Response, error)

type StaticPath struct {
	LocalPath string
	MIMEType  *string
}

type Operator interface {
	ResponseRESTOk(status int, data interface{}) (server.Response, error)
	ResponseRESTError(status int, err error, req *http.Request) (server.Response, error)
	HandleEndpoint(key, serverPath string, endpoint Endpoint) error
	HandleFiles(key, serverPath string, staticPath StaticPath) error

	// ServerHTTP() *http.Server
	Start() error
}
