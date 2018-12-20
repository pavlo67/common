package node_http_jschmhr

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/node"
)

// node_http_jschmhr ----------------------------------------------------------------------------------------------------------

type node_http_jschmhr struct {
	httpServer   *http.Server
	httpServeMux *httprouter.Router

	certFileTLS string
	keyFileTLS  string
}

func New(port int, certFileTLS, keyFileTLS string) (node.Operator, error) {

	if port <= 0 {
		return nil, errors.Errorf("serverOp hasn't started: no correct data for http port: %d", port)
	}

	router := httprouter.New()

	return &node_http_jschmhr{
		httpServer: &http.Server{
			Addr:           ":" + strconv.Itoa(port),
			Handler:        router,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		httpServeMux: router,

		certFileTLS: certFileTLS,
		keyFileTLS:  keyFileTLS,
	}, nil
}

type Handler = httprouter.Handle

func redirect(w http.ResponseWriter, req *http.Request) {
	// remove/add not default ports from req.Host
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

// start wraps and verbalizes http.Server.ListenAndServe method.
func (s *node_http_jschmhr) Start() {
	l.Info("Server is starting on address", s.httpServer.Addr)
	if s.certFileTLS != "" && s.keyFileTLS != "" {
		go http.ListenAndServe(":80", http.HandlerFunc(redirect))
		l.Info(s.httpServer.ListenAndServeTLS(s.certFileTLS, s.keyFileTLS))
	} else {
		l.Info(s.httpServer.ListenAndServe())
	}
	// ??? panic on serverOp fault
}

// HandleFunc wraps and verbalizes node_http_jschmhr.Handler.HandleFunc method.
func (s *node_http_jschmhr) Handle(method, path string, handler Handler) {
	if handler == nil {
		l.Error(method, " --> ", path, "\t!!! NULL HANDLER ISN'T DISPATCHED !!!")
		return
	}
	l.Info(method, " --> ", path)
	switch strings.ToLower(method) {
	case "get":
		s.httpServeMux.GET(path, handler)
	case "post":
		s.httpServeMux.POST(path, handler)
	default:
		l.Error(method, " isn't supported!")
	}
}
