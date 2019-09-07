package server_http

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/basis/auth"
	"github.com/pavlo67/workshop/basis/common"
)

var errNoIdentityOpsMap = errors.New("no map[CredsType]identity.Operator")

func UserWithRequest(r *http.Request, authOps []auth.Operator) (*auth.User, error) {

	var errs common.Errors
	var user *auth.User

	// TOKEN CHECK
	token := r.Header.Get("Token")
	if token != "" {
		user, errs = auth.GetUser([]auth.Creds{{Type: auth.CredsToken, Value: token}}, authOps, errs)
		if user != nil {
			return user, errs.Err()
		}
		// previous errs is added with auth.GetUser()
	}

	tokenJWT := r.Header.Get("JWT")
	if tokenJWT != "" {
		user, errs = auth.GetUser([]auth.Creds{{Type: auth.CredsJWT, Value: token}}, authOps, errs)
		if user != nil {
			return user, errs.Err()
		}
		// previous errs is added with auth.GetUser()
	}

	//// COOKIE CHECK
	//c, _ := r.Cookie("Token") // ErrNoCookie only
	//if c != nil && c.Value != "" {
	//	user, errs = auth.GetUser([]auth.Creds{{Type: auth.CredsToken, Value: c.Value}}, authOps, errs)
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

		credsSignature := []auth.Creds{
			{Type: auth.CredsPublicKeyAddress, Value: publicKeyAddress},
			{Type: auth.CredsContentToSignature, Value: r.URL.Path + "?" + r.URL.RawQuery},
			{Type: auth.CredsNumberToSignature, Value: numberToSignature},
			{Type: auth.CredsSignature, Value: signature},
		}

		user, errs = auth.GetUser(credsSignature, authOps, errs)
		// previous errs is added by auth.GetUser()
	}

	return user, errs.Err()
}
