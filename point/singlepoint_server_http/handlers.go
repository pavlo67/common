package singlepoint_server_http

import (
	"net/http"

	"github.com/pavlo67/punctum/basis/config"
	"github.com/pavlo67/punctum/identity"
	"github.com/pavlo67/punctum/server_http"
)

var endpoints = map[string]config.Endpoint{
	"point": {
		Method:     "GET",
		ServerPath: "/",
	},
}

var restHandlers = map[string]server_http.RESTHandler{
	"point": pointHandler,
}

func pointHandler(_ *identity.User, _ *http.Request, _ map[string]string) (server_http.RESTResponse, error) {
	responseData := server_http.RESTResponse{
		Data: item,
	}

	return responseData, nil
}
