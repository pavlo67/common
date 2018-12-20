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
	Read(token string) (*User, error)
	Allow(user *User, allowedIDs ...basis.ID) bool
}
