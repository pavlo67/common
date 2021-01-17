package server_http

import (
	"net/http"

	"github.com/pavlo67/workshop/common/auth"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/server"
)

const InterfaceKey joiner.InterfaceKey = "server_http"
const PortInterfaceKey joiner.InterfaceKey = "server_http_port"
const NoHTTPSInterfaceKey joiner.InterfaceKey = "server_http_no_https"

type Params map[string]string
type WorkerHTTP func(Operator, *auth.Identity, Params, *http.Request) (server.Response, error)

type StaticPath struct {
	LocalPath string
	MIMEType  *string
}

type Operator interface {
	ResponseRESTError(identity *auth.Identity, status int, err common.Error, req ...*http.Request) (server.Response, error)
	ResponseRESTOk(identity *auth.Identity, data interface{}) (server.Response, error)
	HandleEndpoint(key, serverPath string, endpoint Endpoint) error
	HandleFiles(key, serverPath string, staticPath StaticPath) error

	Start() error
}
