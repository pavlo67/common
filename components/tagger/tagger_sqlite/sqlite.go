package tagger_sqlite

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/libraries/sqllib"
	"github.com/pavlo67/workshop/common/libraries/sqllib/sqllib_sqlite"
	"github.com/pavlo67/workshop/common/selectors"
	"github.com/pavlo67/workshop/components/tagger"
)

const tableDefault = "tags"

var fieldsToSave = []string{"key", "id", "tag"}
var fieldsToSaveStr = strings.Join(fieldsToSave, ", ")

var _ tagger.Operator = &taggerSQLite{}
var _ crud.Cleaner = &taggerSQLite{}

type taggerSQLite struct {
	db    *sql.DB
	table string

	sqlSaveTags, sqlRemoveTags, sqlRemoveAllTags, sqlListTags, sqlListTagged, sqlCountTagged, sqlCountTaggedAll string
	stmSaveTags, stmRemoveTags, stmRemoveAllTags, stmListTags, stmListTagged, stmCountTagged, stmCountTaggedAll *sql.Stmt
}

const onNew = "on taggerSQLite.New(): "

func NewTagger(access config.Access, table string) (tagger.Operator, crud.Cleaner, error) {
	db, err := sqllib_sqlite.Connect(access)
	if err != nil {
		return nil, nil, errors.Wrap(err, onNew)
	}

	if table == "" {
		table = tableDefault
	}

	taggerOp := taggerSQLite{
		db:    db,
		table: table,

		sqlSaveTags:      "INSERT OR IGNORE INTO " + table + " (" + fieldsToSaveStr + ") VALUES (" + strings.Repeat(",? ", len(fieldsToSave))[1:] + ")",
		sqlRemoveTags:    "DELETE FROM " + table + " where key = ? AND id = ? and tag = ?",
		sqlRemoveAllTags: "DELETE FROM " + table + " where key = ? AND id = ?",

		sqlListTags:       "SELECT tag, count(*) FROM " + table + " WHERE key = ? AND id = ? ORDER BY tag",
		sqlCountTagged:    "SELECT tag, count(*) FROM " + table + " WHERE key = ?            GROUP BY tag ORDER BY tag",
		sqlCountTaggedAll: "SELECT tag, count(*) FROM " + table + "            GROUP BY tag ORDER BY tag",
		sqlListTagged:     "SELECT key, id FROM       " + table + " WHERE tag = ?            ORDER BY key, id",
	}

	sqlStmts := []sqllib.SqlStmt{
		{&taggerOp.stmSaveTags, taggerOp.sqlSaveTags},
		{&taggerOp.stmRemoveTags, taggerOp.sqlRemoveTags},
		{&taggerOp.stmRemoveAllTags, taggerOp.sqlRemoveAllTags},

		{&taggerOp.stmListTags, taggerOp.sqlListTags},
		{&taggerOp.stmListTagged, taggerOp.sqlListTagged},
		{&taggerOp.stmCountTagged, taggerOp.sqlCountTagged},
		{&taggerOp.stmCountTaggedAll, taggerOp.sqlCountTaggedAll},
	}

	for _, sqlStmt := range sqlStmts {
		if err := sqllib.Prepare(db, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, nil, errors.Wrap(err, onNew)
		}
	}

	return &taggerOp, &taggerOp, nil
}

const onListTags = "on taggerSQLite.ListTags(): "

func (taggerOp *taggerSQLite) ListTags(key joiner.InterfaceKey, id common.ID, _ *crud.GetOptions) ([]tagger.Tag, error) {
	values := []interface{}{key, id}

	rows, err := taggerOp.stmListTags.Query(values...)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, onListTags+sqllib.CantQuery, taggerOp.sqlListTags, values)
	}
	defer rows.Close()

	var tags []tagger.Tag

	for rows.Next() {
		var tag string

		err = rows.Scan(&tag)
		if err != nil {
			return tags, errors.Wrapf(err, onListTags+sqllib.CantScanQueryRow, taggerOp.sqlListTags, values)
		}

		tags = append(tags, tagger.Tag(tag))
	}
	err = rows.Err()
	if err != nil {
		return tags, errors.Wrapf(err, onListTags+": "+sqllib.RowsError, taggerOp.sqlListTags, values)
	}

	return tags, nil
}

const onSaveTags = "on taggerSQLite.SaveTags(): "

func (taggerOp *taggerSQLite) SaveTags(key joiner.InterfaceKey, id common.ID, tags []tagger.Tag, _ *crud.SaveOptions) error {
	for _, tag := range tags {
		values := []interface{}{key, id, tag}

		_, err := taggerOp.stmSaveTags.Exec(values...)
		if err != nil {
			return errors.Wrapf(err, onSaveTags+sqllib.CantExec, taggerOp.sqlSaveTags, values)
		}
	}

	return nil
}

const onReplaceTags = "on taggerSQLite.ReplaceTags(): "

func (taggerOp *taggerSQLite) ReplaceTags(key joiner.InterfaceKey, id common.ID, tags []tagger.Tag, options *crud.SaveOptions) error {
	values := []interface{}{key, id}

	_, err := taggerOp.stmRemoveAllTags.Exec(values...)
	if err != nil {
		return errors.Wrapf(err, onSaveTags+sqllib.CantExec, taggerOp.sqlRemoveAllTags, values)
	}

	err = taggerOp.SaveTags(key, id, tags, options)
	if err != nil {
		return errors.Wrap(err, onReplaceTags)
	}

	return nil
}

const onRemoveTags = "on taggerSQLite.RemoveTags(): "

func (taggerOp *taggerSQLite) RemoveTags(key joiner.InterfaceKey, id common.ID, tags []tagger.Tag, _ *crud.SaveOptions) error {
	for _, tag := range tags {
		values := []interface{}{key, id, tag}

		_, err := taggerOp.stmRemoveTags.Exec(values...)
		if err != nil {
			return errors.Wrapf(err, onRemoveTags+sqllib.CantExec, taggerOp.sqlRemoveTags, values)
		}
	}

	return nil
}

const onCleanTags = "on taggerSQLite.CleanTags(): "

func (taggerOp *taggerSQLite) CleanTags(joiner.InterfaceKey, *selectors.Term, *crud.SaveOptions) error {
	//for _, tag := range tags {
	//	values := []interface{}{key, id, tag}
	//
	//	_, err := taggerOp.stmRemoveTags.Exec(values...)
	//	if err != nil {
	//		return errors.Wrapf(err, onRemoveTags+sqllib.CantExec, taggerOp.sqlRemoveTags, values)
	//	}
	//}

	// TODO!!!
	return common.ErrNotImplemented
}

const onCountTagged = "on taggerSQLite.CountTagged(): "

func (taggerOp *taggerSQLite) CountTagged(key *joiner.InterfaceKey, _ *crud.GetOptions) (crud.Counter, error) {
	var values []interface{}
	var query string
	var stm *sql.Stmt

	if key == nil {
		query = taggerOp.sqlCountTaggedAll
		stm = taggerOp.stmCountTaggedAll
	} else {
		values = []interface{}{*key}
		query = taggerOp.sqlCountTagged
		stm = taggerOp.stmCountTagged
	}

	rows, err := stm.Query(values...)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, onCountTagged+sqllib.CantQuery, query, values)
	}
	defer rows.Close()

	counter := crud.Counter{}

	for rows.Next() {
		var key string
		var count uint64

		err = rows.Scan(&key, &count)
		if err != nil {
			return counter, errors.Wrapf(err, onCountTagged+sqllib.CantScanQueryRow, query, values)
		}

		counter[key] = count
	}
	err = rows.Err()
	if err != nil {
		return counter, errors.Wrapf(err, onCountTagged+": "+sqllib.RowsError, query, values)
	}

	return counter, nil
}

const onIndexWithTag = "on taggerSQLite.IndexWithTag()"

func (taggerOp *taggerSQLite) IndexWithTag(tag tagger.Tag, _ *crud.GetOptions) (crud.Index, error) {
	values := []interface{}{tag}

	rows, err := taggerOp.stmListTagged.Query(values...)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, onIndexWithTag+sqllib.CantQuery, taggerOp.sqlListTagged, values)
	}
	defer rows.Close()

	index := crud.Index{}

	for rows.Next() {
		var key, id string
		err = rows.Scan(&key, &id)
		if err != nil {
			return index, errors.Wrapf(err, onIndexWithTag+sqllib.CantScanQueryRow, taggerOp.sqlListTagged, values)
		}

		index[key] = append(index[key], common.ID(id))
	}
	err = rows.Err()
	if err != nil {
		return index, errors.Wrapf(err, onIndexWithTag+": "+sqllib.RowsError, taggerOp.sqlListTagged, values)
	}

	return index, nil
}

func (taggerOp *taggerSQLite) Close() error {
	return errors.Wrap(taggerOp.db.Close(), "on taggerSQLite.Close()")
}

func (taggerOp *taggerSQLite) Clean(*selectors.Term) error {
	_, err := taggerOp.db.Exec("DELETE FROM " + taggerOp.table)

	return err
}
