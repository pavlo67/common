package tagger_pg

import (
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/pkg/errors"

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

var fieldsToCountTags = []string{"tag", "is_internal", "parted_size"}
var fieldsToCountTagsStr = strings.Join(fieldsToCountTags, ", ")
var fieldsToUpdateTagsStr = sqllib_pg.WildcardsForUpdate(fieldsToCountTags)

var fieldsToSave = []string{"joiner_key", "id", "tag", "relation"}
var fieldsToInsertStr = strings.Join(fieldsToSave, ", ")
var fieldsToUpdateStr = sqllib_pg.WildcardsForUpdate(fieldsToSave)

var _ tagger.Operator = &tagsPg{}

type tagsPg struct {
	db    *sql.DB
	table string

	ownInterfaceKey joiner.InterfaceKey

	sqlList, sqlIndexTagged, sqlIndexTaggedAll, sqlCountJoinerKeys, sqlCountTags, sqlCountTagsAll, sqlTagPartedSize string
	stmList, stmIndexTagged, stmIndexTaggedAll, stmCountJoinerKeys, stmCountTags, stmCountTagsAll, stmTagPartedSize *sql.Stmt

	// sqlSetTag, sqlGetTag
	sqlAddTag, sqlAddTagged, sqlRemoveTagged string
}

const onNew = "on tagsPg.New(): "

func New(access config.Access, ownInterfaceKey joiner.InterfaceKey) (tagger.Operator, crud.Cleaner, error) {
	db, err := sqllib_pg.Connect(access)
	if err != nil {
		return nil, nil, errors.Wrap(err, onNew)
	}

	tableTagged := tableTaggedDefault
	tableTags := tableTagsDefault
	tableJoined := tableTagged + " LEFT JOIN " + tableTags + " ON joiner_key = '" + joinerKeyTags + "' AND " + tableTagged + ".id  = " + tableTags + ".tag"

	taggerOp := tagsPg{
		db:    db,
		table: tableTagged,

		ownInterfaceKey: ownInterfaceKey,

		sqlAddTagged: "INSERT INTO " + tableTagged + " (" + fieldsToInsertStr + ") VALUES (" + sqllib_pg.WildcardsForInsert(fieldsToSave) + ")" +
			" ON CONFLICT (joiner_key, id, tag) DO UPDATE SET " + fieldsToUpdateStr,

		sqlRemoveTagged:   "DELETE                          FROM " + tableTagged + " WHERE joiner_key = $1 AND id = $2",
		sqlList:           "SELECT tag, relation            FROM " + tableTagged + " WHERE joiner_key = $1 AND id = $2 ORDER BY tag",
		sqlIndexTagged:    "SELECT joiner_key, id, relation FROM " + tableTagged + " WHERE joiner_key = $1 AND tag = $2",
		sqlIndexTaggedAll: "SELECT joiner_key, id, relation FROM " + tableTagged + " WHERE                     tag = $1                                   ORDER BY joiner_key",

		sqlAddTag: "INSERT INTO " + tableTags + " (" + fieldsToCountTagsStr + ") VALUES (" + sqllib_pg.WildcardsForInsert(fieldsToCountTags) + ")" +
			" ON CONFLICT (tag) DO UPDATE SET " + fieldsToUpdateTagsStr,

		sqlTagPartedSize: "SELECT SUM(parted_size) FROM " + tableJoined + " WHERE " + tableTagged + ".tag = $1",

		sqlCountJoinerKeys: "SELECT joiner_key,              COUNT(*) AS cnt, SUM(parted_size) FROM " + tableJoined + " WHERE  " + tableTagged + ".tag = $1 GROUP BY joiner_key              ORDER BY cnt DESC",
		sqlCountTags:       "SELECT " + tableTagged + ".tag, COUNT(*) AS cnt, SUM(parted_size) FROM " + tableJoined + " WHERE joiner_key = $1               GROUP BY " + tableTagged + ".tag ORDER BY cnt DESC",
		sqlCountTagsAll:    "SELECT " + tableTagged + ".tag, COUNT(*) AS cnt, SUM(parted_size) FROM " + tableJoined + "                                     GROUP BY " + tableTagged + ".tag ORDER BY cnt DESC",
	}

	sqlStmts := []sqllib.SqlStmt{
		{Stmt: &taggerOp.stmList, Sql: taggerOp.sqlList},
		{Stmt: &taggerOp.stmIndexTagged, Sql: taggerOp.sqlIndexTagged},
		{Stmt: &taggerOp.stmIndexTaggedAll, Sql: taggerOp.sqlIndexTaggedAll},
		{Stmt: &taggerOp.stmCountJoinerKeys, Sql: taggerOp.sqlCountJoinerKeys},
		{Stmt: &taggerOp.stmCountTags, Sql: taggerOp.sqlCountTags},
		{Stmt: &taggerOp.stmCountTagsAll, Sql: taggerOp.sqlCountTagsAll},
		{Stmt: &taggerOp.stmTagPartedSize, Sql: taggerOp.sqlTagPartedSize},
	}

	for _, sqlStmt := range sqlStmts {
		if err := sqllib.Prepare(db, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, nil, errors.Wrap(err, onNew)
		}
	}

	return &taggerOp, &taggerOp, nil
}

const onAddTags = "on tagsPg.AddTags(): "

func (taggerOp *tagsPg) AddTags(toTag joiner.Link, items []tagger.Tag, _ *crud.SaveOptions) error {
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
		var params []byte
		if len(tag.Params) > 0 {
			params, err = json.Marshal(tag.Params)
			if err != nil {
				return errors.Wrapf(err, onAddTags+"can't marshal .Params(%#v)", tag)
			}
		}

		values := []interface{}{toTag.InterfaceKey, toTag.ID, tag.Label, params}

		if _, err = stmAddTags.Exec(values...); err != nil {
			err = errors.Wrapf(err, onAddTags+": on stmAddTags(%s).Exec(%#v)", taggerOp.sqlAddTagged, values)
			goto ROLLBACK
		}

		// add = append(add, tag.Label)
	}

	if err = taggerOp.countTagChanged(toTag, nil, tx); err != nil {
		err = errors.Wrap(err, onAddTags)
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

const onReplaceTags = "on tagsPg.ReplaceTags(): "

func (taggerOp *tagsPg) ReplaceTags(toTag joiner.Link, items []tagger.Tag, options *crud.SaveOptions) error {

	var tagsFiltered []tagger.Tag
	for _, tag := range items {
		tag.Label = strings.TrimSpace(tag.Label)
		if tag.Label != "" {
			tagsFiltered = append(tagsFiltered, tag)
		}
	}

	tagsOld, err := taggerOp.ListTags(toTag, nil) // TODO!!! use correct options
	if err != nil {
		return errors.Wrap(err, onReplaceTags)
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

	values := []interface{}{toTag.InterfaceKey, toTag.ID}
	_, err = tx.Exec(taggerOp.sqlRemoveTagged, values...)
	if err != nil {
		err = errors.Wrapf(err, onReplaceTags+": on tx.Exec(%s, %#v)", taggerOp.sqlRemoveTagged, values)
		goto ROLLBACK
	}

	for _, tag := range tagsFiltered {

		var params []byte
		if len(tag.Params) > 0 {
			params, err = json.Marshal(tag.Params)
			if err != nil {
				return errors.Wrapf(err, onAddTags+"can't marshal .Params(%#v)", tag)
			}
		}

		values := []interface{}{toTag.InterfaceKey, toTag.ID, tag.Label, params}

		_, err = stmAddTags.Exec(values...)
		if err != nil {
			err = errors.Wrapf(err, onReplaceTags+": on tx.Exec(%s, %#v)", taggerOp.sqlAddTagged, values)
			goto ROLLBACK
		}
	}

	if err = taggerOp.countTagChanged(toTag, tagLabelsRemoved, tx); err != nil {
		err = errors.Wrap(err, onReplaceTags)
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

const onListTags = "on tagsPg.ListTags(): "

func (taggerOp *tagsPg) ListTags(toTag joiner.Link, _ *crud.GetOptions) ([]tagger.Tag, error) {
	values := []interface{}{toTag.InterfaceKey, toTag.ID}

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
		var params []byte

		err = rows.Scan(&tag.Label, &params)
		if err != nil {
			return items, errors.Wrapf(err, onListTags+sqllib.CantScanQueryRow, taggerOp.sqlList, values)
		}

		if len(params) > 0 {
			if err = json.Unmarshal(params, &tag.Params); err != nil {
				return items, errors.Wrapf(err, onListTags+"can't unmarshal .Params (%s)", params)
			}
		}

		items = append(items, tag)
	}
	err = rows.Err()
	if err != nil {
		return items, errors.Wrapf(err, onListTags+": "+sqllib.RowsError, taggerOp.sqlList, values)
	}

	return items, nil
}

const onCountTags = "on tagsPg.CountTags(): "

func (taggerOp *tagsPg) CountTags(key *joiner.InterfaceKey, _ *crud.GetOptions) ([]tagger.TagCount, error) {
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

	counter := []tagger.TagCount{} // to return [] to front-end instead null

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

const onIndexTagged = "on tagsPg.IndexTagged()"

func (taggerOp *tagsPg) IndexTagged(key *joiner.InterfaceKey, label string, _ *crud.GetOptions) (tagger.Index, error) {
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
		var params []byte

		err = rows.Scan(&key, &tagged.ID, &params)
		if err != nil {
			return index, errors.Wrapf(err, onIndexTagged+sqllib.CantScanQueryRow, query, values)
		}

		if len(params) > 0 {
			if err = json.Unmarshal(params, &tagged.Params); err != nil {
				return index, errors.Wrapf(err, onIndexTagged+"can't unmarshal .Params (%s)", params)
			}
		}

		index[joiner.InterfaceKey(key)] = append(index[joiner.InterfaceKey(key)], tagged)
	}
	err = rows.Err()
	if err != nil {
		return index, errors.Wrapf(err, onIndexTagged+": "+sqllib.RowsError, query, values)
	}

	return index, nil
}

func (taggerOp *tagsPg) Close() error {
	return errors.Wrap(taggerOp.db.Close(), "on tagsPg.Close()")
}
