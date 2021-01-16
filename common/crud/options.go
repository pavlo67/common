package crud

import (
	"database/sql"

	"github.com/pavlo67/workshop/common/identity"
	"github.com/pavlo67/workshop/common/selectors"
)

type JoinTo struct {
	Clause string
	Values []interface{}
}

type Options struct {
	ActorKey identity.Key

	Term    *selectors.Term
	JoinTo  JoinTo
	GroupBy []string
	OrderBy []string
	Offset  int64

	Tx *sql.Tx // TODO!!! use some general (non-SQL-specific) interface

	Limit  int64
	Delete bool
}

//type Counter map[string]uint64
//
//type Index map[string][]ID
