package server_http_jschmhr

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/pavlo67/associatio/auth"
	"github.com/pavlo67/associatio/basis"
	"github.com/pavlo67/associatio/server"
	"github.com/pavlo67/associatio/server/controller"
	"github.com/pavlo67/associatio/server/server_http"
)

func ServerPath(ep controller.Endpoint) string {
	path := ep.Path
	if len(path) == 0 || path[0] != '/' {
		path = "/" + path
	}

	if len(ep.ParamNames) < 1 {
		return path
	}

	return path + "/:" + strings.Join(ep.ParamNames, "/:")
}

func (s *serverHTTPJschmhr) HandleHTTP(endpoint controller.Endpoint, workerHTTP server_http.WorkerHTTP) {
	method := endpoint.Method
	serverPath := ServerPath(endpoint)
	s.handleFunc(method, serverPath, func(w http.ResponseWriter, r *http.Request, paramsHR httprouter.Params) {
		user, err := server_http.UserWithRequest(r, s.identOpsMap)
		if err != nil {
			l.Error(err)
		}

		ok, err := auth.HasRights(user, s.identOpsMap, endpoint.AllowedIDs)
		if err != nil {
			l.Error(err)
		}
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var params basis.Params
		if len(paramsHR) > 0 {
			for _, p := range paramsHR {
				params = append(params, basis.Param{Name: p.Key, Value: p.Value})
			}
		}

		responseData, err := workerHTTP(user, params, r)
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

func (s *serverHTTPJschmhr) HandleWorker(endpoint controller.Endpoint, worker controller.Worker) {
	if worker == nil {
		l.Errorf("nil worker for endpoint %#v", endpoint)
		return
	}

	var handler = func(user *auth.User, params basis.Params, r *http.Request) (server.Response, error) {
		var data interface{}
		options := basis.Info{}
		for k, v := range r.URL.Query() {
			options[k] = v
		}

		if endpoint.Method != "" && strings.ToUpper(endpoint.Method) != "GET" {
			var body []byte
			_, err := r.Body.Read(body)
			if err != nil {
				return server.Response{}, err
			}

			if endpoint.DataItem != nil {
				data = reflect.New(reflect.ValueOf(endpoint.DataItem).Elem().Type()).Interface()
				if err = json.Unmarshal(body, data); err != nil {
					return server.Response{}, err
				}
			}
		}

		return worker(user, params, data)
	}

	s.HandleHTTP(endpoint, handler)
}
