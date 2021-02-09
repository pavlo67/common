package auth_server_http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/errata"
	"github.com/pavlo67/common/common/server"
	"github.com/pavlo67/common/common/server/server_http"
)

var Endpoints = server_http.Endpoints{
	auth.IntefaceKeyAuthenticate: authenticateEndpoint,
	auth.IntefaceKeySetCreds:     setCredsEndpoint,
}

//var bodyParams = json.RawMessage(`{
//   "in": "body",
//	"name": "credentials",
//	"description": "user's email/login & password'",
//	"schema": {
//		"type": "object",
//		"required":"password",
//		"properties": {
//			"email":    {"type": "string"},
//			"nickname": {"type": "string"},
//			"password": {"type": "string"}
//          ...
//		}
//	}
//
//}`)

var authenticateEndpoint = server_http.Endpoint{
	Method: "POST",
	//BodyParams: bodyParams,
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, _ server_http.Params, _ *crud.Options) (server.Response, error) {

		credsJSON, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return serverOp.ResponseRESTError(http.StatusBadRequest, errata.KeyableError(errata.WrongBodyKey, common.Map{"error": errors.Wrap(err, "can't read body")}), req)
		}

		var toAuth auth.Creds
		if err = json.Unmarshal(credsJSON, &toAuth); err != nil {
			return serverOp.ResponseRESTError(http.StatusBadRequest, errata.KeyableError(errata.WrongJSONKey, common.Map{"error": errors.Wrapf(err, "can't unmarshal body: %s", credsJSON)}), req)
		}
		toAuth[auth.CredsIP] = req.RemoteAddr

		identity, err := authOp.Authenticate(toAuth)
		if err != nil {
			return serverOp.ResponseRESTError(0, err, req)
		}

		return serverOp.ResponseRESTOk(http.StatusOK, identity)
	},
}

var setCredsEndpoint = server_http.Endpoint{
	Method: "POST",
	//BodyParams: bodyParams,
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, _ server_http.Params, options *crud.Options) (server.Response, error) {

		credsJSON, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return serverOp.ResponseRESTError(http.StatusBadRequest, errata.KeyableError(errata.WrongBodyKey, common.Map{"error": errors.Wrap(err, "can't read body")}), req)
		}

		var toSet auth.Creds
		if err = json.Unmarshal(credsJSON, &toSet); err != nil {
			return serverOp.ResponseRESTError(http.StatusBadRequest, errata.KeyableError(errata.WrongJSONKey, common.Map{"error": errors.Wrapf(err, "can't unmarshal body: %s", credsJSON)}), req)
		}
		toSet[auth.CredsIP] = req.RemoteAddr

		var authID auth.ID
		if options != nil && options.Identity != nil {
			authID = options.Identity.ID
		}

		creds, err := authOp.SetCreds(authID, toSet)
		if err != nil {
			return serverOp.ResponseRESTError(0, err, req)
		}

		return serverOp.ResponseRESTOk(http.StatusOK, creds)
	},
}
