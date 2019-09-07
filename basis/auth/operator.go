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
	ID    common.ID `bson:"id"                json:"id"`
	Nick  string    `bson:"nick"              json:"nick"`
	Creds []Creds   `bson:"creds, omitempty" json:"creds, omitempty"`
	// Accesses []Access `bson:"accesses,omitempty" json:"accesses,omitempty"`
}

type Operator interface {
	// Authorize can require multi-steps (using returned []Creds)...
	Authorize(toAuth ...Creds) (*User, error)

	// SetCreds can require multi-steps (using returned []Creds)...
	SetCreds(user User, toSet ...Creds) ([]Creds, error)
}

// to use with map[CredsType]identity.Operator  --------------------------------------------------------------------

var errNoCreds = errors.New("no creds")
var errNoIdentityOp = errors.New("no identity.Operator")

const onGetUser = "on GetUser()"

func GetUser(creds []Creds, ops []Operator, errs common.Errors) (*User, common.Errors) {
	if len(creds) < 1 {
		return nil, append(errs, errNoCreds)
	}

	for _, op := range ops {
		user, err := op.Authorize(creds...)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, onGetUser+`: on identOp.Authorize(%#v)`, creds))
		}

		if user != nil {
			return user, errs
		}
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
