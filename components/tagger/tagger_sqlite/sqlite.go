package tagger_sqlite

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/libraries/sqllib"
	"github.com/pavlo67/workshop/common/libraries/sqllib/sqllib_sqlite"
	"github.com/pavlo67/workshop/components/crud"
	"github.com/pavlo67/workshop/components/tagger"
)

const limitDefault = 200
const tableDefault = "tags"

var fieldsToSave = []string{"key", "id", "tag"}
var fieldsToSaveStr = strings.Join(fieldsToSave, ", ")

var _ tagger.Operator = &taggerSQLite{}

type taggerSQLite struct {
	limit int
	table string
	db    *sql.DB

	sqlSave, sqlRemove, sqlReset, sqlTags, sqlListTagged string
	stmSave, stmRemove, stmReset, stmTags, stmListTagged *sql.Stmt
}

const onNew = "on taggerSQLite.New(): "

func NewTagger(access config.Access, table string, limit int) (tagger.Operator, error) {
	db, err := sqllib_sqlite.Connect(access)
	if err != nil {
		return nil, errors.Wrap(err, onNew)
	}

	if table == "" {
		table = tableDefault
	}
	if limit <= 0 {
		limit = limitDefault
	}

	taggerOp := taggerSQLite{
		db:    db,
		limit: limit,
		table: table,

		sqlSave:   "INSERT INTO " + table + " (" + fieldsToSaveStr + ") VALUES (" + strings.Repeat(",? ", len(fieldsToSave))[1:] + ")",
		sqlRemove: "DELETE FROM " + table + " where key = ? AND id = ? and tag = ?",
		sqlReset:  "DELETE FROM " + table + " where key = ? AND id = ?",

		sqlTags:       "SELECT tag FROM " + table + " ORDER BY tag WHERE key = ? AND id = ?",
		sqlListTagged: "SELECT key, id FROM " + table + " ORDER BY key, id where tag = ?",
	}

	sqlStmts := []sqllib.SqlStmt{
		{&taggerOp.stmSave, taggerOp.sqlSave},
		{&taggerOp.stmRemove, taggerOp.sqlRemove},
		{&taggerOp.stmReset, taggerOp.sqlReset},

		{&taggerOp.stmTags, taggerOp.sqlTags},
		{&taggerOp.stmListTagged, taggerOp.sqlListTagged},
	}

	for _, sqlStmt := range sqlStmts {
		if err := sqllib.Prepare(db, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, errors.Wrap(err, onNew)
		}
	}

	return &taggerOp, nil
}

const onTags = "on taggerSQLite.Tags(): "

func (taggerOp *taggerSQLite) Tags(key joiner.InterfaceKey, id common.ID, _ *crud.GetOptions) ([]tagger.Tag, error) {
	if len(id) < 1 {
		return nil, errors.New(onTags + "empty ID")
	}

	idNum, err := strconv.ParseUint(string(id), 10, 64)
	if err != nil {
		return nil, errors.Errorf(onTags+"wrong ID (%s)", id)
	}

	rows, err := taggerOp.stmTags.Query(idNum)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, onTags+sqllib.CantQuery, taggerOp.sqlTags, idNum)
	}
	defer rows.Close()

	var tags []tagger.Tag

	for rows.Next() {
		var tag string

		err = rows.Scan(&tag)
		if err != nil {
			return tags, errors.Wrapf(err, onTags+sqllib.CantScanQueryRow, taggerOp.sqlTags, idNum)
		}

		tags = append(tags, tagger.Tag(tag))
	}
	err = rows.Err()
	if err != nil {
		return tags, errors.Wrapf(err, onTags+": "+sqllib.RowsError, taggerOp.sqlTags, idNum)
	}

	return tags, nil
}

const onSave = "on taggerSQLite.Save(): "

func (taggerOp *taggerSQLite) Save(key joiner.InterfaceKey, id common.ID, tags []tagger.Tag, _ *crud.SaveOptions) error {
	for _, tag := range tags {
		values := []interface{}{key, id, tag}

		_, err := taggerOp.stmSave.Exec(values...)
		if err != nil {
			return errors.Wrapf(err, onSave+sqllib.CantExec, taggerOp.sqlSave, values)
		}
	}

	return nil
}

const onReplace = "on taggerSQLite.Replace(): "

func (taggerOp *taggerSQLite) Replace(key joiner.InterfaceKey, id common.ID, tags []tagger.Tag, options *crud.SaveOptions) error {
	values := []interface{}{key, id}

	_, err := taggerOp.stmReset.Exec(values...)
	if err != nil {
		return errors.Wrapf(err, onSave+sqllib.CantExec, taggerOp.sqlReset, values)
	}

	err = taggerOp.Save(key, id, tags, options)
	if err != nil {
		return errors.Wrap(err, onReplace)
	}

	return nil
}

const onRemove = "on taggerSQLite.Remove()"

func (taggerOp *taggerSQLite) Remove(key joiner.InterfaceKey, id common.ID, tags []tagger.Tag, _ *crud.SaveOptions) error {
	for _, tag := range tags {
		values := []interface{}{key, id, tag}

		_, err := taggerOp.stmRemove.Exec(values...)
		if err != nil {
			return errors.Wrapf(err, onRemove+sqllib.CantExec, taggerOp.sqlRemove, values)
		}
	}

	return nil
}

const onListTagged = "on taggerSQLite.ListTagged()"

func (taggerOp *taggerSQLite) ListTagged(tag tagger.Tag, _ *crud.GetOptions) ([]tagger.Tagged, error) {
	values := []interface{}{tag}

	rows, err := taggerOp.stmListTagged.Query(values...)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, onListTagged+sqllib.CantQuery, taggerOp.sqlListTagged, values)
	}
	defer rows.Close()

	var taggedAll []tagger.Tagged

	for rows.Next() {
		var key, id string
		err = rows.Scan(&key, &id)
		if err != nil {
			return taggedAll, errors.Wrapf(err, onListTagged+sqllib.CantScanQueryRow, taggerOp.sqlListTagged, values)
		}

		taggedAll = append(taggedAll, tagger.Tagged{
			InterfaceKey: joiner.InterfaceKey(key),
			ID:           common.ID(id),
		})
	}
	err = rows.Err()
	if err != nil {
		return taggedAll, errors.Wrapf(err, onListTagged+": "+sqllib.RowsError, taggerOp.sqlListTagged, values)
	}

	return taggedAll, nil
}

func (taggerOp *taggerSQLite) Close() error {
	return errors.Wrap(taggerOp.db.Close(), "on taggerSQLite.Close()")
}

func (taggerOp *taggerSQLite) Clean() error {
	_, err := taggerOp.db.Exec("TRUNCATE " + taggerOp.table)

	return err
}
