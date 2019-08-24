package auth

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/basis/common"
	"github.com/pavlo67/workshop/basis/joiner"
)

const InterfaceKey joiner.InterfaceKey = "auth"

const Anyone common.ID = "_"

//type Access struct {
//	TargetID   ID     `bson:"target_id"             json:"target_id"`
//	TargetNick string `bson:"target_nick,omitempty" json:"target_nick,omitempty"`
//	Right      Right  `bson:"right,omitempty"       json:"right,omitempty"`
//}

type User struct {
	ID   common.ID `bson:"id"                 json:"id"`
	Nick string    `bson:"nick"               json:"nick"`
	// Accesses []Access `bson:"accesses,omitempty" json:"accesses,omitempty"`
}

type Operator interface {
	// Authorize can require multi-steps (using returned []Creds)...
	Authorize(toAuth []Creds) (*User, []Creds, error)

	// SetCreds can require multi-steps (using returned []Creds)...
	SetCreds(userID *common.ID, toSet []Creds) (*User, []Creds, error)

	Accepts() ([]CredsType, error)
}

// to use with map[CredsType]identity.Operator  --------------------------------------------------------------------

var errNoCreds = errors.New("no creds")
var errNoIdentityOp = errors.New("no identity.Operator")

const onGetUser = "on GetUser()"

func GetUser(creds []Creds, op Operator, errs common.Errors) (*User, common.Errors) {
	if len(creds) < 1 {
		return nil, append(errs, errNoCreds)
	}

	credsType := creds[0].Type

	if op == nil {
		return nil, append(errs, errors.Wrapf(errNoIdentityOp, onGetUser+": for Authorize with "+string(credsType)))
	}

	user, _, err := op.Authorize(creds)
	if err != nil {
		return nil, append(errs, errors.Wrapf(err, onGetUser+`: on identOp.Authorize(%#v)`, creds))
	}
	if user != nil {
		return user, errs
	}

	return nil, errs
}

// callbacks can be used for partial implementations of identity.Operator (in their own interfaces)
//
// type Callback string
//
// const Confirm Callback = "confirm"
// const SendCode Callback = "send_code"
//
// type Operator interface {
//	// Create stores registration data and (as usual) sends confirmation code to user.
//	Create(creds ...Creds) ([]Message, error)
//
//	AddCallback(key Callback, url string)
// }
