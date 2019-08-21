package administrator

import (
	"github.com/pavlo67/constructor/components/common"
)

type Operator interface {
	AddRight(adminID, userID, objectID common.ID, right Right) error
	CheckRight(adminID, userID, objectID common.ID, right Right) (bool, error)
	RemoveRight(adminID, userID, objectID common.ID, right Right) error
}
