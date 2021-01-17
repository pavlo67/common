package auth

import (
	"errors"
	"fmt"

	"github.com/pavlo67/workshop/common"
)

//type User struct {
//	Key   Key   `bson:",omitempty" json:",omitempty"`
//	Creds Creds `bson:",omitempty" json:",omitempty"`
//}

//func (user *User) KeyYet() Key {
//	if user == nil {
//		return ""
//	}
//
//	return user.Key
//}

type Operator interface {
	// SetCreds sets user's own or temporary (session-generated) creds
	SetCreds(userID ID, toSet Creds) (*Creds, error)

	// Authenticate can require to do .SetCreds first and to usa some session-generated creds
	Authenticate(toAuth Creds) (*Identity, error)
}

// to use with map[CredsType]identity.ActorKey  --------------------------------------------------------------------

var ErrNoIdentityOp = errors.New("no identity.ActorKey")

const onGetIdentity = "on GetIdentity()"

func GetIdentity(creds Creds, ops []Operator, useOperatorAuth bool, errs common.Errors) (*Identity, common.ErrorKey, common.Errors) {
	if len(creds) < 1 {
		return nil, common.NoCredsErr, append(errs, ErrNoCreds)
	}

	for _, op := range ops {
		identity, err := op.Authenticate(creds)
		if err != nil {
			errs = append(errs, fmt.Errorf(onGetIdentity+`: on identOp.Authenticate(%#v): %s`, creds, err))
		}
		if identity != nil {
			return identity, "", errs
		}

		//realm := op.Realm()
		//if (useOperatorAuth && realm == OperatorRealmKey) || (!useOperatorAuth && realm != OperatorRealmKey) {
		//	identity, err := op.Authenticate(creds)
		//	if err != nil {
		//		errs = append(errs, fmt.Errorf(onGetIdentity+`: on identOp.Authenticate(%#v): %s`, creds, err))
		//	}
		//	if identity != nil {
		//		return identity, "", errs
		//	}
		//}
	}

	return nil, common.InvalidCredsErr, errs
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
