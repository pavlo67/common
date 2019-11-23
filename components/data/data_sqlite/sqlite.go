package data_sqlite

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"

	"github.com/pavlo67/workshop/components/selector"
	"github.com/pavlo67/workshop/libraries/sqllib"

	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/instruments/indexer"
	"github.com/pavlo67/workshop/components/marks"
)

const limitDefault = 200

var tableData = "data"

var fieldsToList = []string{"id", "source_url", "types", "title", "summary", "tags", "saved_at"}
var fieldsToListStr = strings.Join(fieldsToList, ", ")

var fieldsToRead = []string{"source_id", "source_time", "source_url", "types", "title", "summary", "details", "embedded", "tags", "saved_at"}
var fieldsToReadStr = strings.Join(fieldsToRead, ", ")

var fieldsToSave = []string{"source_id", "source_time", "source_url", "types", "title", "summary", "details", "embedded", "tags", "indexes", "source_key", "origin"}
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

		sqlHas:  "SELECT count(*) FROM " + tableData + " WHERE source_id = ? AND source_key = ?",
		sqlRead: "SELECT " + fieldsToReadStr + " FROM " + tableData + " WHERE id = ?",

		sqlSave:   "INSERT INTO " + tableData + " (" + fieldsToSaveStr + ") VALUES (" + strings.Repeat(",? ", len(fieldsToSave))[1:] + ")",
		sqlRemove: "DELETE FROM " + tableData + " where ID = ?",

		sqlList: "SELECT " + fieldsToListStr + " FROM " + tableData + " ORDER BY saved_at DESC LIMIT " + strconv.Itoa(limit),
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

func (dataOp *dataSQLite) Has(originKey data.Origin, _ *crud.GetOptions) (uint, error) {
	if len(originKey.Key) < 1 { // || len(originKey.ID) < 1
		return 0, errors.New(onHas + "empty ID")
	}

	values := []interface{}{originKey.ID, originKey.Key}

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

	err = dataOp.stmRead.QueryRow(id).Scan(&item.ID, &item.OriginTime, &item.OriginURL, &item.Type, &item.Title, &item.Summary, &item.Details, &embedded, &tags, &item.SavedAt)
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
			return ids, errs.Append(errors.Wrapf(err, onSave+"can't .marshal: %s", item.Embedded)).Err()
		}

		index, err := json.Marshal(item.Index)
		if err != nil {
			return ids, errs.Append(errors.Wrapf(err, onSave+"can't .marshal: %s", item.Index)).Err()
		}

		values := []interface{}{item.ID, item.OriginTime, item.OriginURL, item.Type, item.Title, item.Summary, item.Details, embedded, strings.Join(item.Tags,
			"\n"), index, item.Key, item.OriginData}

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

func (dataOp *dataSQLite) Remove(*selector.Term, marks.Operator, indexer.Operator, *crud.RemoveOptions) error {
	//		var err error
	//		var values []interface{}
	//		var orderAndLimit, condition, conditionCompleted string
	//
	//		if options != nil {
	//			condition, values, err = selector.Mysql("", options.Selector)
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

func (dataOp *dataSQLite) List(selector *selector.Term, indexerOp indexer.Operator, options *crud.GetOptions) ([]data.Brief, error) {
	var values []interface{}

	rows, err := dataOp.stmList.Query(values...)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, onList+sqllib.CantQuery, dataOp.sqlList, values)
	}
	defer rows.Close()

	var briefs []data.Brief

	for rows.Next() {
		brief := data.Brief{}

		var id int64

		err = rows.Scan(&id, &brief.OriginURL, &brief.Type, &brief.Title, &brief.Summary, &brief.Tags, &brief.SavedAt)
		if err != nil {
			return briefs, errors.Wrapf(err, onList+sqllib.CantScanQueryRow, dataOp.sqlList, values)
		}

		brief.ID = common.ID(strconv.FormatInt(id, 10))
		briefs = append(briefs, brief)
	}
	err = rows.Err()
	if err != nil {
		return briefs, errors.Wrapf(err, onList+": "+sqllib.RowsError, dataOp.sqlList, values)
	}

	return briefs, nil
}

const onCount = "on dataSQLite.Count()"

func (dataOp *dataSQLite) Count(*selector.Term, indexer.Operator, *crud.GetOptions) ([]crud.Part, error) {
	return nil, common.ErrNotImplemented
}

const onReindex = "on dataSQLite.Reindex()"

func (dataOp *dataSQLite) Reindex(*selector.Term, indexer.Operator, *crud.GetOptions) error {
	return common.ErrNotImplemented
}

func (dataOp *dataSQLite) Close() error {
	return errors.Wrap(dataOp.db.Close(), "on dataSQLite.Close()")
}

func (dataOp *dataSQLite) Clean() error {
	_, err := dataOp.db.Exec("TRUNCATE " + tableData)

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
