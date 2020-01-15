package v1_auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/common/server"
	"github.com/pavlo67/workshop/common/server/server_http"
)

var _ = server_http.InitEndpoint(&endpoints, "POST", filelib.RelativePath(filelib.CurrentFile(true), pathBase, prefix, l),
	nil, workerAuth, "")
var _ server_http.WorkerHTTP = workerAuth

func workerAuth(_ *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {

	credsJSON, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return server.ResponseRESTError(http.StatusBadRequest, errors.Wrap(err, "can't read body"))
	}

	r.L.Infof("%s", credsJSON)

	var toAuth auth.Creds
	err = json.Unmarshal(credsJSON, &toAuth)
	if err != nil {
		return server.ResponseRESTError(http.StatusBadRequest, errors.Wrapf(err, "can't unmarshal body: %s", credsJSON))
	}

	user, errs := auth.GetUser(toAuth, r.AuthOps, nil)
	if len(errs) > 0 {
		return server.ResponseRESTError(http.StatusForbidden, errs.Err())
	}
	if user == nil {
		return server.ResponseRESTError(http.StatusForbidden, errors.New("no user authorized"))
	}

	toAddModified, err := r.AuthOpToSetToken.SetCreds(*user, auth.Creds{}) // TODO!!! add custom toAddModified
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

	r.L.Info(user)

	return server.ResponseRESTOk(map[string]interface{}{"user": user})
}
