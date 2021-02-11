package server_http

import (
	"net/http"
)

const (
	CORSAllowHeaders     = "authorization,content-type"
	CORSAllowMethods     = "HEAD,GET,POST,PUT,DELETE,OPTIONS"
	CORSAllowOrigin      = "*"
	CORSAllowCredentials = "true"
)

// REST -------------------------------------------------------------------------------------

type RESTDataMessage struct {
	Info     string `json:"info,omitempty"`
	Redirect string `json:"redirect,omitempty"`
}

// Redirect ----------------------------------------------------------------------------------

func Redirect(w http.ResponseWriter, req *http.Request) {
	// remove/add not default ports from req.Host
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}
