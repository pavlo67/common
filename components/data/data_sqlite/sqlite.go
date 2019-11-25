package data_sqlite

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/libraries/sqllib"
	"github.com/pavlo67/workshop/common/libraries/sqllib/sqllib_sqlite"
	"github.com/pavlo67/workshop/components/crud"
	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/selector"
)

const tableDefault = "data"
const limitDefault = 200

var fieldsToInsert = []string{"title", "summary", "url", "embedded", "tags", "details", "source", "source_key", "source_time", "source_data"}
var fieldsToInsertStr = strings.Join(fieldsToInsert, ", ")

var fieldsToRead = append(fieldsToInsert, "created_at", "updated_at")
var fieldsToReadStr = strings.Join(fieldsToRead, ", ")

var fieldsToList = append([]string{"id"}, fieldsToRead...)
var fieldsToListStr = strings.Join(fieldsToList, ", ")

var _ data.Operator = &dataSQLite{}
var _ crud.Cleaner = &dataSQLite{}

type dataSQLite struct {
	limit int
	db    *sql.DB
	table string

	stmInsert, stmRead, stmRemove, stmList, stmCount *sql.Stmt
	sqlInsert, sqlRead, sqlRemove, sqlList, sqlCount string
}

const onNew = "on dataSQLite.New(): "

func NewData(access config.Access, table string, limit int) (data.Operator, crud.Cleaner, error) {
	db, err := sqllib_sqlite.Connect(access)
	if err != nil {
		return nil, nil, errors.Wrap(err, onNew)
	}

	if table == "" {
		table = tableDefault
	}
	if limit <= 0 {
		limit = limitDefault
	}

	dataOp := dataSQLite{
		db:    db,
		table: table,
		limit: limit,

		sqlRead: "SELECT " + fieldsToReadStr + " FROM " + table + " WHERE id = ?",

		sqlInsert: "INSERT INTO " + table + " (" + fieldsToInsertStr + ") VALUES (" + strings.Repeat(",? ", len(fieldsToInsert))[1:] + ")",
		sqlRemove: "DELETE FROM " + table + " where ID = ?",

		sqlList:  "SELECT " + fieldsToListStr + " FROM " + table + " ORDER BY created_at DESC LIMIT " + strconv.Itoa(limit),
		sqlCount: "SELECT count(*) FROM " + table + " WHERE source = ? AND source_key = ?",
	}

	sqlStmts := []sqllib.SqlStmt{
		{&dataOp.stmInsert, dataOp.sqlInsert},
		{&dataOp.stmRead, dataOp.sqlRead},
		{&dataOp.stmRemove, dataOp.sqlRemove},

		{&dataOp.stmList, dataOp.sqlList},
		{&dataOp.stmCount, dataOp.sqlCount},
	}

	for _, sqlStmt := range sqlStmts {
		if err := sqllib.Prepare(db, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, nil, errors.Wrap(err, onNew)
		}
	}

	return &dataOp, &dataOp, nil
}

const onSave = "on dataSQLite.Save(): "

func (dataOp *dataSQLite) Save(items []data.Item, _ *crud.SaveOptions) ([]common.ID, error) {
	var ids []common.ID

	for _, item := range items {
		embedded, err := json.Marshal(item.Embedded)
		if err != nil {
			return ids, errors.Wrapf(err, onSave+"can't .Marshal(%#v)", item.Embedded)
		}

		tags, err := json.Marshal(item.Tags)
		if err != nil {
			return ids, errors.Wrapf(err, onSave+"can't .Marshal(%#v)", item.Tags)
		}

		details, err := json.Marshal(item.Details)
		if err != nil {
			return ids, errors.Wrapf(err, onSave+"can't .Marshal(%#v)", item.Details)
		}

		values := []interface{}{
			item.Title, item.Summary, item.URL, embedded, tags, details,
			item.Source, item.Key, item.OriginTime, item.OriginData,
		}

		if item.ID != "" {
			// TODO!!! update
			// dataOp.Remove()
			// values = append(values, item.ID)

		} else {
			res, err := dataOp.stmInsert.Exec(values...)
			if err != nil {
				return ids, errors.Wrapf(err, onSave+sqllib.CantExec, dataOp.sqlInsert, values)
			}

			idSQLite, err := res.LastInsertId()
			if err != nil {
				return ids, errors.Wrapf(err, onSave+sqllib.CantGetLastInsertId, dataOp.sqlInsert, values)
			}
			id := common.ID(strconv.FormatInt(idSQLite, 10))

			// TODO!!! save tags

			//var tagItems []tagItem
			//for _, tag := range item.Tags {
			//	tagItems = append(tagItems, tagItem{tag, &id})
			//}

			ids = append(ids, id)
		}
	}

	return ids, nil
}

const onRead = "on dataSQLite.Read(): "

func (dataOp *dataSQLite) Read(id common.ID, _ *crud.GetOptions) (*data.Item, error) {
	if len(id) < 1 {
		return nil, errors.New(onRead + "empty ID")
	}

	idNum, err := strconv.ParseUint(string(id), 10, 64)
	if err != nil {
		return nil, errors.Errorf(onRead+"wrong ID (%s)", id)
	}

	item := data.Item{ID: id}
	var embedded, tags string

	err = dataOp.stmRead.QueryRow(idNum).Scan(
		&item.Title, &item.Summary, &item.URL, &embedded, &tags, &item.DetailsRaw,
		&item.Source, &item.Key, &item.OriginTime, &item.OriginData, &item.CreatedAt, &item.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrapf(err, onRead+sqllib.CantScanQueryRow, dataOp.sqlRead, idNum)
	}

	err = json.Unmarshal([]byte(tags), &item.Tags)
	if err != nil {
		return &item, errors.Wrapf(err, onRead+"can't unmarshal .Tags (%s)", tags)
	}

	err = json.Unmarshal([]byte(embedded), &item.Embedded)
	if err != nil {
		return &item, errors.Wrapf(err, onRead+"can't unmarshal .Embedded (%s)", embedded)
	}

	return &item, nil
}

const onDetails = "on dataSQLite.Details()"

func (dataOp *dataSQLite) Details(item *data.Item, exemplar interface{}) error {
	err := json.Unmarshal(item.DetailsRaw, exemplar)
	if err != nil {
		return errors.Wrapf(err, onDetails+"can't .Unmarshal(%#v)", item.DetailsRaw)
	}

	item.Details = exemplar

	return nil
}

const onRemove = "on dataSQLite.Remove()"

func (dataOp *dataSQLite) Remove(*selector.Term, *crud.RemoveOptions) error {
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

func (dataOp *dataSQLite) List(*selector.Term, *crud.GetOptions) ([]data.Item, error) {
	var values []interface{}

	rows, err := dataOp.stmList.Query(values...)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, onList+sqllib.CantQuery, dataOp.sqlList, values)
	}
	defer rows.Close()

	var items []data.Item

	for rows.Next() {
		var idNum int64
		var item data.Item
		var embedded, tags string

		err := rows.Scan(
			&idNum, &item.Title, &item.Summary, &item.URL, &embedded, &tags, &item.DetailsRaw,
			&item.Source, &item.Key, &item.OriginTime, &item.OriginData, &item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return items, errors.Wrapf(err, onList+sqllib.CantScanQueryRow, dataOp.sqlList, values)
		}

		err = json.Unmarshal([]byte(tags), &item.Tags)
		if err != nil {
			return items, errors.Wrapf(err, onList+"can't unmarshal .Tags (%s)", tags)
		}

		err = json.Unmarshal([]byte(embedded), &item.Embedded)
		if err != nil {
			return items, errors.Wrapf(err, onList+"can't unmarshal .Embedded (%s)", embedded)
		}

		if err != nil {
			return items, errors.Wrapf(err, onList+sqllib.CantScanQueryRow, dataOp.sqlList, values)
		}

		item.ID = common.ID(strconv.FormatInt(idNum, 10))
		items = append(items, item)
	}
	err = rows.Err()
	if err != nil {
		return items, errors.Wrapf(err, onList+": "+sqllib.RowsError, dataOp.sqlList, values)
	}

	return items, nil
}

const onCount = "on dataSQLite.Has(): "

func (dataOp *dataSQLite) Count(*selector.Term, *crud.GetOptions) ([]crud.Part, error) {
	//if len(originKey.Key) < 1 { // || len(originKey.ID) < 1
	//	return 0, errors.New(onCount + "empty ID")
	//}
	//
	//values := []interface{}{originKey.ID, originKey.Key}
	//
	//var cnt uint
	//err := dataOp.stmHas.QueryRow(values...).Scan(&cnt)
	//if err != nil {
	//	return cnt, errors.Wrapf(err, onCount+sqllib.CantScanQueryRow, dataOp.sqlHas, values)
	//}
	//
	return nil, common.ErrNotImplemented
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

func (dataOp *dataSQLite) Close() error {
	return errors.Wrap(dataOp.db.Close(), "on dataSQLite.Close()")
}

func (dataOp *dataSQLite) Clean() error {
	_, err := dataOp.db.Exec("DELETE FROM " + dataOp.table)

	return err
}
