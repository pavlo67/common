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

	"github.com/pavlo67/workshop/components/tagger"
)

const tableDefault = "tagged"
const tableTagsDefault = "tags"

var fieldsToCount = []string{"tag", "is_internal", "parted_size"}
var fieldsToCountStr = strings.Join(fieldsToCount, ", ")

var fieldsToSave = []string{"key", "id", "tag", "relation"}
var fieldsToSaveStr = strings.Join(fieldsToSave, ", ")

var _ tagger.Operator = &tagsSQLite{}

type tagsSQLite struct {
	db          *sql.DB
	table       string
	tableJoined string

	ownInterfaceKey joiner.InterfaceKey

	sqlList, sqlIndexTagged, sqlIndexTaggedAll, sqlCount, sqlCountAll string
	stmList, stmIndexTagged, stmIndexTaggedAll, stmCount, stmCountAll *sql.Stmt

	// sqlSetTag, sqlGetTag
	sqlAddTag, sqlCountTag, sqlAddTagged, sqlRemoveTagged string
}

const onNew = "on tagsSQLite.New(): "

func New(access config.Access, ownInterfaceKey joiner.InterfaceKey) (tagger.Operator, crud.Cleaner, error) {
	db, err := sqllib_sqlite.Connect(access)
	if err != nil {
		return nil, nil, errors.Wrap(err, onNew)
	}

	table := tableDefault
	tableTags := tableTagsDefault
	tableJoined := table + "   LEFT JOIN " + tableTags + " ON " + table + ".tag = " + tableTags + ".tag"
	tableJoinedUp := table + " LEFT JOIN " + tableTags + " ON " + table + ".id  = " + tableTags + ".tag AND key = ?"

	taggerOp := tagsSQLite{
		db:          db,
		table:       table,
		tableJoined: tableJoined,

		ownInterfaceKey: ownInterfaceKey,

		sqlAddTagged:    "INSERT OR REPLACE INTO " + table + " (" + fieldsToSaveStr + ") VALUES (" + strings.Repeat(",? ", len(fieldsToSave))[1:] + ")",
		sqlRemoveTagged: "DELETE FROM " + table + " WHERE key = ? AND id = ?",

		sqlAddTag:   "INSERT OR REPLACE INTO " + tableTags + " (" + fieldsToCountStr + ") VALUES (" + strings.Repeat(",? ", len(fieldsToCount))[1:] + ")",
		sqlCountTag: "SELECT SUM(parted_size) FROM " + tableJoinedUp + " WHERE " + table + ".tag = ?",

		sqlList:           "SELECT tag, relation    FROM " + table + "         WHERE key = ? AND id = ?    ORDER BY tag",
		sqlIndexTagged:    "SELECT key, id, relation                        FROM " + table + "       WHERE key = ? AND tag = ?",
		sqlIndexTaggedAll: "SELECT key, id, relation                        FROM " + table + "       WHERE tag = ?                            ORDER BY key",
		sqlCount:          "SELECT " + table + ".tag, COUNT(*), parted_size FROM " + tableJoined + " WHERE key = ? GROUP BY " + table + ".tag ORDER BY " + table + ".tag",
		sqlCountAll:       "SELECT " + table + ".tag, COUNT(*), parted_size FROM " + tableJoined + "               GROUP BY " + table + ".tag ORDER BY " + table + ".tag",
	}

	sqlStmts := []sqllib.SqlStmt{
		{&taggerOp.stmList, taggerOp.sqlList},
		{&taggerOp.stmIndexTagged, taggerOp.sqlIndexTagged},
		{&taggerOp.stmIndexTaggedAll, taggerOp.sqlIndexTaggedAll},
		{&taggerOp.stmCount, taggerOp.sqlCount},
		{&taggerOp.stmCountAll, taggerOp.sqlCountAll},
	}

	for _, sqlStmt := range sqlStmts {
		if err := sqllib.Prepare(db, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, nil, errors.Wrap(err, onNew)
		}
	}

	return &taggerOp, &taggerOp, nil
}

const onAddTags = "on tagsSQLite.AddTags(): "

func (taggerOp *tagsSQLite) AddTags(key joiner.InterfaceKey, id common.ID, items []tagger.Tag, _ *crud.SaveOptions) error {
	var tagsFiltered []tagger.Tag
	for _, tag := range items {
		tag.Label = strings.TrimSpace(tag.Label)
		if tag.Label != "" {
			tagsFiltered = append(tagsFiltered, tag)
		}
	}

	if len(tagsFiltered) < 1 {
		return nil
	}

	tx, err := taggerOp.db.Begin()
	if err != nil {
		return errors.Wrap(err, onAddTags+": on taggerOp.db.Begin()")
	}

	stmAddTags, err := tx.Prepare(taggerOp.sqlAddTagged)
	if err != nil {
		return errors.Wrapf(err, onAddTags+": on tx.Prepare(%s)", taggerOp.sqlAddTagged)
	}

	// var add []string

	for _, tag := range tagsFiltered {
		values := []interface{}{key, id, tag.Label, tag.Relation}

		if _, err = stmAddTags.Exec(values...); err != nil {
			err = errors.Wrapf(err, onAddTags+": on stmAddTags(%s).Exec(%#v)", taggerOp.sqlAddTagged, values)
			goto ROLLBACK
		}

		// add = append(add, tag.Label)
	}

	if err = taggerOp.countTagChanged(key, id, nil, tx); err != nil {
		err = errors.Wrapf(err, onAddTags+": on taggerOp.countTagChanged(%s, %s, nil, tx)", key, id)
		goto ROLLBACK
	}

	if err = tx.Commit(); err == nil {
		return nil
	}
	err = errors.Wrap(err, onAddTags+": on tx.Commit()")

ROLLBACK:
	if errRollback := tx.Rollback(); errRollback != nil {
		return errors.Wrapf(err, onAddTags+": on tx.Rollback(): %s", errRollback)
	}

	return err
}

const onReplaceTags = "on tagsSQLite.ReplaceTags(): "

func (taggerOp *tagsSQLite) ReplaceTags(key joiner.InterfaceKey, id common.ID, items []tagger.Tag, options *crud.SaveOptions) error {

	var tagsFiltered []tagger.Tag
	for _, tag := range items {
		tag.Label = strings.TrimSpace(tag.Label)
		if tag.Label != "" {
			tagsFiltered = append(tagsFiltered, tag)
		}
	}

	tagsOld, err := taggerOp.ListTags(key, id, nil) // TODO!!! use correct options
	if err != nil {
		return errors.Wrapf(err, onReplaceTags+": on taggerOp.ListTags(%s, %s, nil)", key, id)
	}
	var tagLabelsRemoved []string
	for _, tag := range tagsOld {
		tagLabelsRemoved = append(tagLabelsRemoved, strings.TrimSpace(tag.Label))
	}

	if len(tagsFiltered) < 1 && len(tagLabelsRemoved) < 1 {
		return nil
	}

	tx, err := taggerOp.db.Begin()
	if err != nil {
		return errors.Wrap(err, onReplaceTags+": on taggerOp.db.Begin()")
	}

	stmAddTags, err := tx.Prepare(taggerOp.sqlAddTagged)
	if err != nil {
		return errors.Wrapf(err, onReplaceTags+": on tx.Prepare(%s)", taggerOp.sqlAddTagged)
	}

	//var add []string

	values := []interface{}{key, id}
	_, err = tx.Exec(taggerOp.sqlRemoveTagged, values...)
	if err != nil {
		err = errors.Wrapf(err, onReplaceTags+": on tx.Exec(%s, %#v)", taggerOp.sqlRemoveTagged, values)
		goto ROLLBACK
	}

	for _, tag := range tagsFiltered {
		values := []interface{}{key, id, tag.Label, tag.Relation}

		_, err = stmAddTags.Exec(values...)
		if err != nil {
			err = errors.Wrapf(err, onReplaceTags+": on tx.Exec(%s, %#v)", taggerOp.sqlAddTagged, values)
			goto ROLLBACK
		}
	}

	if err = taggerOp.countTagChanged(key, id, tagLabelsRemoved, tx); err != nil {
		err = errors.Wrapf(err, onReplaceTags+": on taggerOp.countTagChanged(%s, %s, %#v, tx)", key, id, tagLabelsRemoved)
		goto ROLLBACK
	}

	err = tx.Commit()
	if err == nil {
		return nil
	}
	err = errors.Wrap(err, onReplaceTags+": on tx.Commit()")

ROLLBACK:
	errRollback := tx.Rollback()
	if errRollback != nil {
		return errors.Wrapf(err, onReplaceTags+": on tx.Rollback(): %s", errRollback)
	}
	return err
}

const onListTags = "on tagsSQLite.ListTags(): "

func (taggerOp *tagsSQLite) ListTags(key joiner.InterfaceKey, id common.ID, _ *crud.GetOptions) ([]tagger.Tag, error) {
	values := []interface{}{key, id}

	rows, err := taggerOp.stmList.Query(values...)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, onListTags+sqllib.CantQuery, taggerOp.sqlList, values)
	}
	defer rows.Close()

	var items []tagger.Tag

	for rows.Next() {
		var tag tagger.Tag

		err = rows.Scan(&tag.Label, &tag.Relation)
		if err != nil {
			return items, errors.Wrapf(err, onListTags+sqllib.CantScanQueryRow, taggerOp.sqlList, values)
		}

		items = append(items, tag)
	}
	err = rows.Err()
	if err != nil {
		return items, errors.Wrapf(err, onListTags+": "+sqllib.RowsError, taggerOp.sqlList, values)
	}

	return items, nil
}

const onCountTags = "on tagsSQLite.CountTags(): "

func (taggerOp *tagsSQLite) CountTags(key *joiner.InterfaceKey, _ *crud.GetOptions) ([]tagger.TagCount, error) {
	var values []interface{}
	var query string
	var stm *sql.Stmt

	if key == nil {
		query = taggerOp.sqlCountAll
		stm = taggerOp.stmCountAll
	} else {
		values = []interface{}{*key}
		query = taggerOp.sqlCount
		stm = taggerOp.stmCount
	}

	rows, err := stm.Query(values...)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, onCountTags+sqllib.CantQuery, query, values)
	}
	defer rows.Close()

	var counter []tagger.TagCount

	for rows.Next() {
		var count tagger.TagCount
		full := new(uint64)

		err = rows.Scan(&count.Label, &count.Immediate, &full)
		if err != nil {
			return counter, errors.Wrapf(err, onCountTags+sqllib.CantScanQueryRow, query, values)
		}

		if full != nil {
			count.Full = *full
		}

		counter = append(counter, count)
	}
	err = rows.Err()
	if err != nil {
		return counter, errors.Wrapf(err, onCountTags+": "+sqllib.RowsError, query, values)
	}

	return counter, nil
}

const onIndexTagged = "on tagsSQLite.IndexTagged()"

func (taggerOp *tagsSQLite) IndexTagged(key *joiner.InterfaceKey, label string, _ *crud.GetOptions) (tagger.Index, error) {
	var values []interface{}
	var query string
	var stm *sql.Stmt

	if key == nil {
		values = []interface{}{label}
		stm = taggerOp.stmIndexTaggedAll
		query = taggerOp.sqlIndexTaggedAll

	} else {
		values = []interface{}{key, label}
		stm = taggerOp.stmIndexTagged
		query = taggerOp.sqlIndexTagged

	}

	rows, err := stm.Query(values...)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, onIndexTagged+sqllib.CantQuery, query, values)
	}
	defer rows.Close()

	index := tagger.Index{}

	for rows.Next() {
		var key string
		var tagged tagger.Tagged

		err = rows.Scan(&key, &tagged.ID, &tagged.Relation)
		if err != nil {
			return index, errors.Wrapf(err, onIndexTagged+sqllib.CantScanQueryRow, query, values)
		}

		index[joiner.InterfaceKey(key)] = append(index[joiner.InterfaceKey(key)], tagged)
	}
	err = rows.Err()
	if err != nil {
		return index, errors.Wrapf(err, onIndexTagged+": "+sqllib.RowsError, query, values)
	}

	return index, nil
}

func (taggerOp *tagsSQLite) Close() error {
	return errors.Wrap(taggerOp.db.Close(), "on tagsSQLite.Close()")
}
