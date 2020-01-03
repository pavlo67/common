package data_sqlite

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/libraries/sqllib"
	"github.com/pavlo67/workshop/common/selectors"
	"github.com/pavlo67/workshop/common/selectors/logic"
	"github.com/pavlo67/workshop/common/selectors/selectors_sql"
	"github.com/pavlo67/workshop/components/flowcleaner"
)

var _ crud.Cleaner = &dataSQLite{}

const onIDs = "on dataSQLite.IDs()"

func (dataOp *dataSQLite) ids(condition string, values []interface{}) ([]interface{}, error) {
	if strings.TrimSpace(condition) != "" {
		condition = " WHERE " + condition
	}

	query := "SELECT id FROM " + dataOp.table + condition
	stm, err := dataOp.db.Prepare(query)
	if err != nil {
		return nil, errors.Wrapf(err, onIDs+": can't db.Prepare(%s)", query)
	}

	rows, err := stm.Query(values...)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, onIDs+sqllib.CantQuery, query, values)
	}
	defer rows.Close()

	var ids []interface{}

	for rows.Next() {
		var id common.ID

		err := rows.Scan(&id)
		if err != nil {
			return ids, errors.Wrapf(err, onIDs+sqllib.CantScanQueryRow, query, values)
		}

		ids = append(ids, id)
	}
	err = rows.Err()
	if err != nil {
		return ids, errors.Wrapf(err, onIDs+": "+sqllib.RowsError, query, values)
	}

	return ids, nil
}

const onClean = "on dataSQLite.Clean(): "

func (dataOp *dataSQLite) Clean(term *selectors.Term, _ *crud.RemoveOptions) error {
	var termTags *selectors.Term

	condition, values, err := selectors_sql.Use(term)
	if err != nil {
		return errors.Errorf(onClean+"wrong selector (%#v): %s", term, err)
	}

	query := dataOp.sqlClean

	if strings.TrimSpace(condition) != "" {
		ids, err := dataOp.ids(condition, values)
		if err != nil {
			return errors.Wrap(err, onClean+"can't dataOp.ids(condition, values)")
		}
		termTags = logic.AND(selectors.In("key", dataOp.interfaceKey), selectors.In("id", ids...))

		query += " WHERE " + condition

	} else {
		termTags = selectors.In("key", dataOp.interfaceKey) // TODO!!! correct field key

	}

	_, err = dataOp.db.Exec(query, values...)
	if err != nil {
		return errors.Wrapf(err, onClean+sqllib.CantExec, query, values)
	}

	if dataOp.taggerCleaner != nil {
		err = dataOp.taggerCleaner.Clean(termTags, nil)
		if err != nil {
			return errors.Wrap(err, onClean)
		}
	}

	return err
}

func (dataOp *dataSQLite) SelectToClean(options *crud.RemoveOptions) (*selectors.Term, error) {
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
