package auth_server_http

import (
	"encoding/json"
	"net/http"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/errata"
	"github.com/pavlo67/common/common/server"
	"github.com/pavlo67/common/common/server/server_http"
)

var Endpoints = server_http.Endpoints{
	auth.IntefaceKeyAuthenticateHandler: authPassEndpoint,
}

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

var authPassEndpoint = server_http.Endpoint{
	Method:     "POST",
	BodyParams: bodyParams,
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, _ server_http.Params, _ *crud.Options) (server.Response, error) {

		//credsJSON, err := ioutil.ReadAll(req.Body)
		//if err != nil {
		//	return serverOp.ResponseRESTError(http.StatusBadRequest, errata.KeyableError(errata.WrongBodyKey, common.Map{"error": errors.Wrap(err, "can't read body")}), req)
		//}
		//
		//var toAuth auth.Creds
		//if err = json.Unmarshal(credsJSON, &toAuth); err != nil {
		//	return serverOp.ResponseRESTError(http.StatusBadRequest, errata.KeyableError(errata.WrongJSONKey, common.Map{"error": errors.Wrapf(err, "can't unmarshal body: %s", credsJSON)}), req)
		//}
		//
		//toAuth[auth.CredsIP] = req.RemoteAddr
		//
		//identity, errorKey, errs := auth.GetIdentity(toAuth, authOps, false, nil)
		//if identity != nil {
		//	result := common.Map{} // "user": persons.Item{Identity: *identity}
		//	if errorKey != "" {
		//		result[server.ErrorKey] = errorKey
		//	}
		//	if len(errs) > 0 {
		//		result["errors"] = errs.Err()
		//	}
		//	return serverOp.ResponseRESTOk(http.StatusOK, result)
		//}
		//
		//if errorKey == "" {
		//	errorKey = errata.NoCredsKey
		//}
		//
		//if len(errs) > 0 {
		//	return serverOp.ResponseRESTError(0, errata.KeyableError(errorKey, common.Map{"error": errs.Err()}), req)
		//}
		//
		//return serverOp.ResponseRESTError(0, errata.KeyableError(errorKey, common.Map{"error": "no identity authorized"}), req)

		return serverOp.ResponseRESTError(0, errata.CommonError(errata.NotImplementedKey), req)
	},
}

//if identity.Token == "" && authOpToSetToken != nil {
//	toAddModified, err := authOpToSetToken.SetCreds(
//		identity.ID,
//		auth.Creds{
//			auth.CredsNickname: identity.Nickname,
//			auth.CredsRoles:    identity.Roles,
//			auth.CredsToSet:    auth.CredsJWT,
//		},
//	)
//	if err != nil || toAddModified == nil {
//		return serverOp.ResponseRESTError(identity, 0, fmt.Errorf("can't create Token. got %s / %#v", err, toAddModified), req)
//	}
//	identity.Token, _ = toAddModified.String(auth.CredsJWT)
//	// TODO!!! add CompanyID, OperatorAccountID
//}
