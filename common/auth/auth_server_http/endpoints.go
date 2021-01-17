package auth_server_http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/server"
	"github.com/pavlo67/workshop/common/server/server_http"

	"github.com/pkg/errors"
)

var authEndpoint = server_http.Endpoint{
	Method: "POST",
	BodyParams: json.RawMessage(`{
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

}`),

	WorkerHTTP: func(serverOp server_http.Operator, identity *auth.Identity, _ server_http.Params, req *http.Request) (server.Response, error) {
		credsJSON, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return serverOp.ResponseRESTError(identity, http.StatusBadRequest, common.KeyableError(common.WrongBodyErr, nil, errors.Wrap(err, "can't read body")))
		}

		l.Debugf("%s", credsJSON)

		var toAuth auth.Creds
		err = json.Unmarshal(credsJSON, &toAuth)
		if err != nil {
			return serverOp.ResponseRESTError(identity, http.StatusBadRequest, common.KeyableError(common.WrongJSONErr, nil, errors.Wrapf(err, "can't unmarshal body: %s", credsJSON)))
		}

		toAuth[auth.CredsIP] = req.RemoteAddr

		//for _, authOp := range authOps {
		//	l.Infof("%#v", authOp)
		//}
		//l.Infof("authOps length = %d", len(authOps))

		identity, errorKey, errs := auth.GetIdentity(toAuth, authOps, false, nil)
		if identity != nil {
			result := common.Map{} // "user": persons.Item{Identity: *identity}
			if errorKey != "" {
				result[server.ErrorKey] = errorKey
			}
			if len(errs) > 0 {
				result["errors"] = errs.Err()
			}
			return serverOp.ResponseRESTOk(identity, result)
		}

		if errorKey == "" {
			errorKey = common.NoCredsErr
		}

		if len(errs) > 0 {
			return serverOp.ResponseRESTError(identity, 0, common.KeyableError(errorKey, nil, errs.Err()))
		}

		return serverOp.ResponseRESTError(identity, 0, common.KeyableError(errorKey, nil, errors.New("no identity authorized")))
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
//		return serverOp.ResponseRESTError(identity, 0, errors.Errorf("can't create JWT. got %s / %#v", err, toAddModified), req)
//	}
//	identity.JWT, _ = toAddModified.String(auth.CredsJWT)
//	// TODO!!! add CompanyID, OperatorAccountID
//}
