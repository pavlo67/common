package demo_server_http

//var htmlHandlers = map[string]server_http.HTMLHandler{
//	"root":  rootHandler,
//	"sect1": section1Handler,
//}
//
//var endpoints = map[string]server_http.Endpoint{
//	"root":  {Method: "GET", ServerPath: "/"},
//	"sect1": {Method: "GET", ServerPath: "/section1"},
//}
//
//func rootHandler(_ *auth.User, _ *http.Request, _ map[string]string) (server_http.HTMLResponse, error) {
//	responseData := server_http.HTMLResponse{
//		Data: map[string]string{
//			"caput":  "Про цей сервер",
//			"title":  "про себе",
//			"corpus": "!!!",
//		},
//	}
//
//	return responseData, nil
//}
//
//func section1Handler(_ *auth.User, _ *http.Request, _ map[string]string) (server_http.HTMLResponse, error) {
//	responseData := server_http.HTMLResponse{
//		Data: map[string]string{
//			"caput":  "Розділ 1",
//			"title":  "перший розділ сервера",
//			"corpus": "???",
//		},
//	}
//
//	return responseData, nil
//}
