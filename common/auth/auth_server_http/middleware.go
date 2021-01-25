package auth_server_http

import (
	"net/http"
	"strings"

	"github.com/pavlo67/common/common"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/server/server_http"
)

var _ server_http.OnRequest = &onRequest{}

type onRequest struct{}

func (*onRequest) Options(r *http.Request) (*crud.Options, error) {
	//if r == nil {
	//	return nil, errors.New("no server_http.Request in RequestOptions(...)")
	//}

	var errorKey errors.Key
	var errs errors.Errors
	var identity *auth.Identity

	tokenJWT := r.Header.Get("Authorization")

	if tokenJWT != "" {
		tokenJWT = strings.Replace(tokenJWT, "Bearer ", "", 1)
		identity, errorKey, errs = auth.GetIdentity(auth.Creds{auth.CredsJWT: tokenJWT}, authOps, false, errs)

	} else {
		errorKey = errors.NoCredsErr

	}

	err := errors.KeyableError(errorKey, common.Map{"error": errs.Err()})
	if identity == nil {
		return nil, err
	}

	return &crud.Options{Identity: identity}, err
}

// TOKEN CHECK
//token := r.Header.Get("Token")
//if token != "" {
//	user, errs = auth.GetIdentity(auth.Creds{auth.CredsToken: token}, authOps, errs)
//	if user != nil {
//		return user, errs.Error()
//	}
//	// previous errs is added with auth.GetIdentity()
//}
//

//// COOKIE CHECK
//c, _ := r.Cookie("Token") // ErrNoCookie only
//if c != nil && c.Left != "" {
//	user, errs = auth.GetIdentity([]auth.Creds{{TypeKey: auth.CredsToken, Left: c.Left}}, authOps, errs)
//	if user != nil {
//		return user, errs.Error()
//	}
//	// previous errs is added with auth.GetIdentity()
//}

//// SIGNATURE CHECK
//signature := r.Header.Get("Signature")
//if signature != "" && r.URL != nil {
//	publicKeyAddress := r.Header.Get("Public-IDStr-Address")
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
