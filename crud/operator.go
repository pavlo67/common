package crud

import (
	"github.com/pavlo67/punctum/auth"
	"github.com/pavlo67/punctum/basis/selectors"
	"github.com/pavlo67/punctum/starter/joiner"
)

const InterfaceKey joiner.InterfaceKey = "crud"

type ReadOptions struct {
	Selector *selectors.Term `json:"selector,omitempty"`
	SortBy   []string        `json:"sort_by,omitempty"`
}

// Operator is a common interface to manage create/read/update/delete operations
type Operator interface {
	Mapper

	Create(userIS auth.ID, native interface{}) (id string, err error)

	// Read returns crud item (accordingly to requester's rights).
	Read(userIS auth.ID, id string) (interface{}, error)

	// ReadList returns crud items list (accordingly to requester's rights).
	ReadList(userIS auth.ID, options ReadOptions) ([]interface{}, *uint64, error)

	// Update changes crud item (accordingly to requester's rights).
	Update(userIS auth.ID, id string, native interface{}) error

	// Update deletes crud item (accordingly to requester's rights).
	Delete(userIS auth.ID, id string) error
}

type Cleaner func() error
