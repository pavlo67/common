package server_http

import (
	"net/http"

	"github.com/pavlo67/punctum/auth"
	"github.com/pavlo67/punctum/starter/joiner"
)

const InterfaceKey joiner.InterfaceKey = "server_http"

type Operator interface {
	Start()

	HandleFile(serverPath, localPath string, mimeType *string) error
	HandleString(serverPath, str string, mimeType *string)
	HandleFuncRaw(method, serverPath string, rawHandler RawHandler, allowedIDs ...auth.ID)
	HandleFuncHTML(method, serverPath string, htmlHandler HTMLHandler, allowedIDs ...auth.ID)
	HandleTemplatorHTML(templatorHTML Templator)
	HandleFuncREST(method, serverPath string, restHandler RESTHandler, allowedIDs ...auth.ID)
	HandleFuncBinary(method, serverPath string, binaryHandler BinaryHandler, allowedIDs ...auth.ID)
}

// !!! requires internal variables (so it can't be a simple function only)
type Templator interface {
	Context(*auth.User, *http.Request, map[string]string) map[string]string
}

type RawHandler func(*auth.User, *http.Request, map[string]string, http.ResponseWriter) error
type HTMLHandler func(*auth.User, *http.Request, map[string]string) (HTMLResponse, error)
type RESTHandler func(*auth.User, *http.Request, map[string]string) (RESTResponse, error)
type BinaryHandler func(*auth.User, *http.Request, map[string]string) (BinaryResponse, error)
