package server_http

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/auth"
	"github.com/pavlo67/punctum/basis"
)

var errNoIdentityOpsMap = errors.New("no map[CredsType]identity.Operator")

func UserWithRequest(r *http.Request, identOpsMap map[auth.CredsType][]auth.Operator) (*auth.User, error) {
	if identOpsMap == nil {
		return nil, errNoIdentityOpsMap
	}

	var errs basis.Errors
	var user *auth.User

	// TOKEN_CHECK
	token := r.Header.Get("Token")
	if token == "" {
		c, _ := r.Cookie("Token") // ErrNoCookie only
		if c == nil || c.Value == "" {
			goto SIGNATURE_CHECK
		}
	}

	user, errs = auth.GetUser([]auth.Creds{{Type: auth.CredsToken, Value: token}}, identOpsMap, errs)
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

	credsSignature := []auth.Creds{
		{Type: auth.CredsSignature, Value: signature},
		{Type: auth.CredsContentToSignature, Value: contentToSignature},
		{Type: auth.CredsPublicKeyAddress, Value: publicKeyAddress},
	}

	user, errs = auth.GetUser(credsSignature, identOpsMap, errs)
	return user, errs.Err()
}
