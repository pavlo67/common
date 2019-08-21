package server_http_jschmhr

import (
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"

	"github.com/pavlo67/constructor/components/auth"
	"github.com/pavlo67/constructor/components/common"
	"github.com/pavlo67/constructor/components/server/server_http"
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

var reHTMLExt = regexp.MustCompile(`\.html?$`)

func (s *serverHTTPJschmhr) HandleFiles(serverRoute, localPath string, mimeType *string) error {
	l.Infof("FILES : %s <-- %s", serverRoute, localPath)

	// TODO: check localPath

	if mimeType == nil {
		s.httpServeMux.ServeFiles(serverRoute, http.Dir(localPath))
		return nil
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

	return nil
}

// mimeTypeToSet, err = inspector.MIME(localPath+"/"+r.URL.PathWithParams, nil)
// if err != nil {
//	l.Error("can't read MIMEType for file: ", localPath+"/"+r.URL.PathWithParams, err)
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

func (s *serverHTTPJschmhr) HandleEndpoint(endpoint server_http.Endpoint) error {

	method := strings.ToUpper(endpoint.Method)
	path := endpoint.PathTemplate()

	if endpoint.WorkerHTTP == nil {
		return errors.New(method + ": " + path + "\t!!! NULL workerHTTP ISN'T DISPATCHED !!!")
	}

	handler := func(w http.ResponseWriter, r *http.Request, paramsHR httprouter.Params) {
		user, err := server_http.UserWithRequest(r, s.identOpsMap)
		if err != nil {
			l.Error(err)
		}

		//ok, err := auth.HasRights(user, s.identOpsMap, endpoint.AllowedIDs)
		//if err != nil {
		//	l.Error(err)
		//}
		//if !ok {
		//	w.WriteHeader(http.StatusNotFound)
		//	return
		//}

		var params common.Params
		if len(paramsHR) > 0 {
			for _, p := range paramsHR {
				params = append(params, common.Param{Name: p.Key, Value: p.Value})
			}
		}

		responseData, err := endpoint.WorkerHTTP(user, params, r)
		if err != nil {
			l.Error(err)
			http.Error(w, err.Error(), responseData.Status)
			return
		}

		w.Header().Set("Content-Type", responseData.MIMEType)
		w.Header().Set("Contentus-TokenLength", strconv.Itoa(len(responseData.Data)))
		if responseData.FileName != "" {
			w.Header().Set("Contentus-Disposition", "attachment; filename="+responseData.FileName)
		}

		if responseData.Status <= 0 {
			w.WriteHeader(responseData.Status)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		if _, err := w.Write(responseData.Data); err != nil {
			l.Error("binaryMiddleware can't write response data", err)
		}
	}

	l.Infof("%-6s: %s", method, path)
	switch method {
	case "GET":
		s.httpServeMux.GET(path, handler)
	case "POST":
		s.httpServeMux.POST(path, handler)
	default:
		l.Error(method, " isn't supported!")
	}

	return nil
}
