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

const tableDefault = "tagged"
const tableTagsDefault = "tags"

var fieldsToCount = []string{"tag", "is_internal", "parted_size"}
var fieldsToCountStr = strings.Join(fieldsToCount, ", ")

var fieldsToSave = []string{"key", "id", "tag", "relation"}
var fieldsToSaveStr = strings.Join(fieldsToSave, ", ")

var _ tagger.Operator = &taggerSQLite{}
var _ crud.Cleaner = &taggerSQLite{}

type taggerSQLite struct {
	db          *sql.DB
	table       string
	tableJoined string

	ownInterfaceKey joiner.InterfaceKey

	sqlListTags, sqlIndexWithTag, sqlCountTagged, sqlCountTaggedAll string
	stmListTags, stmIndexWithTag, stmCountTagged, stmCountTaggedAll *sql.Stmt

	// sqlSetTag, sqlGetTag
	sqlAddTag, sqlCountTag, sqlAddTagged, sqlRemoveTagged string
}

const onNew = "on taggerSQLite.New(): "

func New(access config.Access, ownInterfaceKey joiner.InterfaceKey) (tagger.Operator, crud.Cleaner, error) {
	db, err := sqllib_sqlite.Connect(access)
	if err != nil {
		return nil, nil, errors.Wrap(err, onNew)
	}

	table := tableDefault
	tableTags := tableTagsDefault
	tableJoined := table + "   LEFT JOIN " + tableTags + " ON " + table + ".tag = " + tableTags + ".tag"
	tableJoinedUp := table + " LEFT JOIN " + tableTags + " ON " + table + ".id  = " + tableTags + ".tag AND key = ?"

	taggerOp := taggerSQLite{
		db:          db,
		table:       table,
		tableJoined: tableJoined,

		ownInterfaceKey: ownInterfaceKey,

		sqlAddTagged:    "INSERT OR REPLACE INTO " + table + " (" + fieldsToSaveStr + ") VALUES (" + strings.Repeat(",? ", len(fieldsToSave))[1:] + ")",
		sqlRemoveTagged: "DELETE FROM " + table + " WHERE key = ? AND id = ?",

		// sqlGetTag: "SELECT " + fieldsToCountStr + " FROM " + tableTags + " WHERE tag = ?",
		// sqlSetTag: "UPDATE " + tableTags + " SET is_internal = ?, parted_size = ? WHERE tag = ?",

		sqlAddTag: "INSERT OR REPLACE INTO " + tableTags + " (" + fieldsToCountStr + ") VALUES (" + strings.Repeat(",? ", len(fieldsToCount))[1:] + ")",

		// CASE parted_size WHEN NULL THEN 0 ELSE parted_size END
		sqlCountTag: "SELECT SUM(parted_size) FROM " + tableJoinedUp + " WHERE " + table + ".tag = ?",
		sqlListTags: "SELECT tag, relation    FROM " + table + "         WHERE key = ? AND id = ?    ORDER BY tag",

		sqlIndexWithTag:   "SELECT key, id, relation                        FROM " + table + "       WHERE tag = ?                                       ORDER BY key, id",
		sqlCountTagged:    "SELECT " + table + ".tag, COUNT(*), parted_size FROM " + tableJoined + " WHERE key = ?            GROUP BY " + table + ".tag ORDER BY " + table + ".tag",
		sqlCountTaggedAll: "SELECT " + table + ".tag, COUNT(*), parted_size FROM " + tableJoined + "                          GROUP BY " + table + ".tag ORDER BY " + table + ".tag",
	}

	sqlStmts := []sqllib.SqlStmt{
		{&taggerOp.stmListTags, taggerOp.sqlListTags},
		{&taggerOp.stmIndexWithTag, taggerOp.sqlIndexWithTag},
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

const onAddTags = "on taggerSQLite.AddTags(): "

func (taggerOp *taggerSQLite) AddTags(key joiner.InterfaceKey, id common.ID, tags []tagger.Tag, _ *crud.SaveOptions) error {
	tx, err := taggerOp.db.Begin()
	if err != nil {
		return errors.Wrap(err, onAddTags+": on taggerOp.db.Begin()")
	}

	stmAddTags, err := tx.Prepare(taggerOp.sqlAddTagged)
	if err != nil {
		return errors.Wrapf(err, onAddTags+": on tx.Prepare(%s)", taggerOp.sqlAddTagged)
	}

	// var add []string

	for _, tag := range tags {
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

const onReplaceTags = "on taggerSQLite.ReplaceTags(): "

func (taggerOp *taggerSQLite) ReplaceTags(key joiner.InterfaceKey, id common.ID, tags []tagger.Tag, options *crud.SaveOptions) error {

	tagsOld, err := taggerOp.ListTags(key, id, nil) // TODO!!! use correct options
	if err != nil {
		return errors.Wrapf(err, onReplaceTags+": on taggerOp.ListTags(%s, %s, nil)", key, id)
	}
	var tagLabelsRemoved []string
	for _, tag := range tagsOld {
		tagLabelsRemoved = append(tagLabelsRemoved, tag.Label)
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

	for _, tag := range tags {
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
		var tag tagger.Tag

		err = rows.Scan(&tag.Label, &tag.Relation)
		if err != nil {
			return tags, errors.Wrapf(err, onListTags+sqllib.CantScanQueryRow, taggerOp.sqlListTags, values)
		}

		tags = append(tags, tag)
	}
	err = rows.Err()
	if err != nil {
		return tags, errors.Wrapf(err, onListTags+": "+sqllib.RowsError, taggerOp.sqlListTags, values)
	}

	return tags, nil
}

const onCountTagged = "on taggerSQLite.CountTagged(): "

func (taggerOp *taggerSQLite) CountTagged(key *joiner.InterfaceKey, _ *crud.GetOptions) (tagger.Counter, error) {
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

	counter := tagger.Counter{}

	for rows.Next() {
		var key string
		var count tagger.TaggedCount

		err = rows.Scan(&key, &count.Immediate, &count.Full)
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

func (taggerOp *taggerSQLite) IndexWithTag(label string, _ *crud.GetOptions) (tagger.Index, error) {
	values := []interface{}{label}

	rows, err := taggerOp.stmIndexWithTag.Query(values...)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, onIndexWithTag+sqllib.CantQuery, taggerOp.sqlIndexWithTag, values)
	}
	defer rows.Close()

	index := tagger.Index{}

	for rows.Next() {
		var key string
		var tagged tagger.Tagged

		err = rows.Scan(&key, &tagged.ID, &tagged.Relation)
		if err != nil {
			return index, errors.Wrapf(err, onIndexWithTag+sqllib.CantScanQueryRow, taggerOp.sqlIndexWithTag, values)
		}

		index[joiner.InterfaceKey(key)] = append(index[joiner.InterfaceKey(key)], tagged)
	}
	err = rows.Err()
	if err != nil {
		return index, errors.Wrapf(err, onIndexWithTag+": "+sqllib.RowsError, taggerOp.sqlIndexWithTag, values)
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
