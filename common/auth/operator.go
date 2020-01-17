package auth

import (
	"github.com/pavlo67/workshop/common/identity"
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/joiner"
)

const InterfaceKey joiner.InterfaceKey = "auth"

const Anyone common.ID = "_"

//type Access struct {
//	TargetID   Key     `bson:"target_id"             json:"target_id"`
//	TargetNick string `bson:"target_nick,omitempty" json:"target_nick,omitempty"`
//	Right      Right  `bson:"right,omitempty"       json:"right,omitempty"`
//}

type User struct {
	Key      identity.Key `bson:",omitempty" json:",omitempty"`
	Nickname string       `bson:",omitempty" json:",omitempty"`
	Creds    Creds        `bson:",omitempty" json:",omitempty"`
}

type Operator interface {
	// SetCreds can require multi-steps (using returned Creds)...
	SetCreds(user *User, toSet Creds) (*User, *Creds, error)

	// InitAuthSession starts an auth session if it's required
	InitAuthSession(toInit Creds) (*Creds, error)

	// Authorize can require multi-steps (using returned Creds)...
	Authorize(toAuth Creds) (*User, error)
}

// to use with map[CredsType]identity.Actor  --------------------------------------------------------------------

var ErrNoIdentityOp = errors.New("no identity.Actor")

const onGetUser = "on GetUser()"

func GetUser(creds Creds, ops []Operator, errs common.Errors) (*User, common.Errors) {
	if len(creds.Values) < 1 {
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
