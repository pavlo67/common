package demo_server_http_jsschmhr

import (
	"net/http"

	"github.com/pavlo67/punctum/identity"
	"github.com/pavlo67/punctum/server_http"
)

var htmlHandlers = map[string]server_http.HTMLHandler{
	"root": root,
}

func root(_ *identity.User, _ *http.Request, _ map[string]string) (server_http.HTMLResponse, error) {
	responseData := server_http.HTMLResponse{
		Data: map[string]string{
			"caput":  "Про цей сервер",
			"title":  "про себе",
			"corpus": "!!!",
		},
	}

	return responseData, nil
}
