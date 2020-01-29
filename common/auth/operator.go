package auth

import (
	"errors"
	"fmt"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/identity"
	"github.com/pavlo67/workshop/common/joiner"
)

const AuthorizeHandlerKey joiner.InterfaceKey = "authorize_handler"
const SetCredsHandlerKey joiner.InterfaceKey = "set_creds_handler"
const GetCredsHandlerKey joiner.InterfaceKey = "get_creds_handler"

type User struct {
	Key   identity.Key `bson:",omitempty" json:",omitempty"`
	Creds Creds        `bson:",omitempty" json:",omitempty"`
}

func (user *User) KeyYet() identity.Key {
	if user == nil {
		return ""
	}

	return user.Key
}

type Operator interface {
	// SetCreds sets user's own or temporary (session-generated) creds
	SetCreds(userKey identity.Key, toSet Creds) (*Creds, error)

	// Authorize can require to do .SetCreds first and to usa some session-generated creds
	Authorize(toAuth Creds) (*User, error)
}

// to use with map[CredsType]identity.ActorKey  --------------------------------------------------------------------

var ErrNoIdentityOp = errors.New("no identity.ActorKey")

const onGetUser = "on GetUser()"

func GetUser(creds Creds, ops []Operator, errs common.Errors) (*User, common.Errors) {
	if len(creds) < 1 {
		return nil, append(errs, ErrNoCreds)
	}

	for _, op := range ops {
		user, err := op.Authorize(creds)
		if err != nil {
			errs = append(errs, fmt.Errorf(onGetUser+`: on identOp.Authorize(%#v): %s`, creds, err))
		}

		if user != nil {
			return user, errs
		}
	}

	return nil, errs
}

// callbacks can be used for partial implementations of identity.ActorKey (in their own interfaces)
//
// type Callback string
//
// const Confirm Callback = "confirm"
// const SendCode Callback = "send_code"
//
// type ActorKey interface {
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
