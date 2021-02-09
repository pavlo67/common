package auth_server_http

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/errata"
	"github.com/pavlo67/common/common/server/server_http"
)

var _ server_http.OnRequestMiddleware = &onRequestMiddleware{}

func OnRequestMiddleware(authJWTOp auth.Operator) (server_http.OnRequestMiddleware, error) {
	if authJWTOp == nil {
		return nil, errors.New("no authJWTOp")
	}

	return &onRequestMiddleware{
		authJWTOp: authJWTOp,
	}, nil
}

type onRequestMiddleware struct {
	authJWTOp auth.Operator
}

var reBearer = regexp.MustCompile(`^\s*Bearer(\s|%[fF]20)*`)

const onOptions = "on onRequestMiddleware.Options()"

func (orm *onRequestMiddleware) Options(r *http.Request) (*crud.Options, error) {
	//if r == nil {
	//	return nil, errors.New("no server_http.Request in RequestOptions(...)")
	//}

	var identity *auth.Identity
	if tokenJWT := r.Header.Get("Authorization"); tokenJWT != "" {
		tokenJWT = reBearer.ReplaceAllString(tokenJWT, "")
		var err error
		if identity, err = orm.authJWTOp.Authenticate(auth.Creds{auth.CredsJWT: tokenJWT}); err != nil {
			return nil, errata.CommonError(err, onOptions)
		}
	}

	if identity == nil {
		return nil, nil
	}

	return &crud.Options{Identity: identity}, nil
}

//// SIGNATURE CHECK
//signature := r.Header.Get("Signature")
//if signature != "" && r.URL != nil {
//	publicKeyAddress := r.Header.Get("Public-Key-Address")
//	numberToSignature := r.Header.Get("Number-To-Signature")
//
//	credsSignature := auth.Creds{
//		Values: map[auth.CredsType]string{
//			auth.CredsPublicKeyBase58:    publicKeyAddress,
//			auth.CredsContentToSignature: r.URL.Path + "?" + r.URL.RawQuery,
//			auth.CredsKeyToSignature:     numberToSignature,
//			auth.CredsSignature:          signature,
//		},
//	}
//
//	user, errs = auth.GetIdentity(credsSignature, authOps, errs)
//	// previous errs is added by auth.GetIdentity()
//}
//var errNoIdentityOpsMap = errors.New("no map[CredsType]identity.UserKey")
