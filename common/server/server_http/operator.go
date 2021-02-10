package server_http

import (
	"net/http"

	"github.com/pavlo67/common/common"

	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/server"
)

const OnRequestInterfaceKey common.InterfaceKey = "server_http_on_request"
const InterfaceKey common.InterfaceKey = "server_http"

type Params map[string]string
type WorkerHTTP func(Operator, *http.Request, Params, *crud.Options) (server.Response, error)

type OnRequestMiddleware interface {
	Options(r *http.Request) (*crud.Options, error)
}

type StaticPath struct {
	LocalPath string
	MIMEType  *string
}

type Operator interface {
	ResponseRESTOk(status int, data interface{}) (server.Response, error)
	ResponseRESTError(status int, err error, req *http.Request) (server.Response, error)
	HandleEndpoint(key common.InterfaceKey, serverPath string, endpoint Endpoint) error
	HandleFiles(key common.InterfaceKey, serverPath string, staticPath StaticPath) error

	Start() error
	Addr() (port int, https bool)
}
