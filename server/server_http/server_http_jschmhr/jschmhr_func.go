package server_http_jschmhr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cbroglie/mustache"
	"github.com/julienschmidt/httprouter"

	"strings"

	"github.com/pavlo67/punctum/auth"
	"github.com/pavlo67/punctum/server"
	"github.com/pavlo67/punctum/server/router"
	"github.com/pavlo67/punctum/server/server_http"
)

func ServerPath(ep router.Endpoint) string {
	if len(ep.ParamNames) < 1 {
		return ep.ServerPath
	}
	return ep.ServerPath + "/:" + strings.Join(ep.ParamNames, "/:")
}

func (s *serverHTTPJschmhr) HandleRaw(endpoint router.Endpoint, rawHandler server_http.RawHandler, allowedIDs []auth.ID) {
	l.Fatal("func (s *serverHTTPJschmhr) HandleFuncRaw() isn't implemented!!!")
}

func (s *serverHTTPJschmhr) HandleHTML(endpoint router.Endpoint, htmlHandler server_http.HTMLHandler, allowedIDs []auth.ID) {
	method := endpoint.Method
	serverPath := ServerPath(endpoint)
	s.handleFunc(method, serverPath, func(w http.ResponseWriter, r *http.Request, paramsHR httprouter.Params) {
		user, err := server_http.UserWithRequest(r, s.identOpsMap)
		if err != nil {
			l.Error(err)
		}

		var params router.Params
		if len(paramsHR) > 0 {
			for _, p := range paramsHR {
				params = append(params, router.Param{Name: p.Key, Value: p.Value})
			}
		}

		var context map[string]string
		if s.templator != nil {
			context = s.templator.Context(user, r, params)
		}

		ok, err := auth.HasRights(user, s.identOpsMap, allowedIDs)
		if err != nil {
			l.Error(err)
		}
		if !ok {
			w.Header().Set("Content-Type", "text/html")

			res, err := mustache.Render(s.htmlTemplate, context)
			if err != nil {
				l.Error(err)
			}
			fmt.Fprint(w, res)
			return
		}

		responseData, err := htmlHandler(user, r, params)
		if err != nil {
			l.Error(err)
		}

		if context == nil && len(responseData.Data) > 0 {
			context = map[string]string{}
		}
		for k, v := range responseData.Data {
			context[k] = v
		}

		res, err := mustache.Render(s.htmlTemplate, context)
		if err != nil {
			l.Error(err)
		}

		w.Header().Set("Content-Type", "text/html")
		if responseData.Status > 0 {
			w.WriteHeader(responseData.Status)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		if _, err := w.Write([]byte(res)); err != nil {
			l.Error("htmlMiddleware can't write response data", err)
		}
	})
}

func (s *serverHTTPJschmhr) HandleTemplatorHTML(templatorHTML server_http.Templator) {
	s.templator = templatorHTML
}

func (s *serverHTTPJschmhr) HandleREST(endpoint router.Endpoint, restHandler server_http.RESTHandler, allowedIDs []auth.ID) {
	method := endpoint.Method
	serverPath := ServerPath(endpoint)
	s.handleFunc(method, serverPath, func(w http.ResponseWriter, r *http.Request, paramsHR httprouter.Params) {
		user, err := server_http.UserWithRequest(r, s.identOpsMap)
		if err != nil {
			l.Error(err)
		}

		ok, err := auth.HasRights(user, s.identOpsMap, allowedIDs)
		if err != nil {
			l.Error(err)
		}
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var params router.Params
		if len(paramsHR) > 0 {
			for _, p := range paramsHR {
				params = append(params, router.Param{Name: p.Key, Value: p.Value})
			}
		}

		responseData, err := restHandler(user, r, params)
		if err != nil {
			l.Error(err)
		}

		jsonBytes, err := json.Marshal(responseData.Data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if responseData.Status > 0 {
			w.WriteHeader(responseData.Status)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		if _, err := w.Write(jsonBytes); err != nil {
			l.Error("restMiddleware can't write response data", err)
		}
	})

}

func (s *serverHTTPJschmhr) HandleWorker(endpoint router.Endpoint, workerFunc router.WorkerFunc, allowedIDs []auth.ID) {
	if workerFunc == nil {
		l.Errorf("nil worker for endpoint %#v", endpoint)
		return
	}

	var restHandler = func(user *auth.User, r *http.Request, params router.Params) (server.DataResponse, error) {
		var body []byte
		_, err := r.Body.Read(body)
		if err != nil {
			return server.DataResponse{}, err
		}

		return workerFunc(params, body)
	}

	s.HandleREST(endpoint, restHandler, allowedIDs)
}

func (s *serverHTTPJschmhr) HandleBinary(endpoint router.Endpoint, binaryHandler server_http.BinaryHandler, allowedIDs []auth.ID) {
	method := endpoint.Method
	serverPath := ServerPath(endpoint)
	s.handleFunc(method, serverPath, func(w http.ResponseWriter, r *http.Request, paramsHR httprouter.Params) {
		user, err := server_http.UserWithRequest(r, s.identOpsMap)
		if err != nil {
			l.Error(err)
		}

		ok, err := auth.HasRights(user, s.identOpsMap, allowedIDs)
		if err != nil {
			l.Error(err)
		}
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var params router.Params
		if len(paramsHR) > 0 {
			for _, p := range paramsHR {
				params = append(params, router.Param{Name: p.Key, Value: p.Value})
			}
		}

		responseData, err := binaryHandler(user, r, params)
		if err != nil {
			l.Error(err)
			http.Error(w, err.Error(), responseData.Status)
			return
		}

		w.Header().Set("Content-Type", responseData.MIMEType)
		w.Header().Set("Contentus-TokenLength", strconv.Itoa(len(responseData.Data)))
		if responseData.FileName != "" {
			w.Header().Set("Contentus-Disposition", "attachment; filename="+responseData.FileName)
		}

		if responseData.Status <= 0 {
			w.WriteHeader(responseData.Status)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		if _, err := w.Write(responseData.Data); err != nil {
			l.Error("binaryMiddleware can't write response data", err)
		}
	})
}
