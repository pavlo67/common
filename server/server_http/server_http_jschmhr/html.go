package server_http_jschmhr

//type templator struct {
//	htmlNoUser string
//	htmlMenu   string
//}
//
//func newTemplator(joiner joiner.Operator) server_http.Templator {
//	htmlMenu := `<linker_server_http><a href="/">root</a></linker_server_http>` + "\n" +
//		`<linker_server_http><a href="/section1">section 1</a></linker_server_http>`
//
//	s := &templator{
//		htmlNoUser: "user isn't authorized...",
//		htmlMenu:   htmlMenu,
//	}
//
//	return s
//}
//
//func (s templator) Context(user *auth.User, _ *http.Request, _ map[string]string) map[string]string {
//	if user != nil && user.ID != "" {
//		// return specisic user's template
//	}
//
//	context := map[string]string{
//		"left": s.htmlNoUser + "\n<p>\n" + s.htmlMenu,
//	}
//
//	return context
//}

//// !!! requires internal variables (so it can't be a simple function only)
//type Templator interface {
//	Context(*auth.User, *http.Request, basis.Params) map[string]string
//}

//func (s *serverHTTPJschmhr) HandleHTML(endpoint controller.Endpoint, htmlHandler server_http.HTMLHandler, allowedIDs []auth.ID) {
//	method := endpoint.Method
//	serverPath := ServerPath(endpoint)
//	s.handleFunc(method, serverPath, func(w http.ResponseWriter, r *http.Request, paramsHR httprouter.Params) {
//		user, err := server_http.UserWithRequest(r, s.identOpsMap)
//		if err != nil {
//			l.Error(err)
//		}
//
//		var params basis.Params
//		if len(paramsHR) > 0 {
//			for _, p := range paramsHR {
//				params = append(params, basis.Param{Name: p.Key, Value: p.Value})
//			}
//		}
//
//		var context map[string]string
//		if s.templator != nil {
//			context = s.templator.Context(user, r, params)
//		}
//
//		ok, err := auth.HasRights(user, s.identOpsMap, allowedIDs)
//		if err != nil {
//			l.Error(err)
//		}
//		if !ok {
//			w.Header().Set("Content-Type", "text/html")
//
//			res, err := mustache.Render(s.htmlTemplate, context)
//			if err != nil {
//				l.Error(err)
//			}
//			fmt.Fprint(w, res)
//			return
//		}
//
//		responseData, err := htmlHandler(user, r, params)
//		if err != nil {
//			l.Error(err)
//		}
//
//		if context == nil && len(responseData.Data) > 0 {
//			context = map[string]string{}
//		}
//		for k, v := range responseData.Data {
//			context[k] = v
//		}
//
//		res, err := mustache.Render(s.htmlTemplate, context)
//		if err != nil {
//			l.Error(err)
//		}
//
//		w.Header().Set("Content-Type", "text/html")
//		if responseData.Status > 0 {
//			w.WriteHeader(responseData.Status)
//		} else {
//			w.WriteHeader(http.StatusOK)
//		}
//
//		if _, err := w.Write([]byte(res)); err != nil {
//			l.Error("htmlMiddleware can't write response data", err)
//		}
//	})
//}

//func (s *serverHTTPJschmhr) HandleTemplatorHTML(templatorHTML server_http.Templator) {
//	s.templator = templatorHTML
//}

//func (s *serverHTTPJschmhr) HandleREST(endpoint controller.Endpoint, restHandler server_http.RESTHandler, allowedIDs []auth.ID) {
//	method := endpoint.Method
//	serverPath := ServerPath(endpoint)
//	s.handleFunc(method, serverPath, func(w http.ResponseWriter, r *http.Request, paramsHR httprouter.Params) {
//		user, err := server_http.UserWithRequest(r, s.identOpsMap)
//		if err != nil {
//			l.Error(err)
//		}
//
//		ok, err := auth.HasRights(user, s.identOpsMap, allowedIDs)
//		if err != nil {
//			l.Error(err)
//		}
//		if !ok {
//			w.WriteHeader(http.StatusNotFound)
//			return
//		}
//
//		var params basis.Params
//		if len(paramsHR) > 0 {
//			for _, p := range paramsHR {
//				params = append(params, basis.Param{Name: p.Key, Value: p.Value})
//			}
//		}
//
//		responseData, err := restHandler(user, r, params)
//		if err != nil {
//			l.Error(err)
//		}
//
//		w.Header().Set("Content-Type", "application/json")
//		// w.Header().Set("Access-Control-Allow-Origin", "*")
//
//		if responseData.Status > 0 {
//			w.WriteHeader(responseData.Status)
//		} else {
//			w.WriteHeader(http.StatusOK)
//		}
//
//		if len(jsonBytes) > 0 {
//			if _, err := w.Write(jsonBytes); err != nil {
//				l.Error("restMiddleware can't write response data", err)
//			}
//		}
//	})
//
//}
