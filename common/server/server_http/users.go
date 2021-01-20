package server_http

import (
	"net/http"
	"strings"

	"github.com/pavlo67/workshop/common/auth"

	"github.com/pavlo67/workshop/common/errors"
)

var errNoIdentityOpsMap = errors.New("no map[CredsType]identity.UserKey")

func IdentityWithRequest(r *http.Request, authOps []auth.Operator) (*auth.Identity, errors.Key, error) {

	var errs errors.Errors
	var errorKey errors.Key
	var identity *auth.Identity

	tokenJWT := r.Header.Get("Authorization")

	if tokenJWT != "" {
		tokenJWT = strings.Replace(tokenJWT, "Bearer ", "", 1)

		var postfix string
		if len(tokenJWT) >= len(OperatorJWTKey) && tokenJWT[len(tokenJWT)-len(OperatorJWTKey):] == OperatorJWTKey {
			postfix = OperatorJWTKey
			identity, errorKey, errs = auth.GetIdentity(auth.Creds{auth.CredsJWT: tokenJWT[:len(tokenJWT)-len(OperatorJWTKey)]}, authOps, true, errs)
		} else {
			identity, errorKey, errs = auth.GetIdentity(auth.Creds{auth.CredsJWT: tokenJWT}, authOps, false, errs)
		}

		if identity != nil {
			identity.JWT += postfix
		}
		// previous errs is added with auth.GetIdentity()

	} else {
		errorKey = errors.NoCredsErr

	}

	return identity, errorKey, errs.Err()
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
