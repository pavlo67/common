package auth_http

import (
	"encoding/json"

	"github.com/pavlo67/common/common/crud"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/server/server_http"
)

var _ auth.Operator = &authHTTP{}

type authHTTP struct {
	serverConfig server_http.Config
}

const onNew = "on authHTTP.New()"

func New(serverConfig server_http.Config) (auth.Operator, error) {
	authOp := authHTTP{
		serverConfig: serverConfig,
	}

	return &authOp, nil
}

func (authOp *authHTTP) SetCreds(authID auth.ID, toSet auth.Creds) (*auth.Creds, error) {
	ep := authOp.serverConfig.EndpointsSettled[auth.IntefaceKeySetCreds]
	serverURL := authOp.serverConfig.Host + authOp.serverConfig.Port + ep.Path

	requestBody, err := json.Marshal(toSet)
	if err != nil {
		return nil, errors.Wrapf(err, onAuthenticate+": can't marshal toSet(%#v)", toSet)
	}

	var creds *auth.Creds
	if err := server_http.Request(serverURL, ep, requestBody, creds, &crud.Options{Identity: &auth.Identity{ID: authID}}, l); err != nil {
		return nil, err
	}

	return creds, nil
}

// Authenticate can require to call .SetCreds() first and to use some session-generated creds
const onAuthenticate = "on authHTTP.Authenticate()"

func (authOp *authHTTP) Authenticate(toAuth auth.Creds) (*auth.Identity, error) {
	ep := authOp.serverConfig.EndpointsSettled[auth.IntefaceKeyAuthenticate]
	serverURL := authOp.serverConfig.Host + authOp.serverConfig.Port + ep.Path

	requestBody, err := json.Marshal(toAuth)
	if err != nil {
		return nil, errors.Wrapf(err, onAuthenticate+": can't marshal toAuth(%#v)", toAuth)
	}

	var identity *auth.Identity
	if err := server_http.Request(serverURL, ep, requestBody, &identity, nil, l); err != nil {
		return nil, err
	}

	return identity, nil
}
