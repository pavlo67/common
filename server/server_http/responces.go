package server_http

import (
	"github.com/pavlo67/constructor/basis"
)

// REST -------------------------------------------------------------------------------------

type RESTDataMessage struct {
	Info     string `json:"info,omitempty"`
	Redirect string `json:"redirect,omitempty"`
}

type RESTDataError struct {
	Error basis.Errors `json:"error,omitempty"`
}

//func RESTError(err error) server.Response {
//	return server.Response{
//		Status: http.StatusOK,
//		Data:   RESTDataError{basis.Errors{err}},
//	}
//}
