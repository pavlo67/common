package auth_server_http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/server"
	"github.com/pavlo67/common/common/server/server_http"

	"github.com/pavlo67/common/common/errors"
)

var bodyParams = json.RawMessage(`{
   "in": "body",
	"name": "credentials",
	"description": "user's email/login & password'",
	"schema": {
		"type": "object",
		"required":"password",
		"properties": {
			"email":    {"type": "string"},
			"nickname": {"type": "string"},
			"password": {"type": "string"}
		}
	}

}`)

var authEndpoint = server_http.Endpoint{
	Method:     "POST",
	BodyParams: bodyParams,
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, _ server_http.Params, _ *crud.Options) (server.Response, error) {

		credsJSON, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return serverOp.ResponseRESTError(http.StatusBadRequest, errors.KeyableError(errors.Wrap(err, "can't read body"), errors.WrongBodyErr, nil), req)
		}

		var toAuth auth.Creds
		if err = json.Unmarshal(credsJSON, &toAuth); err != nil {
			return serverOp.ResponseRESTError(http.StatusBadRequest, errors.KeyableError(errors.Wrapf(err, "can't unmarshal body: %s", credsJSON), errors.WrongJSONErr, nil), req)
		}

		toAuth[auth.CredsIP] = req.RemoteAddr

		identity, errorKey, errs := auth.GetIdentity(toAuth, authOps, false, nil)
		if identity != nil {
			result := common.Map{} // "user": persons.Item{Identity: *identity}
			if errorKey != "" {
				result[server.ErrorKey] = errorKey
			}
			if len(errs) > 0 {
				result["errors"] = errs.Err()
			}
			return serverOp.ResponseRESTOk(http.StatusOK, result)
		}

		if errorKey == "" {
			errorKey = errors.NoCredsErr
		}

		if len(errs) > 0 {
			return serverOp.ResponseRESTError(0, errors.KeyableError(errs.Err(), errorKey, nil), req)
		}

		return serverOp.ResponseRESTError(0, errors.KeyableError(errors.New("no identity authorized"), errorKey, nil), req)
	},
}

//if identity.JWT == "" && authOpToSetToken != nil {
//	toAddModified, err := authOpToSetToken.SetCreds(
//		identity.ID,
//		auth.Creds{
//			auth.CredsNickname: identity.Nickname,
//			auth.CredsRoles:    identity.Roles,
//			auth.CredsToSet:    auth.CredsJWT,
//		},
//	)
//	if err != nil || toAddModified == nil {
//		return serverOp.ResponseRESTError(identity, 0, fmt.Errorf("can't create JWT. got %s / %#v", err, toAddModified), req)
//	}
//	identity.JWT, _ = toAddModified.String(auth.CredsJWT)
//	// TODO!!! add CompanyID, OperatorAccountID
//}
