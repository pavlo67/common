package confidence_v1

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/basis/auth"
	"github.com/pavlo67/workshop/basis/server"
	"github.com/pavlo67/workshop/basis/server/server_http"

	"github.com/pavlo67/workshop/apps/rest/confidence/confidence_routes"
	"github.com/pavlo67/workshop/basis/common/filelib"
)

var ToInit = server_http.InitEndpoint(&confidence_routes.Endpoints, "POST", filelib.RelativePath(confidence_routes.Prefix, filelib.CurrentFile(true)), nil, workerAuth, "")
var _ server_http.WorkerHTTP = workerAuth

func workerAuth(_ *auth.User, _ server_http.Params, req *http.Request) (server.Response, error) {

	credsJSON, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return server.ResponseRESTError(http.StatusBadRequest, errors.Wrap(err, "can't read body"))
	}

	var toAuth []auth.Creds
	err = json.Unmarshal(credsJSON, &toAuth)
	if err != nil {
		return server.ResponseRESTError(http.StatusBadRequest, errors.Wrapf(err, "can't unmarshal body: %s", credsJSON))
	}

	user, creds, err := confidence_routes.AuthOp.Authorize(toAuth...)
	if err != nil {
		return server.ResponseRESTError(http.StatusForbidden, err)
	}

	return server.ResponseRESTOk(map[string]interface{}{"user": user, "creds": creds})
}
