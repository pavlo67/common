package identity

import (
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/program"
)

const InterfaceKey program.InterfaceKey = "identity"

type User struct {
	ID       basis.ID `bson:"id"       json:"id"`
	Nickname string   `bson:"nickname" json:"nickname"`
}

type Operator interface {
	// Use can require multi-steps...
	// It autenticates the user (using his own or admin's creds) and can also add, update or remove user's creds.
	// It can create a new user if no toAuth-creds sent
	Identify(toAuth []Creds, toSet ...Creds) (*User, basis.Values, error)
	HasRights(user *User, allowedIDs ...basis.ID) bool
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
