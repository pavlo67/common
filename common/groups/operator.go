package groups

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/identity"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/selectors"
)

const InterfaceKey joiner.InterfaceKey = "groups"

type Item struct {
	auth.User `bson:",omitempty" json:",omitempty"`

	Title   string
	Details string

	OwnerKey  identity.Key `bson:",omitempty" json:",omitempty"`
	ViewerKey identity.Key `bson:",omitempty" json:",omitempty"`
	History   crud.History `bson:",omitempty" json:",omitempty"`
}

type Right string

// Operator is a common interface to implement an equality of user rights
type Operator interface {
	Save(Item, *crud.SaveOptions) (identity.Key, error)
	Remove(identity.Key, *crud.RemoveOptions) error

	Read(identity.Key, *crud.GetOptions) (*Item, error)
	List(*selectors.Term, *crud.GetOptions) ([]Item, error)
	// Count(*selectors.Term, *crud.GetOptions) (uint64, error)

	UserAccesses(user identity.Key) ([]identity.Name, error)
	UserRights(user, object identity.Key) ([]Right, error)

	//SetRights(is, groupIS, memberIS auth.ID, rightsToTime map[rights.Right]*time.Time) error
	//AddRight(adminID, userID, objectID common.ID, right Right) error
	//RemoveRight(adminID, userID, objectID common.ID, right Right) error
}

const onHasAccessTo = "on groups.HasAccessTo(): "

func HasAccessTo(groupsOp Operator, user, object identity.Key, hasEspecialRight Right) (bool, error) {
	if groupsOp == nil {
		return false, errors.New(onHasAccessTo + "no groups.Operator")
	}

	rights, err := groupsOp.UserRights(user, object)
	if err != nil {
		return false, errors.Wrapf(err, onHasAccessTo+"can't groupsOp.UserRights(%s, %s)", user, object)
	}

	if hasEspecialRight == "" {
		return true, nil
	}

	for _, right := range rights {
		if right == hasEspecialRight {
			return true, nil
		}
	}

	return false, nil
}
