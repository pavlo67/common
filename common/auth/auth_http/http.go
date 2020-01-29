package auth_http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/server"
	"github.com/pavlo67/workshop/common/server/server_http"
)

const onNew = "on auth_http.New(): "

func New() (authorizeHandler, setCredsHandler, getCredsHandler *server_http.Endpoint, err error) {
	return &authorizeEndpoint, &setCredsEndpoint, &getCredsEndpoint, nil
}

var setCredsEndpoint = server_http.Endpoint{
	Method: "POST",
	WorkerHTTP: func(user *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {
		return server.ResponseRESTError(http.StatusBadRequest, errors.Errorf("the .SetCreds() method is temporary unavailable"))

		//if authOpToInit == nil {
		//	return server.ResponseRESTOk(map[string]interface{}{})
		//}
		//
		//credsJSON, err := ioutil.ReadAll(req.Body)
		//if err != nil {
		//	return server.ResponseRESTError(http.StatusBadRequest, errors.Wrap(err, "can't read body"))
		//}
		//
		//l.Infof("%s", credsJSON)
		//
		//var toInit auth.Creds
		//err = json.Unmarshal(credsJSON, &toInit)
		//if err != nil {
		//	return server.ResponseRESTError(http.StatusBadRequest, errors.Wrapf(err, "can't unmarshal body: %s", credsJSON))
		//}
		//
		//toInit[auth.CredsIP] = req.RemoteAddr
		//
		//var userKey identity.Key
		//if user != nil {
		//	userKey = user.Key
		//}
		//
		//creds, err := authOpToInit.SetCreds(userKey, toInit)
		//if err != nil {
		//	return server.ResponseRESTError(http.StatusForbidden, err)
		//}
		//if creds == nil {
		//	return server.ResponseRESTError(http.StatusForbidden, errors.New("no creds to init auth session"))
		//}
		//
		//return server.ResponseRESTOk(map[string]interface{}{"creds": *creds})
	},
}

var getCredsEndpoint = server_http.Endpoint{
	Method: "POST",
	WorkerHTTP: func(user *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {
		return server.ResponseRESTOk(map[string]interface{}{"user": user})
	},
}

var authorizeEndpoint = server_http.Endpoint{
	Method: "POST",
	WorkerHTTP: func(_ *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {
		credsJSON, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return server.ResponseRESTError(http.StatusBadRequest, errors.Wrap(err, "can't read body"))
		}

		l.Infof("%s", credsJSON)

		var toAuth auth.Creds
		err = json.Unmarshal(credsJSON, &toAuth)
		if err != nil {
			return server.ResponseRESTError(http.StatusBadRequest, errors.Wrapf(err, "can't unmarshal body: %s", credsJSON))
		}

		toAuth[auth.CredsIP] = req.RemoteAddr

		user, errs := auth.GetUser(toAuth, authOps, nil)
		if len(errs) > 0 {
			return server.ResponseRESTError(http.StatusForbidden, errs.Err())
		}
		if user == nil {
			return server.ResponseRESTError(http.StatusForbidden, errors.New("no user authorized"))
		}

		toAddModified, err := authOpToSetToken.SetCreds(
			user.Key,
			auth.Creds{
				auth.CredsNickname: user.Creds[auth.CredsNickname],
				auth.CredsEmail:    user.Creds[auth.CredsEmail],
				auth.CredsToSet:    string(auth.CredsJWT),
			},
		)
		// TODO!!! add custom toAddModified

		if err != nil {
			return server.ResponseRESTError(http.StatusInternalServerError, errors.Wrap(err, "can't create JWT"))
		}

		if toAddModified != nil {
			if user.Creds == nil {
				user.Creds = map[auth.CredsType]string{}
			}

			for t, c := range *toAddModified {
				user.Creds[t] = c
			}
		}

		l.Infof("user: %#v", user)

		return server.ResponseRESTOk(map[string]interface{}{"user": user})
	},
}
