package packs_pg

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/libraries/sqllib"
	"github.com/pavlo67/workshop/common/libraries/sqllib/sqllib_pg"
	"github.com/pavlo67/workshop/common/selectors"
	"github.com/pavlo67/workshop/common/selectors/selectors_sql"
)

var _ crud.Cleaner = &packsPg{}

const onClean = "on packsPg.Clean(): "

func (packsOp *packsPg) Clean(term *selectors.Term, _ *crud.RemoveOptions) error {
	condition, values, err := selectors_sql.Use(term)
	if err != nil {
		return errors.Errorf(onClean+"wrong selector (%#v): %s", term, err)
	}

	query := packsOp.sqlClean
	if strings.TrimSpace(condition) != "" {
		query += " WHERE " + sqllib_pg.CorrectWildcards(condition)
	}

	_, err = packsOp.db.Exec(query, values...)
	if err != nil {
		return errors.Wrapf(err, onClean+sqllib.CantExec, query, values)
	}

	return nil
}

const onSelectToClean = "on packsSQLite.SelectToClean(): "

func (packsOp *packsPg) SelectToClean(options *crud.RemoveOptions) (*selectors.Term, error) {

	var limit uint64

	if options != nil && options.Limit > 0 {
		limit = options.Limit
	} else {
		return nil, errors.New(onSelectToClean + "no clean limit is defined")
	}

	queryMax := "SELECT MAX(id) from " + packsOp.table

	var maxID uint64
	row := packsOp.db.QueryRow(queryMax)

	err := row.Scan(&maxID)
	if err != nil {
		return nil, errors.Errorf(onSelectToClean+": error on query (%s)", queryMax)
	}

	return selectors.Binary(selectors.Le, "id", selectors.Value{V: maxID - limit}), nil

}
