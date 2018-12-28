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

	"github.com/pavlo67/punctum/identity"
	"github.com/pavlo67/punctum/server_http"
)

var _ server_http.Operator = &serverHTTPJschmhr{}

type serverHTTPJschmhr struct {
	httpServer   *http.Server
	httpServeMux *httprouter.Router
	certFileTLS  string
	keyFileTLS   string
	identOpsMap  map[identity.CredsType][]identity.Operator

	htmlTemplate string
	templator    server_http.Templator
}

func New(port int, certFileTLS, keyFileTLS string, identOpsMap map[identity.CredsType][]identity.Operator, htmlTemplate string) (server_http.Operator, error) {
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

		htmlTemplate: htmlTemplate,
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

func (s *serverHTTPJschmhr) HandleFile(serverPath, localPath string, mimeType *string) error {
	l.Infof("FILES : %s <-- %s", serverPath, localPath)

	// TODO: check localPath

	if mimeType == nil {
		s.httpServeMux.ServeFiles(serverPath, http.Dir(localPath))
		return nil
	}

	//fileServer := http.FileServer(http.Dir(localPath))
	s.httpServeMux.GET(serverPath, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

	return nil
}

// mimeTypeToSet, err = inspector.MIME(localPath+"/"+r.URL.Path, nil)
// if err != nil {
//	l.Error("can't read MIMEType for file: ", localPath+"/"+r.URL.Path, err)
// }

func (s *serverHTTPJschmhr) HandleString(serverPath, str string, mimeType *string) {
	s.handleFunc("GET", serverPath, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		if mimeType != nil {
			// "application/javascript"
			w.Header().Set("Content-Type", *mimeType)
		}
		w.Write([]byte(str))
	})
}
