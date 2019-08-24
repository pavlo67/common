package data_sqlite

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/basis/common"
	"github.com/pavlo67/workshop/basis/common/sqllib"
	"github.com/pavlo67/workshop/basis/crud"
	"github.com/pavlo67/workshop/basis/selectors"

	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/instruments/indexer"
	"github.com/pavlo67/workshop/components/marks"
)

const limitDefault = 200

var tabledata = "datas"

var fieldsToList = []string{"id", "source_time", "source_url", "types", "title", "summary", "tags"}
var fieldsToListStr = strings.Join(fieldsToRead, ", ")

var fieldsToRead = []string{"source_id", "source_time", "source_url", "types", "title", "summary", "details", "href", "embedded", "tags"}
var fieldsToReadStr = strings.Join(fieldsToRead, ", ")

var fieldsToSave = []string{"source_id", "source_time", "source_url", "types", "title", "summary", "details", "href", "embedded", "tags", "indexes", "source_key", "origin"}
var fieldsToSaveStr = strings.Join(fieldsToSave, ", ")

var _ data.Operator = &dataSQLite{}

type dataSQLite struct {
	limit int
	db    *sql.DB

	stmHas, stmRead, stmSave, stmRemove, stmList *sql.Stmt
	sqlHas, sqlRead, sqlSave, sqlRemove, sqlList string
}

const onNew = "on dataSQLite.New(): "

func New(db *sql.DB, limit int) (data.Operator, error) {
	if db == nil {
		return nil, errors.New(onNew + "no db")
	}

	if limit <= 0 {
		limit = limitDefault
	}

	dataOp := dataSQLite{
		db:    db,
		limit: limit,

		sqlHas:  "SELECT count(*) FROM " + tabledata + " WHERE source_id = ? AND source_key = ?",
		sqlRead: "SELECT " + fieldsToReadStr + " FROM " + tabledata + " WHERE id = ?",

		sqlSave:   "INSERT INTO " + tabledata + " (" + fieldsToSaveStr + ") VALUES (" + strings.Repeat(",? ", len(fieldsToSave))[1:] + ")",
		sqlRemove: "DELETE FROM " + tabledata + " where ID = ?",

		sqlList: "SELECT " + fieldsToListStr + " FROM " + tabledata + " ORDER BY saved_at DESC LIMIT " + strconv.Itoa(limit),
	}

	sqlStmts := []sqllib.SqlStmt{
		{&dataOp.stmHas, dataOp.sqlHas},
		{&dataOp.stmRead, dataOp.sqlRead},

		{&dataOp.stmSave, dataOp.sqlSave},
		{&dataOp.stmRemove, dataOp.sqlRemove},

		{&dataOp.stmList, dataOp.sqlList},
	}

	for _, sqlStmt := range sqlStmts {
		if err := sqllib.Prepare(db, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, errors.Wrap(err, onNew)
		}
	}

	return &dataOp, nil
}

const onHas = "on dataSQLite.Has(): "

func (dataOp *dataSQLite) Has(originKey data.OriginKey, _ *crud.GetOptions) (uint, error) {
	if len(originKey.SourceKey) < 1 { // || len(originKey.SourceID) < 1
		return 0, errors.New(onHas + "empty ID")
	}

	values := []interface{}{originKey.SourceID, originKey.SourceKey}

	var cnt uint
	err := dataOp.stmHas.QueryRow(values...).Scan(&cnt)
	if err != nil {
		return cnt, errors.Wrapf(err, onHas+sqllib.CantScanQueryRow, dataOp.sqlHas, values)
	}

	return cnt, nil
}

const onRead = "on dataSQLite.Read(): "

func (dataOp *dataSQLite) Read(idStr common.ID, _ *crud.GetOptions) (*data.Item, error) {
	if len(idStr) < 1 {
		return nil, errors.New(onRead + "empty ID")
	}

	id, err := strconv.ParseUint(string(idStr), 10, 64)
	if err != nil {
		return nil, errors.Errorf(onRead+"wrong ID (%s)", idStr)
	}

	var item data.Item
	var embedded, tags string

	err = dataOp.stmRead.QueryRow(id).Scan(&item.SourceID, &item.SourceTime, &item.SourceURL, &item.Type, &item.Title, &item.Summary, &item.Details, &item.Href, &embedded, &tags)
	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrapf(err, onRead+sqllib.CantScanQueryRow, dataOp.sqlRead, id)
	}

	item.Tags = strings.Split(tags, "\n")
	err = json.Unmarshal([]byte(embedded), &item.Embedded)
	if err != nil {
		return &item, errors.Wrapf(err, onRead+"can't unmarshal .Embedded (%s)", embedded)
	}

	return &item, nil
}

const onSave = "on dataSQLite.Save(): "

func (dataOp *dataSQLite) Save(items []data.Item, marksOp marks.Operator, indexerOp indexer.Operator, _ *crud.SaveOptions) ([]common.ID, error) {
	var ids []common.ID
	var errs common.Errors

	for _, item := range items {
		embedded, err := json.Marshal(item.Embedded)
		if err != nil {
			return ids, errs.Append(errors.Wrapf(err, onSave+"can't .Marshal: %s", item.Embedded)).Err()
		}

		index, err := json.Marshal(item.Index)
		if err != nil {
			return ids, errs.Append(errors.Wrapf(err, onSave+"can't .Marshal: %s", item.Index)).Err()
		}

		values := []interface{}{item.SourceID, item.SourceTime, item.SourceURL, item.Type, item.Title, item.Summary, item.Details, item.Href, embedded, strings.Join(item.Tags,
			"\n"), index, item.SourceKey, item.Origin}

		res, err := dataOp.stmSave.Exec(values...)
		if err != nil {
			return ids, errs.Append(errors.Wrapf(err, onSave+sqllib.CantExec, dataOp.sqlSave, values)).Err()
		}

		idSQLite, err := res.LastInsertId()
		if err != nil {
			return ids, errs.Append(errors.Wrapf(err, onSave+sqllib.CantGetLastInsertId, dataOp.sqlSave, values)).Err()
		}
		id := common.ID(strconv.FormatInt(idSQLite, 10))

		//var tagItems []tagItem
		//for _, tag := range item.Tags {
		//	tagItems = append(tagItems, tagItem{tag, &id})
		//}

		ids = append(ids, id)
		//errs = errs.Append(dataOp.saveTags(tagItems))
	}

	return ids, errs.Err()
}

const onRemove = "on dataSQLite.Remove()"

func (dataOp *dataSQLite) Remove(selectors.Term, marks.Operator, indexer.Operator, *crud.RemoveOptions) error {
	//		var err error
	//		var values []interface{}
	//		var orderAndLimit, condition, conditionCompleted string
	//
	//		if options != nil {
	//			condition, values, err = selectors.Mysql("", options.Selector)
	//			if err != nil {
	//				return crud.Result{}, errors.Wrapf(err, onDelete+"bad selector ('%#v')", options.Selector)
	//			}
	//
	//			conditionCompleted = condition
	//			if strings.TrimSpace(conditionCompleted) != "" {
	//				conditionCompleted = " where " + conditionCompleted
	//			}
	//
	//			orderAndLimit = mysqllib.OrderAndLimit(options.OrderBy, options.Limits)
	//		}
	//
	//		if strings.TrimSpace(condition) != "" {
	//			condition = "where " + condition
	//		}
	//
	//		sqlQuery := dsOp.sqlDelete + " " + condition + " " + orderAndLimit
	//		res, err := dsOp.db.Exec(sqlQuery, values...)
	//		if err != nil {
	//			return crud.Result{}, errors.Wrapf(err, onDelete+"can't exec SQL: %s, %s", sqlQuery, values)
	//		}
	//		cnt, err := res.RowsAffected()
	//		if err != nil {
	//			return crud.Result{}, errors.Wrapf(err, onDelete+"can't get RowsAffected(): %s, %s", sqlQuery, values)
	//		}
	//		return crud.Result{cnt}, nil
	//	}

	return common.ErrNotImplemented
}

const onList = "on dataSQLite.List()"

func (dataOp *dataSQLite) List(selector selectors.Term, indexerOp indexer.Operator, options *crud.GetOptions) ([]crud.Brief, error) {
	var values []interface{}

	rows, err := dataOp.stmList.Query(values...)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, onList+sqllib.CantQuery, dataOp.sqlList, values)
	}
	defer rows.Close()

	var briefs []crud.Brief

	for rows.Next() {
		brief := crud.Brief{Info: common.Info{}}

		var sourceTime *time.Time
		var sourceURL, tags string

		err = rows.Scan(&brief.ID, &sourceTime, &sourceURL, &brief.Type, &brief.Title, &brief.Summary, &tags)
		if err != nil {
			return briefs, errors.Wrapf(err, onList+sqllib.CantScanQueryRow, dataOp.sqlList, values)
		}

		if sourceTime != nil {
			brief.Info["source_time"] = *sourceTime
		}
		brief.Info["source_url"] = sourceURL
		brief.Info["tags"] = strings.Split(tags, "\n")

		briefs = append(briefs, brief)
	}
	err = rows.Err()
	if err != nil {
		return briefs, errors.Wrapf(err, onList+": "+sqllib.RowsError, dataOp.sqlList, values)
	}

	return briefs, nil
}

const onCount = "on dataSQLite.Count()"

func (dataOp *dataSQLite) Count(selectors.Term, indexer.Operator, *crud.GetOptions) ([]crud.Part, error) {
	return nil, common.ErrNotImplemented
}

const onReindex = "on dataSQLite.Reindex()"

func (dataOp *dataSQLite) Reindex(selectors.Term, indexer.Operator, *crud.GetOptions) error {
	return common.ErrNotImplemented
}

func (dataOp *dataSQLite) Close() error {
	return errors.Wrap(dataOp.db.Close(), "on dataSQLite.Close()")
}

func (dataOp *dataSQLite) Clean() error {
	_, err := dataOp.db.Exec("TRUNCATE " + tabledata)

	return err
}

//const onLastKey = "on datastoreMySQL.LastKey()"
//
//func (dsOp *datastoreMySQL) LastKey(class data.Type, options *crud.ReadOptions) (string, error) {
//
//	// TODO: use options!!!
//
//	values := []interface{}{string(class)}
//	rows, err := dsOp.stmLastKey.Query(values...)
//	if err == sql.ErrNoRows {
//		return "", nil
//	} else if err != nil {
//		return "", errors.Wrapf(err, onLastKey+"can't query (sql='%s', values='%#v')", dsOp.sqlLastKey, values)
//	}
//	defer rows.Close()
//	if rows.Next() {
//		var lastKey string
//		err = rows.Scan(&lastKey)
//		if err != nil {
//			return "", errors.Wrapf(err, onLastKey+"can't scan query row (sql='%s', values='%#v')", dsOp.sqlLastKey, values)
//		}
//		return lastKey, nil
//	}
//	err = rows.Err()
//	if err != nil {
//		return "", errors.Wrapf(err, onLastKey+"on rows.Err() (sql='%s', values='%#v')", dsOp.sqlLastKey, values)
//	}
//
//	return "", nil
//}
