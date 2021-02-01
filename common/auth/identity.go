package auth

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/rbac"
)

type ID = common.IDStr

type Identity struct {
	ID       ID         `json:",omitempty"`
	Nickname string     `json:",omitempty"`
	Roles    rbac.Roles `json:",omitempty"`
	Creds    common.Map `json:",omitempty"`
}
