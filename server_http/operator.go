package server_http

import (
	"errors"
	"net/http"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/config"
	"github.com/pavlo67/punctum/basis/joiner"
	"github.com/pavlo67/punctum/identity"
)

const InterfaceKey joiner.InterfaceKey = "server_http"

type Operator interface {
	Start()

	HandleFile(serverPath, localPath string, mimeType *string) error
	HandleString(serverPath, str string, mimeType *string)
	HandleFuncRaw(method, serverPath string, rawHandler RawHandler, allowedIDs ...identity.ID)
	HandleFuncHTML(method, serverPath string, htmlHandler HTMLHandler, allowedIDs ...identity.ID)
	HandleTemplatorHTML(templatorHTML Templator)
	HandleFuncREST(method, serverPath string, restHandler RESTHandler, allowedIDs ...identity.ID)
	HandleFuncBinary(method, serverPath string, binaryHandler BinaryHandler, allowedIDs ...identity.ID)
}

// !!! requires internal variables (so it can't be a simple function only)
type Templator interface {
	Context(*identity.User, *http.Request, map[string]string) map[string]string
}

type RawHandler func(*identity.User, *http.Request, map[string]string, http.ResponseWriter) error
type HTMLHandler func(*identity.User, *http.Request, map[string]string) (HTMLResponse, error)
type RESTHandler func(*identity.User, *http.Request, map[string]string) (RESTResponse, error)
type BinaryHandler func(*identity.User, *http.Request, map[string]string) (BinaryResponse, error)

func InitEndpoints(op Operator, endpoints map[string]config.Endpoint, htmlHandlers map[string]HTMLHandler, restHandlers map[string]RESTHandler, binaryHandlers map[string]BinaryHandler, allowedIDs []identity.ID) basis.Errors {
	var errs basis.Errors

	for key, ep := range endpoints {
		if ep.Method == "FILE" {
			op.HandleFile(ep.ServerPath, ep.LocalPath, nil)
		} else if htmlHandler, ok := htmlHandlers[key]; ok {
			op.HandleFuncHTML(ep.Method, ep.ServerPath, htmlHandler, allowedIDs...)
		} else if restHandler, ok := restHandlers[key]; ok {
			op.HandleFuncREST(ep.Method, ep.ServerPath, restHandler, allowedIDs...)
		} else if binaryHandler, ok := binaryHandlers[key]; ok {
			op.HandleFuncBinary(ep.Method, ep.ServerPath, binaryHandler, allowedIDs...)
		} else {
			errs = append(errs, errors.New("no handler for endpoint: "+key))
		}
	}

	return errs
}

//// !!! non-recursively
//func HandleDir(srvOp Operator, path, localPath string) error {
//	if srvOp == nil {
//		return errors.Wrap(basis.ErrNull, "no serverhttp.Operator to HandleDir")
//	}
//
//	path = strings.TrimSpace(path)
//	if path == "" || path[len(path)-1] != '/' {
//		path += "/"
//	}
//
//	files, err := ioutil.ReadDir(localPath)
//	if err != nil {
//		return errors.Wrapf(err, "on ioutil.ReadDir(\"%s\")", localPath)
//	}
//
//	var errs basis.Errors
//
//	for _, file := range files {
//		if file.IsDir() {
//			err = srvOp.HandleFile("/"+file.Name()+"/*filepath", localPath+file.Name()+"/")
//			if err != nil {
//				errs = append(errs, err)
//			}
//		}
//	}
//
//	return errs.Err()
//}
