package server_http

import (
	"net/http"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/server"
)

// HTML -------------------------------------------------------------------------------------

type HTMLResponse struct {
	Status int
	Data   map[string]string
}

func HTMLError(status int, label string) HTMLResponse {
	if status == 0 {
		status = http.StatusInternalServerError
	}

	return HTMLResponse{
		Status: status,
		Data:   map[string]string{"corpus": label},
	}
}

// REST -------------------------------------------------------------------------------------

type RESTDataMessage struct {
	Info     string `json:"info,omitempty"`
	Redirect string `json:"redirect,omitempty"`
}

type RESTDataError struct {
	Error basis.Errors `json:"error,omitempty"`
}

func RESTError(err error) server.DataResponse {
	return server.DataResponse{
		Status: http.StatusOK,
		Data:   RESTDataError{basis.Errors{err}},
	}
}
