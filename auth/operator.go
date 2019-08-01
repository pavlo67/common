package auth

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/associatio/basis"
	"github.com/pavlo67/associatio/starter/joiner"
)

const InterfaceKey joiner.InterfaceKey = "auth"

type ID string

const Anyone ID = "_"

//type Access struct {
//	TargetID   ID     `bson:"target_id"             json:"target_id"`
//	TargetNick string `bson:"target_nick,omitempty" json:"target_nick,omitempty"`
//	Right      Right  `bson:"right,omitempty"       json:"right,omitempty"`
//}

type User struct {
	ID   ID     `bson:"id"                 json:"id"`
	Nick string `bson:"nick"               json:"nick"`
	// Accesses []Access `bson:"accesses,omitempty" json:"accesses,omitempty"`
}

type Operator interface {
	// SetCreds can require multi-steps (using returned []Creds)...
	SetCreds(userID *ID, toSet []Creds, toAuth ...Creds) (*User, []Creds, error)

	// Authorize can require multi-steps (using returned []Creds)...
	Authorize(toAuth ...Creds) (*User, []Creds, error)

	Accepts() ([]CredsType, error)
}

// to use with map[CredsType]identity.Operator  --------------------------------------------------------------------

var errNoCreds = errors.New("no creds")
var errNoIdentityOp = errors.New("no identity.Operator")

const onGetUser = "on GetUser()"

func GetUser(creds []Creds, identOpsMap map[CredsType][]Operator, errs basis.Errors) (*User, basis.Errors) {
	if len(creds) < 1 {
		return nil, append(errs, errNoCreds)
	}
	credsType := creds[0].Type

	if len(identOpsMap[credsType]) < 1 {
		return nil, append(errs, errors.Wrap(errNoIdentityOp, "for authorize with "+string(credsType)))
	}

	for _, identOp := range identOpsMap[credsType] {
		if identOp == nil {
			errs = append(errs, errors.Wrapf(errNoIdentityOp, onGetUser+": for Authorize with "+string(credsType)))
			continue
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
