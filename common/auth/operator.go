package auth

import (
	"github.com/pavlo67/workshop/common/identity"
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/joiner"
)

const AuthorizeHandlerKey joiner.InterfaceKey = "auth_handler"
const SetCredsHandlerKey joiner.InterfaceKey = "auth_set_creds_handler"

type User struct {
	Key   identity.Key `bson:",omitempty" json:",omitempty"`
	Creds Creds        `bson:",omitempty" json:",omitempty"`
}

type Operator interface {
	// SetCreds sets user's own or temporary (session-generated) creds
	SetCreds(userKey identity.Key, toSet Creds) (*Creds, error)

	// Authorize can require to do .SetCreds first and to usa some session-generated creds
	Authorize(toAuth Creds) (*User, error)
}

// to use with map[CredsType]identity.Actor  --------------------------------------------------------------------

var ErrNoIdentityOp = errors.New("no identity.Actor")

const onGetUser = "on GetUser()"

func GetUser(creds Creds, ops []Operator, errs common.Errors) (*User, common.Errors) {
	if len(creds) < 1 {
		return nil, append(errs, ErrNoCreds)
	}

	for _, op := range ops {
		user, err := op.Authorize(creds)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, onGetUser+`: on identOp.Authorize(%#v)`, creds))
		}

		if user != nil {
			return user, errs
		}
	}

	return nil, errs
}

// callbacks can be used for partial implementations of identity.Actor (in their own interfaces)
//
// type Callback string
//
// const Confirm Callback = "confirm"
// const SendCode Callback = "send_code"
//
// type Actor interface {
//	// Create stores registration data and (as usual) sends confirmation code to user.
//	Create(creds ...Creds) ([]Message, error)
//
//	AddCallback(key Callback, url string)
// }

//const Anyone common.ID = "_"

//type Access struct {
//	TargetID   Key     `bson:"target_id"             json:"target_id"`
//	TargetNick string `bson:"target_nick,omitempty" json:"target_nick,omitempty"`
//	Right      Right  `bson:"right,omitempty"       json:"right,omitempty"`
//}
