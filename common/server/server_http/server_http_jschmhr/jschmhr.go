package server_http_jschmhr

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/libraries/strlib"
	"github.com/pavlo67/workshop/common/server"
	"github.com/pavlo67/workshop/common/server/server_http"
)

var _ server_http.Operator = &serverHTTPJschmhr{}

type serverHTTPJschmhr struct {
	httpServer   *http.Server
	httpServeMux *httprouter.Router

	port        int
	certFileTLS string
	keyFileTLS  string
	authOps     []auth.Operator

	handledOptions []string
}

func New(port int, certFileTLS, keyFileTLS string, authOps []auth.Operator, noEventsOp bool) (server_http.Operator, error) {
	if port <= 0 {
		return nil, errors.Errorf("on server_http_jschmhr.New(): wrong port = %d", port)
	}

	if !noEventsOp {
		//if eventsOpSystem == nil {
		//	return nil, errors.New("on server_http_jschmhr.New(): no events.OperatorSystem")
		//} else if eventsOp == nil {
		//	return nil, errors.New("on server_http_jschmhr.New(): no events.Operator")
		//}
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

		port: port,

		certFileTLS: certFileTLS,
		keyFileTLS:  keyFileTLS,

		authOps: authOps,
		//eventsOp:       eventsOp,
		//eventsOpSystem: eventsOpSystem,
	}, nil
}

//const onNotifyByREST = "on serverHTTPJschmhr.NotifyByREST()"
//
//func (s *serverHTTPJschmhr) NotifyByREST(identity *auth.Identity, responseBytes []byte) []byte {
//	if s.eventsOpSystem == nil {
//		return responseBytes
//	}
//
//	crudOptions := crud.Options{Identity: identity}
//
//	eventsToNotify, err := s.eventsOp.ListEventsToNotify(false, &crudOptions)
//	if err != nil {
//		l.Error(onNotifyByREST + ": " + err.Error())
//	}
//
//	if len(eventsToNotify) < 1 {
//		return responseBytes
//	}
//
//	// l.Infof("EVENTS TO NOTIFY: %#v, %#v", crudOptions, eventsToNotify)
//
//	responseBytesChanged, err := jsonlib.AddKeyValue(responseBytes, "Alerts", eventsToNotify)
//	if err != nil {
//		l.Errorf(onNotifyByREST+": %s (%s)", err, responseBytes)
//		responseBytesChanged = responseBytes
//	} else {
//		// l.Infof(onNotifyByREST+": %s", responseBytesChanged)
//	}
//
//	for _, eventToNotify := range eventsToNotify {
//		if err2 := s.eventsOpSystem.SaveNotification(eventToNotify.ID, alerts.PopUp, err); err2 != nil {
//			l.Errorf(onNotifyByREST+": %s", err2)
//		}
//	}
//
//	return responseBytesChanged
//}

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

func (s *serverHTTPJschmhr) ResponseRESTError(identity *auth.Identity, status int, keyableErr common.Error, req ...*http.Request) (server.Response, error) {
	if keyableErr == nil {
		keyableErr = common.KeyableError("", nil, errors.Errorf("unknown error with status %d", status))
	}
	key := keyableErr.Key()

	data := common.Map{server.ErrorKey: key}
	if key == common.NoCredsErr || key == common.InvalidCredsErr {
		status = http.StatusUnauthorized
	} else if key == common.OverdueRightsErr || key == common.NoUserErr || key == common.NoRightsErr {
		status = http.StatusForbidden
	} else if status == 0 || status == http.StatusOK {
		status = http.StatusInternalServerError
	}

	if os.Getenv("ENV") != "production" {
		data["details"] = keyableErr.Error()
	}

	var err error
	if len(req) > 0 && req[0] != nil {
		err = errors.Errorf("ERROR on %s %s, got: %s", req[0].Method, req[0].URL, keyableErr.Err())
		// TODO: add body[:2048] for debugging
	} else {
		err = keyableErr.Err()
	}

	jsonBytes, _ := json.Marshal(data)
	//if identity != nil {
	//	jsonBytes = s.NotifyByREST(identity, jsonBytes)
	//}
	return server.Response{Status: status, Data: jsonBytes}, err
}

func (s *serverHTTPJschmhr) ResponseRESTOk(identity *auth.Identity, data interface{}) (server.Response, error) {
	if data == nil {
		return server.Response{Status: http.StatusOK}, nil
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return server.Response{Status: http.StatusInternalServerError}, errors.Wrapf(err, "can't marshal pbxm (%#v)", data)
	}
	//if identity != nil {
	//	jsonBytes = s.NotifyByREST(identity, jsonBytes)
	//}

	return server.Response{Status: http.StatusOK, Data: jsonBytes}, nil
}

func (s *serverHTTPJschmhr) HandleOptions(key, serverPath string) {
	if strlib.In(s.handledOptions, serverPath) {
		//l.Infof("- %#v", s.handledOptions)
		return
	}

	s.httpServeMux.OPTIONS(serverPath, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		l.Infof("%-10s: OPTIONS %s", key, serverPath)
		w.Header().Set("Access-Control-Allow-Origin", server_http.CORSAllowOrigin)
		w.Header().Set("Access-Control-Allow-Headers", server_http.CORSAllowHeaders)
		w.Header().Set("Access-Control-Allow-Methods", server_http.CORSAllowMethods)
		w.Header().Set("Access-Control-Allow-Credentials", server_http.CORSAllowCredentials)
	})

	s.handledOptions = append(s.handledOptions, serverPath)
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

func (s *serverHTTPJschmhr) HandleEndpoint(key, serverPath string, endpoint server_http.Endpoint) error {

	method := strings.ToUpper(endpoint.Method)
	path := endpoint.PathTemplate(serverPath)

	if endpoint.WorkerHTTP == nil {
		return errors.New(method + ": " + path + "\t!!! NULL workerHTTP ISN'T DISPATCHED !!!")
	}

	s.HandleOptions(key, path)

	handler := func(w http.ResponseWriter, r *http.Request, paramsHR httprouter.Params) {
		identity, _, err := server_http.IdentityWithRequest(r, s.authOps)
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

		responseData, err := endpoint.WorkerHTTP(s, identity, params, r)
		if err != nil {
			l.Error(err)

			//http.Critical(w, string(responseData.AccountData), responseData.Status)
			//return
		}

		// l.Infof("responseData: %#v", responseData)

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

// mimeTypeToSet, err = inspector.MIME(localPath+"/"+r.ExportID.PathWithParams, nil)
// if err != nil {
//	l.ErrStr("can't read MIMEType for file: ", localPath+"/"+r.ExportID.PathWithParams, err)
// }

//func (s *serverHTTPJschmhr) HandleGetString(serverRoute, str string, mimeType *string) {
//	s.handleFunc("GET", serverRoute, func(w http.ResponseWriter, r *http.Request, params httprouter.Content) {
//		if mimeType != nil {
//			// "application/javascript"
//			w.Header().Set("Content-TypeKey", *mimeType)
//		}
//		w.Write([]byte(str))
//	})
//}
