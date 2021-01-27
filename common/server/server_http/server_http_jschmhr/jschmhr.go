package server_http_jschmhr

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/libraries/strlib"
	"github.com/pavlo67/common/common/server"
	"github.com/pavlo67/common/common/server/server_http"
)

var _ server_http.Operator = &serverHTTPJschmhr{}

type serverHTTPJschmhr struct {
	httpServer   *http.Server
	httpServeMux *httprouter.Router

	port        int
	tlsCertFile string
	tlsKeyFile  string

	onRequest server_http.OnRequest

	secretENVsToLower []string
}

func New(port int, tlsCertFile, tlsKeyFile string, onRequest server_http.OnRequest, secretENVs []string) (server_http.Operator, error) {
	if port <= 0 {
		return nil, fmt.Errorf("on server_http_jschmhr.New(): wrong port = %d", port)
	}

	if onRequest == nil {
		return nil, errata.New("on server_http_jschmhr.New(): no server_http.OnRequest")
	}

	var secretENVsToLower []string
	for _, secretENV := range secretENVs {
		secretENVsToLower = append(secretENVsToLower, strings.ToLower(secretENV))
	}

	router := httprouter.New()

	return &serverHTTPJschmhr{
		httpServer: &http.Server{
			Handler:        router,
			ReadTimeout:    60 * time.Second,
			WriteTimeout:   60 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		httpServeMux: router,
		port:         port,
		tlsCertFile:  tlsCertFile,
		tlsKeyFile:   tlsKeyFile,

		onRequest: onRequest,

		secretENVsToLower: secretENVsToLower,
	}, nil
}

// start wraps and verbalizes http.Server.ListenAndServe method.
func (s *serverHTTPJschmhr) Start() error {
	if s == nil {
		return errata.New("no serverOp to start")
	}

	s.httpServer.Addr = ":" + strconv.Itoa(s.port)
	l.Info("Server is starting on address ", s.httpServer.Addr)

	if s.tlsCertFile != "" && s.tlsKeyFile != "" {
		return s.httpServer.ListenAndServeTLS(s.tlsCertFile, s.tlsKeyFile)
	}

	return s.httpServer.ListenAndServe()
}

//func (s *serverHTTPJschmhr) ServerHTTP() *http.Server {
//	return s.httpServer
//}

func (s *serverHTTPJschmhr) ResponseRESTError(status int, err error, req *http.Request) (server.Response, error) {
	commonErr := errata.CommonError(err)

	key := commonErr.Key()
	data := common.Map{server.ErrorKey: key}

	if status == 0 || status == http.StatusOK {
		if key == errata.NoCredsKey || key == errata.InvalidCredsKey {
			status = http.StatusUnauthorized
		} else if key == errata.OverdueRightsErr || key == errata.NoUserKey || key == errata.NoRightsKey {
			status = http.StatusForbidden
		} else if status == 0 || status == http.StatusOK {
			status = http.StatusInternalServerError

		} else {
			status = http.StatusInternalServerError
		}
	}

	if !strlib.In(s.secretENVsToLower, strings.ToLower(os.Getenv("ENV"))) {
		data["details"] = commonErr.Error()
	}

	if req != nil {
		err = fmt.Errorf("ERROR on %s %s, got: %s", req.Method, req.URL, commonErr.Error())
		// TODO: add body[:2048] for debugging
	} else {
		err = commonErr
	}

	jsonBytes, errJSON := json.Marshal(data)
	if errJSON != nil {
		l.Errorf("ERROR marshalling error data (%#v): %s", data, errJSON)
	}
	return server.Response{Status: status, Data: jsonBytes}, err
}

func (s *serverHTTPJschmhr) ResponseRESTOk(status int, data interface{}) (server.Response, error) {
	if status == 0 {
		status = http.StatusOK
	}

	if data == nil {
		return server.Response{Status: status}, nil
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return server.Response{Status: http.StatusInternalServerError}, errata.Wrapf(err, "can't marshal json (%#v)", data)
	}

	return server.Response{Status: status, Data: jsonBytes}, nil
}

func (s *serverHTTPJschmhr) HandleEndpoint(key, serverPath string, endpoint server_http.Endpoint) error {

	method := strings.ToUpper(endpoint.Method)
	path := endpoint.PathTemplate(serverPath)

	if endpoint.WorkerHTTP == nil {
		return errata.New(method + ": " + path + "\t!!! NULL workerHTTP ISN'T DISPATCHED !!!")
	}

	s.HandleOptions(key, path)

	handler := func(w http.ResponseWriter, r *http.Request, paramsHR httprouter.Params) {
		options, err := s.onRequest.Options(r)
		if err != nil {
			l.Error(err)
		}

		var params server_http.Params
		if len(paramsHR) > 0 {
			params = server_http.Params{}
			for _, p := range paramsHR {
				params[p.Key] = p.Value
			}
		}

		w.Header().Set("Access-Control-Allow-Origin", server_http.CORSAllowOrigin)
		w.Header().Set("Access-Control-Allow-Headers", server_http.CORSAllowHeaders)
		w.Header().Set("Access-Control-Allow-Methods", server_http.CORSAllowMethods)
		w.Header().Set("Access-Control-Allow-Credentials", server_http.CORSAllowCredentials)

		responseData, err := endpoint.WorkerHTTP(s, r, params, options)
		if err != nil {
			l.Error(err)
		}

		if responseData.MIMEType != "" {
			w.Header().Set("Content-Type", responseData.MIMEType)
		}
		w.Header().Set("Content-Length", strconv.Itoa(len(responseData.Data)))
		if responseData.FileName != "" {
			w.Header().Set("Content-Disposition", "attachment; filename="+responseData.FileName)
		}

		if responseData.Status > 0 {
			w.WriteHeader(responseData.Status)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		if _, err := w.Write(responseData.Data); err != nil {
			l.Error("can't write response", err)
		}
	}

	l.Infof("%-10s: %s %s", key, method, path)
	switch method {
	case "GET":
		s.httpServeMux.GET(path, handler)
	case "POST":
		s.httpServeMux.POST(path, handler)
	case "PUT":
		s.httpServeMux.PUT(path, handler)
	case "DELETE":
		s.httpServeMux.DELETE(path, handler)
	default:
		l.Error(method, " isn't supported!")
	}

	return nil
}

func (s *serverHTTPJschmhr) HandleOptions(key, serverPath string) {
	//if strlib.In(s.handledOptions, serverPath) {
	//	//l.Infof("- %#v", s.handledOptions)
	//	return
	//}

	s.httpServeMux.OPTIONS(serverPath, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		l.Infof("%-10s: OPTIONS %s", key, serverPath)
		w.Header().Set("Access-Control-Allow-Origin", server_http.CORSAllowOrigin)
		w.Header().Set("Access-Control-Allow-Headers", server_http.CORSAllowHeaders)
		w.Header().Set("Access-Control-Allow-Methods", server_http.CORSAllowMethods)
		w.Header().Set("Access-Control-Allow-Credentials", server_http.CORSAllowCredentials)
	})

	//s.handledOptions = append(s.handledOptions, serverPath)
}

var reHTMLExt = regexp.MustCompile(`\.html?$`)

func (s *serverHTTPJschmhr) HandleFiles(key, serverPath string, staticPath server_http.StaticPath) error {
	l.Infof("%-10s: FILES %s <-- %s", key, serverPath, staticPath.LocalPath)

	// TODO: check localPath

	if staticPath.MIMEType == nil {
		// TODO!!! CORS

		s.httpServeMux.ServeFiles(serverPath, http.Dir(staticPath.LocalPath))
		return nil
	}

	s.HandleOptions(key, serverPath)

	//fileServer := http.FileServer(http.Dir(localPath))
	s.httpServeMux.GET(serverPath, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", server_http.CORSAllowOrigin)
		w.Header().Set("Access-Control-Allow-Headers", server_http.CORSAllowHeaders)
		w.Header().Set("Access-Control-Allow-Methods", server_http.CORSAllowMethods)
		w.Header().Set("Access-Control-Allow-Credentials", server_http.CORSAllowCredentials)

		if staticPath.MIMEType != nil && *staticPath.MIMEType != "" {
			w.Header().Set("Content-Type", *staticPath.MIMEType)
		}

		OpenFile, err := os.Open(staticPath.LocalPath + "/" + p.ByName("filepath"))
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

// mimeTypeToSet, err = inspector.MIME(localPath+"/"+r.ExportID.PathWithParams, nil)
// if err != nil {
//	l.ErrStr("can't read MIMEType for file: ", localPath+"/"+r.ExportID.PathWithParams, err)
// }

//func (s *serverHTTPJschmhr) HandleGetString(serverRoute, str string, mimeType *string) {
//	s.handleFunc("GET", serverRoute, func(w http.ResponseWriter, r *http.Request, params httprouter.Content) {
//		if mimeType != nil {
//			// "application/javascript"
//			w.Header().Set("Content-Type", *mimeType)
//		}
//		w.Write([]byte(str))
//	})
//}
