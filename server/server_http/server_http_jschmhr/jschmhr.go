package server_http_jschmhr

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"

	"io"
	"os"
	"regexp"

	"github.com/pavlo67/associatio/auth"
	"github.com/pavlo67/associatio/server/server_http"
)

var _ server_http.Operator = &serverHTTPJschmhr{}

type serverHTTPJschmhr struct {
	httpServer   *http.Server
	httpServeMux *httprouter.Router
	certFileTLS  string
	keyFileTLS   string
	identOpsMap  map[auth.CredsType][]auth.Operator

	//htmlTemplate string
	//templator    server_http.Templator
}

func New(port int, certFileTLS, keyFileTLS string, identOpsMap map[auth.CredsType][]auth.Operator) (server_http.Operator, error) {
	if port <= 0 {
		return nil, errors.Errorf("serverOp hasn't started: no correct data for http port: %d", port)
	}

	if len(identOpsMap) < 1 {
		l.Warn("no one identity.Operator for serverHTTPJschmhr.New()")
	}

	router := httprouter.New()

	return &serverHTTPJschmhr{
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

		identOpsMap: identOpsMap,
	}, nil
}

// start wraps and verbalizes http.Server.ListenAndServe method.
func (s *serverHTTPJschmhr) Start() {
	l.Info("Server is starting on address ", s.httpServer.Addr)

	var err error

	if s.certFileTLS != "" && s.keyFileTLS != "" {
		go http.ListenAndServe(":80", http.HandlerFunc(server_http.Redirect))
		err = s.httpServer.ListenAndServeTLS(s.certFileTLS, s.keyFileTLS)
	} else {
		err = s.httpServer.ListenAndServe()
	}

	if err != nil {
		l.Error(err)
	}
}

func (s *serverHTTPJschmhr) handleFunc(method, path string, handler httprouter.Handle) {
	if handler == nil {
		l.Error(method, ": ", path, "\t!!! NULL HANDLER ISN'T DISPATCHED !!!")
		return
	}
	l.Infof("%-6s: %s", method, path)
	switch strings.ToLower(method) {
	case "get":
		s.httpServeMux.GET(path, handler)
	case "post":
		s.httpServeMux.POST(path, handler)
	default:
		l.Error(method, " isn't supported!")
	}
}

var reHTMLExt = regexp.MustCompile(`\.html?$`)

func (s *serverHTTPJschmhr) HandleFiles(serverRoute, localPath string, mimeType *string) {
	l.Infof("FILES : %s <-- %s", serverRoute, localPath)

	// TODO: check localPath

	if mimeType == nil {
		s.httpServeMux.ServeFiles(serverRoute, http.Dir(localPath))
		return
	}

	//fileServer := http.FileServer(http.Dir(localPath))
	s.httpServeMux.GET(serverRoute, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", *mimeType)
		OpenFile, err := os.Open(localPath + "/" + p.ByName("filepath"))
		defer OpenFile.Close()
		if err != nil {
			l.Error(err)
		} else {
			io.Copy(w, OpenFile)
		}

		//if mimeType != nil {
		//}
		//fileServer.ServeHTTP(w, r)
	})

	return
}

// mimeTypeToSet, err = inspector.MIME(localPath+"/"+r.URL.WithParams, nil)
// if err != nil {
//	l.Error("can't read MIMEType for file: ", localPath+"/"+r.URL.WithParams, err)
// }

//func (s *serverHTTPJschmhr) HandleGetString(serverRoute, str string, mimeType *string) {
//	s.handleFunc("GET", serverRoute, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
//		if mimeType != nil {
//			// "application/javascript"
//			w.Header().Set("Content-Type", *mimeType)
//		}
//		w.Write([]byte(str))
//	})
//}
