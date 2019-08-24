package rector

import (
	"github.com/pavlo67/workshop/basis/joiner"
	"github.com/pavlo67/workshop/basis/common"
)

const InterfaceKey joiner.InterfaceKey = "rector"

type Operator interface {
	AddRight(adminID, userID, objectID common.ID, right Right) error
	CheckRight(adminID, userID, objectID common.ID, right Right) (bool, error)
	RemoveRight(adminID, userID, objectID common.ID, right Right) error
}
