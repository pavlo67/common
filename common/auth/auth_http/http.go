package auth_http

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/errata"
	"github.com/pavlo67/common/common/server/server_http"
)

var _ auth.Operator = &authHTTP{}

type authHTTP struct {
	serverConfig server_http.Config

	logfile string
}

const onNew = "on authHTTP.New()"

func New(serverConfig server_http.Config) (auth.Operator, error) {
	authOp := authHTTP{
		serverConfig: serverConfig,
	}

	return &authOp, nil
}

func (authOp *authHTTP) SetCreds(userID auth.ID, toSet auth.Creds) (*auth.Creds, error) {
	return nil, errata.NotImplemented

}

// Authenticate can require to do .SetCreds first and to use some session-generated creds

const onAuthenticate = "on authHTTP.Authenticate()"

func (authOp *authHTTP) Authenticate(toAuth auth.Creds) (*auth.Identity, error) {
	ep := authOp.serverConfig.EndpointsSettled[auth.IntefaceKeyAuthenticateHandler]
	serverURL := authOp.serverConfig.Host + authOp.serverConfig.Port + ep.Path

	requestBody, err := json.Marshal(toAuth)
	if err != nil {
		return nil, errors.Wrapf(err, onAuthenticate+": can't marshal toAuth(%#v)", toAuth)
	}

	var identity auth.Identity
	if err := server_http.Request(serverURL, ep, requestBody, &identity, nil, authOp.logfile); err != nil {
		return nil, err
	}

	return &identity, nil
}
