package auth

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/pavlo67/workshop/common/rbac"
)

func TestJSON(t *testing.T) {
	identity := Identity{
		ID:       "1",
		Nickname: "2",
		Roles:    rbac.Roles{rbac.RoleUser},
	}

	bytes, err := json.Marshal(identity)

	log.Printf("%s / %s", bytes, err)

}
