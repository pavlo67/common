package v1_auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/apps/confidence/confidence_routes"
	"github.com/pavlo67/workshop/common/server"
	"github.com/pavlo67/workshop/common/server/server_http"
	"github.com/pavlo67/workshop/components/auth"
	"github.com/pavlo67/workshop/libraries/filelib"
	"github.com/pavlo67/workshop/libraries/strlib"
)

var _ = server_http.InitEndpoint(&confidence_routes.Endpoints, "POST", filelib.RelativePath(filelib.CurrentFile(true), confidence_routes.BasePath, confidence_routes.Prefix),
	nil, workerRegister, "")
var _ server_http.WorkerHTTP = workerRegister

func workerRegister(_ *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {

	var testModePath string

	if confidence_routes.Cfg != nil && strlib.In(confidence_routes.Cfg.ServerHTTP.Testers, req.RemoteAddr) {
		testModePath = req.Header.Get("Test-Mode-Path")
	}

	if testModePath != "" {

	}

	credsJSON, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return server.ResponseRESTError(http.StatusBadRequest, errors.Wrap(err, "can't read body"))
	}

	// log.Printf("%s", credsJSON)

	var toSet auth.Creds
	err = json.Unmarshal(credsJSON, &toSet)
	if err != nil {
		return server.ResponseRESTError(http.StatusBadRequest, errors.Wrapf(err, "can't unmarshal body: %s", credsJSON))
	}

	//
	//user, errs := auth.GetUser(toSet, confidence_routes.AuthOps, nil)
	//if len(errs) > 0 {
	//	return server.ResponseRESTError(http.StatusForbidden, errs.Err())
	//}
	//if user == nil {
	//	return server.ResponseRESTError(http.StatusForbidden, errors.New("no user authorized"))
	//}
	//
	//toAddModified, err := confidence_routes.AuthOpToSetToken.SetCreds(*user, auth.Creds{}) // TODO!!! add custom toAddModified
	//if err != nil {
	//	return server.ResponseRESTError(http.StatusInternalServerError, errors.Wrap(err, "can't create JWT"))
	//}
	//
	//if toAddModified != nil {
	//	for t, c := range toAddModified.Values {
	//		user.Creds.Values[t] = c
	//	}
	//}
	//
	return server.ResponseRESTOk(map[string]interface{}{"user": nil})
}
