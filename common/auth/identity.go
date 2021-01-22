package auth

import (
	"time"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/rbac"
)

type ID = common.IDStr

type Identity struct {
	ID           ID         `json:",omitempty"`
	Nickname     string     `json:",omitempty"`
	Roles        rbac.Roles `json:",omitempty"`
	JWT          string     `json:",omitempty"`
	RefreshToken string     `json:",omitempty"`
	ExpiredAt    *time.Time `json:"-"`
	ReAuthData   common.Map `json:"-"`
}
