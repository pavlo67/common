package server_http

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/auth"
)

var errNoIdentityOpsMap = errors.New("no map[CredsType]identity.Actor")

func UserWithRequest(r *http.Request, authOps []auth.Operator) (*auth.User, error) {

	var errs common.Errors
	var user *auth.User

	// TOKEN CHECK
	token := r.Header.Get("Token")
	if token != "" {
		user, errs = auth.GetUser(auth.Creds{Values: map[auth.CredsType]string{auth.CredsToken: token}}, authOps, errs)
		if user != nil {
			return user, errs.Err()
		}
		// previous errs is added with auth.GetUser()
	}

	tokenJWT := r.Header.Get("JWT")
	if tokenJWT != "" {
		user, errs = auth.GetUser(auth.Creds{Values: map[auth.CredsType]string{auth.CredsJWT: token}}, authOps, errs)
		if user != nil {
			return user, errs.Err()
		}
		// previous errs is added with auth.GetUser()
	}

	//// COOKIE CHECK
	//c, _ := r.Cookie("Token") // ErrNoCookie only
	//if c != nil && c.Left != "" {
	//	user, errs = auth.GetUser([]auth.Creds{{TypeKey: auth.CredsToken, Left: c.Left}}, authOps, errs)
	//	if user != nil {
	//		return user, errs.Err()
	//	}
	//	// previous errs is added with auth.GetUser()
	//}

	// SIGNATURE CHECK
	signature := r.Header.Get("Signature")
	if signature != "" && r.URL != nil {
		publicKeyAddress := r.Header.Get("Public-Key-Address")
		numberToSignature := r.Header.Get("Number-To-Signature")

		credsSignature := auth.Creds{
			Values: map[auth.CredsType]string{
				auth.CredsPublicKeyBase58:    publicKeyAddress,
				auth.CredsContentToSignature: r.URL.Path + "?" + r.URL.RawQuery,
				auth.CredsKeyToSignature:     numberToSignature,
				auth.CredsSignature:          signature,
			},
		}

		user, errs = auth.GetUser(credsSignature, authOps, errs)
		// previous errs is added by auth.GetUser()
	}

	return user, errs.Err()
}
