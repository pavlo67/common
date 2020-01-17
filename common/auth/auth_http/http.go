package auth_http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pavlo67/workshop/common/libraries/encrlib"

	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/server"
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/server/server_http"
)

const onNew = "on auth_http.New(): "

func New() (authHandler, initAuthSessionHandler *server_http.Endpoint, err error) {
	return &authEndpoint, &initAuthSessionEndpoint, nil
}

var initAuthSessionEndpoint = server_http.Endpoint{
	Method: "POST",
	WorkerHTTP: func(_ *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {
		if authOpToInit == nil {
			return server.ResponseRESTOk(map[string]interface{}{})
		}

		toInit := auth.Creds{
			Cryptype: encrlib.NoCrypt,
			Values: auth.Values{
				auth.CredsIP: req.RemoteAddr,
			},
		}

		creds, err := authOpToInit.InitAuth(toInit)
		if err != nil {
			return server.ResponseRESTError(http.StatusForbidden, err)
		}
		if creds == nil {
			return server.ResponseRESTError(http.StatusForbidden, errors.New("no creds to init auth session"))
		}

		return server.ResponseRESTOk(map[string]interface{}{"creds": *creds})
	},
}

var authEndpoint = server_http.Endpoint{
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

		toAuth.Values[auth.CredsIP] = req.RemoteAddr

		user, errs := auth.GetUser(toAuth, authOps, nil)
		if len(errs) > 0 {
			return server.ResponseRESTError(http.StatusForbidden, errs.Err())
		}
		if user == nil {
			return server.ResponseRESTError(http.StatusForbidden, errors.New("no user authorized"))
		}

		_, toAddModified, err := authOpToSetToken.SetCreds(user, auth.Creds{}) // TODO!!! add custom toAddModified
		if err != nil {
			return server.ResponseRESTError(http.StatusInternalServerError, errors.Wrap(err, "can't create JWT"))
		}

		if toAddModified != nil {
			if user.Creds.Values == nil {
				user.Creds.Values = map[auth.CredsType]string{}
			}

			for t, c := range toAddModified.Values {
				user.Creds.Values[t] = c
			}
		}

		l.Infof("user: %#v", user)

		return server.ResponseRESTOk(map[string]interface{}{"user": user})
	},
}
