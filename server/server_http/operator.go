package server_http

import (
	"errors"
	"net/http"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/config"
	"github.com/pavlo67/punctum/basis/program"
	"github.com/pavlo67/punctum/identity"
	"github.com/pavlo67/punctum/server"
)

const InterfaceKey program.InterfaceKey = "server_http"

type Operator interface {
	server.Operator

	HandleFile(serverPath, localPath string, mimeType *string) error
	HandleString(serverPath, str string, mimeType *string)
	HandleFuncRaw(method, serverPath string, rawHandler HandlerRaw, allowedIDs ...basis.ID)
	HandleFuncHTML(method, serverPath string, htmlHandler HandlerHTML, allowedIDs ...basis.ID)
	HandleTemplatorHTML(templatorHTML Templator)
	HandleFuncREST(method, serverPath string, restHandler HandlerREST, allowedIDs ...basis.ID)
	HandleFuncBinary(method, serverPath string, binaryHandler HandlerBinary, allowedIDs ...basis.ID)
}

type Templator func(*identity.User, *http.Request, map[string]string) map[string]string
type HandlerRaw func(*identity.User, *http.Request, map[string]string, http.ResponseWriter) error
type HandlerHTML func(*identity.User, *http.Request, map[string]string) (HTMLResponse, error)
type HandlerREST func(*identity.User, *http.Request, map[string]string) (RESTResponse, error)
type HandlerBinary func(*identity.User, *http.Request, map[string]string) (BinaryResponse, error)

func InitEndpoints(op Operator, endpoints map[string]config.Endpoint, htmlHandlers map[string]HandlerHTML, restHandlers map[string]HandlerREST, binaryHandlers map[string]HandlerBinary, allowedIDs []basis.ID) basis.Errors {
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
//		return errors.Wrap(basis.ErrNullItem, "no serverhttp.Operator to HandleDir")
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
