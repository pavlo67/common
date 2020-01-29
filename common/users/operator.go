package users

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/identity"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/selectors"
)

const InterfaceKey joiner.InterfaceKey = "users"

const UserKeyFieldName = "key"
const EmailFieldName = "email"
const NicknameFieldName = "nickname"
const VerifiedFieldName = "verified"

type Item struct {
	auth.User `bson:",omitempty" json:",omitempty"`

	Allowed  bool           `bson:",omitempty" json:",omitempty"`
	ToVerify []Verification `bson:",omitempty" json:",omitempty"`
	History  crud.History   `bson:",omitempty" json:",omitempty"`
}

type Verification struct {
	CredsType auth.CredsType `bson:",omitempty" json:",omitempty"`
	Value     string         `bson:",omitempty" json:",omitempty"`
	Open      bool           `bson:",omitempty" json:",omitempty"`
	History   crud.History   `bson:",omitempty" json:",omitempty"`
}

type Operator interface {
	Save(Item, *crud.SaveOptions) (identity.Key, error)
	Remove(identity.Key, *crud.RemoveOptions) error

	Read(identity.Key, *crud.GetOptions) (*Item, error)
	List(*selectors.Term, *crud.GetOptions) ([]Item, error)
	Count(*selectors.Term, *crud.GetOptions) (uint64, error)

	CheckPassword(password, passHash string) bool

	Allow() error
	SetVerification(auth.CredsType, string, bool) error
	Verify(auth.CredsType, string, common.Errors) error
}
