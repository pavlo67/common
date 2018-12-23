package identity

import (
	"github.com/pavlo67/punctum/basis"
	"github.com/pkg/errors"
)

var errNoIdentityOp = errors.New("no identity.Operator")

const onGetUser = "on GetUserData()"

func GetUserWithSignature(contentToSignature, publicKeyAddress, signature string, identOpsMap map[CredsType][]Operator, errs basis.Errors) (*User,
	basis.Errors) {
	if len(identOpsMap[CredsSignature]) < 1 {
		return nil, append(errs, errors.Wrap(errNoIdentityOp, "for authorize with CredsSignature"))
	}

	for _, identOp := range identOpsMap[CredsSignature] {
		if identOp == nil {
			errs = append(errs, errors.Wrapf(errNoIdentityOp, onGetUser+": for Authorize with CredsSignature"))
			continue
		}

		creds := []Creds{
			{Type: CredsContentToSignature, Value: contentToSignature},
			{Type: CredsPublicKeyAddress, Value: publicKeyAddress},
			{Type: CredsSignature, Value: signature},
		}

		user, _, err := identOp.Authorize(creds...)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, onGetUser+`: on identOp.Authorize(%#v)`, creds))
		}
		if user != nil {
			return user, errs
		}
	}

	return nil, errs
}

func GetUserWithToken(token string, identOpsMap map[CredsType][]Operator, errs basis.Errors) (*User, basis.Errors) {
	if len(identOpsMap[CredsToken]) < 1 {
		return nil, append(errs, errors.Wrap(errNoIdentityOp, "for authorize with CredsToken"))
	}

	for _, identOp := range identOpsMap[CredsToken] {
		if identOp == nil {
			errs = append(errs, errors.Wrapf(errNoIdentityOp, onGetUser+": for Authorize with CredsToken"))
			continue
		}

		creds := Creds{Type: CredsToken, Value: token}
		user, _, err := identOp.Authorize(creds)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, onGetUser+`: on identOp.Authorize(%#v)`, creds))
		}
		if user != nil {
			return user, errs
		}
	}

	return nil, errs
}
