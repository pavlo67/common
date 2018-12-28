package demo_server_http

import (
	"net/http"

	"github.com/pavlo67/punctum/basis/config"
	"github.com/pavlo67/punctum/identity"
	"github.com/pavlo67/punctum/server_http"
)

var htmlHandlers = map[string]server_http.HTMLHandler{
	"root":  rootHandler,
	"sect1": section1Handler,
}

var endpoints = map[string]config.Endpoint{
	"root":  {Method: "GET", ServerPath: "/"},
	"sect1": {Method: "GET", ServerPath: "/section1"},
}

func rootHandler(_ *identity.User, _ *http.Request, _ map[string]string) (server_http.HTMLResponse, error) {
	responseData := server_http.HTMLResponse{
		Data: map[string]string{
			"caput":  "Про цей сервер",
			"title":  "про себе",
			"corpus": "!!!",
		},
	}

	return responseData, nil
}

func section1Handler(_ *identity.User, _ *http.Request, _ map[string]string) (server_http.HTMLResponse, error) {
	responseData := server_http.HTMLResponse{
		Data: map[string]string{
			"caput":  "Розділ 1",
			"title":  "перший розділ сервера",
			"corpus": "???",
		},
	}

	return responseData, nil
}
