package identity

import (
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/program"
)

const InterfaceKey program.InterfaceKey = "identity"

type Access struct {
	ID    basis.ID `bson:"id"              json:"id"`
	Right Right    `bson:"right,omitempty" json:"right,omitempty"`
	Label string   `bson:"label,omitempty" json:"label,omitempty"`
}

type User struct {
	ID       basis.ID `bson:"id"                 json:"id"`
	Nickname string   `bson:"nickname"           json:"nickname"`
	Accesses []Access `bson:"accesses,omitempty" json:"accesses,omitempty"`
}

type Operator interface {
	// SetCreds can require multi-steps (using returned []Creds)...
	SetCreds(userID *basis.ID, toSet []Creds, toAuth ...Creds) (*User, []Creds, error)

	// Authorize can require multi-steps (using returned []Creds)...
	Authorize(toAuth ...Creds) (*User, []Creds, error)

	Accepts() ([]CredsType, error)
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
