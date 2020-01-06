package tagger_pg

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/libraries/sqllib"
	"github.com/pavlo67/workshop/common/libraries/sqllib/sqllib_pg"

	"github.com/pavlo67/workshop/components/tagger"
)

const tableTaggedDefault = "tagged"
const tableTagsDefault = "tags"

const joinerKeyTags = "tags"

var fieldsToCount = []string{"tag", "is_internal", "parted_size"}
var fieldsToCountStr = strings.Join(fieldsToCount, ", ")

var fieldsToSave = []string{"joiner_key", "id", "tag", "relation"}
var fieldsToSaveStr = strings.Join(fieldsToSave, ", ")

var _ tagger.Operator = &tagsSQLite{}

type tagsSQLite struct {
	db    *sql.DB
	table string

	ownInterfaceKey joiner.InterfaceKey

	sqlList, sqlIndexTagged, sqlIndexTaggedAll, sqlCountJoinerKeys, sqlCountTags, sqlCountTagsAll string
	stmList, stmIndexTagged, stmIndexTaggedAll, stmCountJoinerKeys, stmCountTags, stmCountTagsAll *sql.Stmt

	// sqlSetTag, sqlGetTag
	sqlAddTag, sqlCountTagFull, sqlAddTagged, sqlRemoveTagged string
}

const onNew = "on tagsSQLite.New(): "

func New(access config.Access, ownInterfaceKey joiner.InterfaceKey) (tagger.Operator, crud.Cleaner, error) {
	db, err := sqllib_pg.Connect(access)
	if err != nil {
		return nil, nil, errors.Wrap(err, onNew)
	}

	tableTagged := tableTaggedDefault
	tableTags := tableTagsDefault
	tableJoined := tableTagged + " LEFT JOIN " + tableTags + " ON joiner_key = '" + joinerKeyTags + "' AND " + tableTagged + ".id  = " + tableTags + ".tag"

	taggerOp := tagsSQLite{
		db:    db,
		table: tableTagged,

		ownInterfaceKey: ownInterfaceKey,

		// TODO: on conflict REPLACE
		sqlAddTagged: "INSERT INTO " + tableTagged + " (" + fieldsToSaveStr + ") VALUES (" + sqllib_pg.WildcardsForInsert(fieldsToSave) + ")",

		sqlRemoveTagged:   "DELETE                          FROM " + tableTagged + " WHERE joiner_key = $1 AND id = $2",
		sqlList:           "SELECT tag, relation            FROM " + tableTagged + " WHERE joiner_key = $1 AND id = $2 ORDER BY tag",
		sqlIndexTagged:    "SELECT joiner_key, id, relation FROM " + tableTagged + " WHERE joiner_key = $1 AND tag = $2",
		sqlIndexTaggedAll: "SELECT joiner_key, id, relation FROM " + tableTagged + " WHERE                     tag = $1                                   ORDER BY joiner_key",

		// TODO: on conflict REPLACE
		sqlAddTag:       "INSERT INTO " + tableTags + " (" + fieldsToCountStr + ") VALUES (" + sqllib_pg.WildcardsForInsert(fieldsToCount) + ")",
		sqlCountTagFull: "SELECT SUM(parted_size) FROM " + tableJoined + " WHERE " + tableTagged + ".tag = $1",

		sqlCountJoinerKeys: "SELECT joiner_key, COUNT(*) AS cnt FROM " + tableTagged + " WHERE tag = $1        GROUP BY joiner_key ORDER BY cnt DESC",
		sqlCountTags:       "SELECT tag,        COUNT(*) AS cnt FROM " + tableTagged + " WHERE joiner_key = $1 GROUP BY tag        ORDER BY cnt DESC",
		sqlCountTagsAll:    "SELECT tag,        COUNT(*) AS cnt FROM " + tableTagged + "                       GROUP BY tag        ORDER BY cnt DESC",
	}

	sqlStmts := []sqllib.SqlStmt{
		{&taggerOp.stmList, taggerOp.sqlList},
		{&taggerOp.stmIndexTagged, taggerOp.sqlIndexTagged},
		{&taggerOp.stmIndexTaggedAll, taggerOp.sqlIndexTaggedAll},
		{&taggerOp.stmCountJoinerKeys, taggerOp.sqlCountJoinerKeys},
		{&taggerOp.stmCountTags, taggerOp.sqlCountTags},
		{&taggerOp.stmCountTagsAll, taggerOp.sqlCountTagsAll},
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
		query = taggerOp.sqlCountTagsAll
		stm = taggerOp.stmCountTagsAll
	} else {
		values = []interface{}{*key}
		query = taggerOp.sqlCountTags
		stm = taggerOp.stmCountTags
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
