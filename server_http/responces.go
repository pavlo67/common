package server_http

import (
	"net/http"

	"github.com/pavlo67/punctum/basis"
)

func Redirect(w http.ResponseWriter, req *http.Request) {
	// remove/add not default ports from req.Host
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

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

type RESTDataMessage struct {
	Info     string `json:"info,omitempty"`
	Redirect string `json:"redirect,omitempty"`
}

type RESTDataError struct {
	Error basis.Errors `json:"error,omitempty"`
}

type RESTResponse struct {
	Status int         `json:"status,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func RESTError(err error) RESTResponse {
	return RESTResponse{
		Status: http.StatusOK,
		Data:   RESTDataError{basis.Errors{err}},
	}
}

type BinaryResponse struct {
	Status   int
	MIMEType string
	Data     []byte
	FileName string
}
