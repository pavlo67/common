package rights

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/joiner"
)

const InterfaceKey joiner.InterfaceKey = "rector"

type Operator interface {
	AddRight(adminID, userID, objectID common.ID, right Right) error
	CheckRight(adminID, userID, objectID common.ID, right Right) (bool, error)
	RemoveRight(adminID, userID, objectID common.ID, right Right) error
}
