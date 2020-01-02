package flowcleaner_sqlite

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/libraries/sqllib/sqllib_sqlite"
	"github.com/pavlo67/workshop/common/selectors"

	"github.com/pavlo67/workshop/components/flow"
	"github.com/pavlo67/workshop/components/flowcleaner"
)

var _ crud.Cleaner = &flowCleanerSQLite{}

type flowCleanerSQLite struct {
	db        *sql.DB
	table     string
	tableTags string

	interfaceKey joiner.InterfaceKey
}

const onNew = "on flowCleanerSQLite.New(): "

func New(access config.Access, table, tableTags string, interfaceKey joiner.InterfaceKey) (crud.Cleaner, error) {
	db, err := sqllib_sqlite.Connect(access)
	if err != nil {
		return nil, errors.Wrap(err, onNew)
	}

	if table == "" {
		table = flow.CollectionDefault
	}

	fcOp := flowCleanerSQLite{
		db:        db,
		table:     table,
		tableTags: tableTags,

		interfaceKey: interfaceKey,
	}

	return &fcOp, nil
}

const onClean = "on flowCleanerSQLite.Clean(): "

func (fcOp *flowCleanerSQLite) Clean(_ *selectors.Term, options *crud.RemoveOptions) error {
	var limit uint64 = flowcleaner.FlowLimitDefault

	if options != nil && options.Limit > 0 {
		limit = options.Limit
	}

	queryMax := "SELECT MAX(id) from " + fcOp.table

	var maxID uint64
	row := fcOp.db.QueryRow(queryMax)

	err := row.Scan(&maxID)
	if err != nil {
		return errors.Errorf(onClean+": error on query (%s)", queryMax)
	}

	queryDelete := "DELETE from " + fcOp.table + " WHERE id <= ?"
	res, err := fcOp.db.Exec(queryDelete, maxID-limit)
	if err != nil {
		return errors.Errorf(onClean+": error on query (%s)", queryDelete)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Errorf(onClean+": error on res.RowsAffected(%s)", queryDelete)
	}

	l.Infof(onClean+": res.RowsAffected() = %d", rowsAffected)

	if fcOp.tableTags != "" {
		// TODO!!!
	}

	return nil
}

func (fcOp *flowCleanerSQLite) Close() error {
	return errors.Wrap(fcOp.db.Close(), "on flowCleanerSQLite.Close()")
}
