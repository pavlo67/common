package server_http

import (
	"net/http"

	"github.com/pavlo67/workshop/common/crud"

	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/server"
)

const InterfaceKey joiner.InterfaceKey = "server_http"
const PortInterfaceKey joiner.InterfaceKey = "server_http_port"
const NoHTTPSInterfaceKey joiner.InterfaceKey = "server_http_no_https"

type Params map[string]string
type WorkerHTTP func(Operator, *crud.Options, Params, *http.Request) (server.Response, error)

type StaticPath struct {
	LocalPath string
	MIMEType  *string
}

type Operator interface {
	ResponseRESTError(options *crud.Options, status int, err error, req ...*http.Request) (server.Response, error)
	ResponseRESTOk(options *crud.Options, data interface{}) (server.Response, error)
	HandleEndpoint(key, serverPath string, endpoint Endpoint) error
	HandleFiles(key, serverPath string, staticPath StaticPath) error

	Start() error
}
