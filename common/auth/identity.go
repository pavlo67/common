package auth

import (
	"encoding/json"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/rbac"
)

type ID common.IDStr

type Identity struct {
	ID       ID         `json:",omitempty"`
	Nickname string     `json:",omitempty"`
	Roles    rbac.Roles `json:",omitempty"`
	creds    common.Map `json:",omitempty"`
}

func (identity *Identity) Creds(key string) string {
	if identity == nil {
		return ""
	}

	return identity.creds.StringDefault(key, "")
}

func (identity *Identity) SetCreds(key string, value interface{}) {
	if identity == nil {
		return
	}

	if identity.creds == nil {
		identity.creds = common.Map{key: value}
		return
	}

	identity.creds[key] = value
}

type IdentityForMarshalling struct {
	ID       ID         `json:",omitempty"`
	Nickname string     `json:",omitempty"`
	Roles    rbac.Roles `json:",omitempty"`
	Creds    common.Map `json:",omitempty"`
}

func (identity Identity) MarshalJSON() ([]byte, error) {

	ifm := IdentityForMarshalling{
		ID:       identity.ID,
		Nickname: identity.Nickname,
		Roles:    identity.Roles,
		Creds:    identity.creds,
	}

	return json.Marshal(ifm)
}

func (identity *Identity) UnmarshalJSON(jsonBytes []byte) error {
	if len(jsonBytes) < 1 {
		return nil
	}

	var ifm IdentityForMarshalling
	if err := json.Unmarshal(jsonBytes, &ifm); err != nil {
		return err
	}

	identity.ID = ifm.ID
	identity.Nickname = ifm.Nickname
	identity.Roles = ifm.Roles
	identity.creds = ifm.Creds

	return nil
}
