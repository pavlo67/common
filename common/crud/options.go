package crud

import (
	"database/sql"

	"github.com/pavlo67/workshop/common/identity"
)

type SaveOptions struct {
	ActorKey identity.Key

	// TODO!!! use some general (non-SQL-specific) interface
	Tx *sql.Tx

	// TODO??? check if item.Key exists and if it should be existing (insert vs. replace)
}

type GetOptions struct {
	ActorKey identity.Key
	GroupBy  []string
	OrderBy  []string
	Offset   uint64
	Limit    uint64
}

type RemoveOptions struct {
	ActorKey identity.Key
	Limit    uint64
	Delete   bool
}
