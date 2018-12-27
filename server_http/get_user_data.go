package server_http

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/identity"
)

var errNoIdentityOpsMap = errors.New("no map[CredsType]identity.Operator")

var errEmptyToken = errors.New("empty token")

const onGetUserData = "on UserWithRequest()"

func UserWithRequest(r *http.Request, identOpsMap map[identity.CredsType][]identity.Operator) (*identity.User, error) {
	if identOpsMap == nil {
		return nil, errNoIdentityOpsMap
	}

	var errs basis.Errors
	var user *identity.User

	// TOKEN_CHECK
	token := r.Header.Get("Token")
	if token == "" {
		c, _ := r.Cookie("Token") // ErrNoCookie only
		if c == nil || c.Value == "" {
			goto SIGNATURE_CHECK
		}
	}

	user, errs = identity.GetUserWithToken(token, identOpsMap, errs)
	if user != nil {
		return user, errs.Err()
	}

SIGNATURE_CHECK:
	signature := r.Header.Get("Signature")
	if signature == "" {
		return nil, errs.Err()
	}

	contentToSignature := r.Header.Get("Content-To-Signature") + r.RemoteAddr
	publicKeyAddress := r.Header.Get("Public-Key-Address")

	user, errs = identity.GetUserWithSignature(contentToSignature, publicKeyAddress, signature, identOpsMap, errs)
	return user, errs.Err()
}
