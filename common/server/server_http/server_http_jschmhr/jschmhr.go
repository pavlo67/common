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

	"github.com/pavlo67/workshop/common/server/server_http"
	"github.com/pavlo67/workshop/components/auth"
)

var _ server_http.Operator = &serverHTTPJschmhr{}

type serverHTTPJschmhr struct {
	httpServer   *http.Server
	httpServeMux *httprouter.Router

	port        int
	certFileTLS string
	keyFileTLS  string
	authOps     []auth.Operator
}

func New(port int, certFileTLS, keyFileTLS string, authOps []auth.Operator) (server_http.Operator, error) {
	if port <= 0 {
		return nil, errors.Errorf("on server_http_jschmhr.New(): wrong port = %d", port)
	}

	router := httprouter.New()

	return &serverHTTPJschmhr{
		httpServer: &http.Server{
			Handler:        router,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   60 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		httpServeMux: router,

		port: port,

		certFileTLS: certFileTLS,
		keyFileTLS:  keyFileTLS,

		authOps: authOps,
	}, nil
}

// start wraps and verbalizes http.Server.ListenAndServe method.
func (s *serverHTTPJschmhr) Start() error {
	if s == nil {
		return errors.Errorf("no serverOp to start")
	}

	s.httpServer.Addr = ":" + strconv.Itoa(s.port)

	l.Info("Server is starting on address ", s.httpServer.Addr)

	if s.certFileTLS != "" && s.keyFileTLS != "" {
		go http.ListenAndServe(":80", http.HandlerFunc(server_http.Redirect))
		return s.httpServer.ListenAndServeTLS(s.certFileTLS, s.keyFileTLS)
	}

	return s.httpServer.ListenAndServe()
}

var reHTMLExt = regexp.MustCompile(`\.html?$`)

func (s *serverHTTPJschmhr) HandleFiles(serverRoute, localPath string, mimeType *string) error {
	l.Infof("FILES : %s <-- %s", serverRoute, localPath)

	// TODO: check localPath

	if mimeType == nil {
		// TODO!!! CORS

		s.httpServeMux.ServeFiles(serverRoute, http.Dir(localPath))
		return nil
	}

	s.httpServeMux.OPTIONS(serverRoute, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		l.Infof("OPTIONS: %s", serverRoute)

		w.Header().Set("Access-Control-Allow-Origin", server_http.CORSAllowOrigin)
		w.Header().Set("Access-Control-Allow-Headers", server_http.CORSAllowHeaders)
		w.Header().Set("Access-Control-Allow-Methods", server_http.CORSAllowMethods)
		w.Header().Set("Access-Control-Allow-Credentials", server_http.CORSAllowCredentials)
		w.Header().Set("Content-Type", *mimeType)
	})

	//fileServer := http.FileServer(http.Dir(localPath))
	s.httpServeMux.GET(serverRoute, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", server_http.CORSAllowOrigin)
		w.Header().Set("Access-Control-Allow-Headers", server_http.CORSAllowHeaders)
		w.Header().Set("Access-Control-Allow-Methods", server_http.CORSAllowMethods)
		w.Header().Set("Access-Control-Allow-Credentials", server_http.CORSAllowCredentials)
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
		user, err := server_http.UserWithRequest(r, s.authOps)
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

		var params server_http.Params
		if len(paramsHR) > 0 {
			params = server_http.Params{}
			for _, p := range paramsHR {
				params[p.Key] = p.Value
			}
		}

		//var params server_http.Params
		//if len(paramsHR) > 0 {
		//	for _, p := range paramsHR {
		//		params = append(params, server_http.Param{Name: p.Key, Value: p.Value})
		//	}
		//}

		w.Header().Set("Access-Control-Allow-Origin", server_http.CORSAllowOrigin)
		w.Header().Set("Access-Control-Allow-Headers", server_http.CORSAllowHeaders)
		w.Header().Set("Access-Control-Allow-Methods", server_http.CORSAllowMethods)
		w.Header().Set("Access-Control-Allow-Credentials", server_http.CORSAllowCredentials)

		responseData, err := endpoint.WorkerHTTP(user, params, r)
		if err != nil {
			l.Error(err)
			http.Error(w, err.Error(), responseData.Status)
			return
		}
		w.Header().Set("Content-Type", responseData.MIMEType)
		w.Header().Set("Content-Length", strconv.Itoa(len(responseData.Data)))
		if responseData.FileName != "" {
			w.Header().Set("Content-Disposition", "attachment; filename="+responseData.FileName)
		}

		if responseData.Status <= 0 {
			w.WriteHeader(responseData.Status)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		if _, err := w.Write(responseData.Data); err != nil {
			l.Error("can't write response data", err)
		}
	}

	s.httpServeMux.OPTIONS(path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// l.Infof("OPTIONS: %s", path)

		w.Header().Set("Access-Control-Allow-Origin", server_http.CORSAllowOrigin)
		w.Header().Set("Access-Control-Allow-Headers", server_http.CORSAllowHeaders)
		w.Header().Set("Access-Control-Allow-Methods", server_http.CORSAllowMethods)
		w.Header().Set("Access-Control-Allow-Credentials", server_http.CORSAllowCredentials)
		// w.Header().Set("Content-Type", *mimeType)
	})

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
