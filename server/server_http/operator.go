package server_http

import (
	"net/http"

	"github.com/pavlo67/punctum/auth"
	"github.com/pavlo67/punctum/server"
	"github.com/pavlo67/punctum/starter/joiner"
)

const InterfaceKey joiner.InterfaceKey = "server_http"

type Operator interface {
	Start()

	HandleFile(serverRoute, localPath string, mimeType *string) error
	HandleString(serverRoute, str string, mimeType *string)
	HandleFuncRaw(method, serverRoute string, rawHandler RawHandler, allowedIDs ...auth.ID)
	HandleFuncHTML(method, serverRoute string, htmlHandler HTMLHandler, allowedIDs ...auth.ID)
	HandleTemplatorHTML(templatorHTML Templator)
	HandleFuncREST(method, serverRoute string, restHandler RESTHandler, allowedIDs ...auth.ID)
	HandleFuncBinary(method, serverRoute string, binaryHandler BinaryHandler, allowedIDs ...auth.ID)
}

// !!! requires internal variables (so it can't be a simple function only)
type Templator interface {
	Context(*auth.User, *http.Request, server.RouteParams) map[string]string
}

type RawHandler func(*auth.User, *http.Request, server.RouteParams, http.ResponseWriter) error
type BinaryHandler func(*auth.User, *http.Request, server.RouteParams) (server.BinaryResponse, error)
type HTMLHandler func(*auth.User, *http.Request, server.RouteParams) (HTMLResponse, error)
type RESTHandler func(*auth.User, *http.Request, server.RouteParams) (RESTResponse, error)
