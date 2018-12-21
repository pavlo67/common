package server_http_jschmhr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cbroglie/mustache"
	"github.com/julienschmidt/httprouter"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/identity"
	"github.com/pavlo67/punctum/server_http"
)

func (s *server_http_jschmhr) getUserData(r *http.Request) *identity.User {
	if s.identOp == nil {
		return nil
	}

	token := r.Header.Get("Authorization")
	if token == "" {
		c, err := r.Cookie("token")
		if err != nil {
			return nil
		}
		token = c.Value
		if token == "" {
			return nil
		}
	}
	user, _, err := s.identOp.Identify([]identity.Creds{{Type: identity.CredsToken, Value: token}})
	if err != nil {
		l.Error(err)
		return nil
	}
	return user
}

func (s *server_http_jschmhr) HandleFuncRaw(method, serverPath string, rawHandler server_http.HandlerRaw, allowedIDs ...basis.ID) {
	l.Fatal("func (s *server_http_jschmhr) HandleFuncRaw() isn't implemented!!!")
}

func (s *server_http_jschmhr) HandleFuncHTML(method, serverPath string, htmlHandler server_http.HandlerHTML, allowedIDs ...basis.ID) {
	s.handleFunc(method, serverPath, func(w http.ResponseWriter, r *http.Request, paramsHR httprouter.Params) {
		user := s.getUserData(r)

		var params map[string]string
		if len(paramsHR) > 0 {
			params = map[string]string{}
			for _, p := range paramsHR {
				if _, ok := params[p.Key]; !ok {
					params[p.Key] = p.Value
				}
			}
		}

		var context map[string]string
		if s.templator != nil {
			context = s.templator(user, r, params)
		}

		if allowedIDs != nil && !s.identOp.HasRights(user, allowedIDs...) {
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

func (s *server_http_jschmhr) HandleTemplatorHTML(templatorHTML server_http.Templator) {
	s.templator = templatorHTML
}

func (s *server_http_jschmhr) HandleFuncREST(method, serverPath string, restHandler server_http.HandlerREST, allowedIDs ...basis.ID) {
	s.handleFunc(method, serverPath, func(w http.ResponseWriter, r *http.Request, paramsHR httprouter.Params) {
		user := s.getUserData(r)

		if allowedIDs != nil && !s.identOp.HasRights(user, allowedIDs...) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var params map[string]string
		if len(paramsHR) > 0 {
			params = map[string]string{}
			for _, p := range paramsHR {
				if _, ok := params[p.Key]; !ok {
					params[p.Key] = p.Value
				}
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

func (s *server_http_jschmhr) HandleFuncBinary(method, serverPath string, binaryHandler server_http.HandlerBinary, allowedIDs ...basis.ID) {
	s.handleFunc(method, serverPath, func(w http.ResponseWriter, r *http.Request, paramsHR httprouter.Params) {
		user := s.getUserData(r)

		if allowedIDs != nil && !s.identOp.HasRights(user, allowedIDs...) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var params map[string]string
		if len(paramsHR) > 0 {
			params = map[string]string{}
			for _, p := range paramsHR {
				if _, ok := params[p.Key]; !ok {
					params[p.Key] = p.Value
				}
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

		if responseData.Status > 0 {
			w.WriteHeader(responseData.Status)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		if _, err := w.Write(responseData.Data); err != nil {
			l.Error("binaryMiddleware can't write response data", err)
		}
	})

}
