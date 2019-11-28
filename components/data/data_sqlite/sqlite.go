package data_sqlite

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/libraries/sqllib"
	"github.com/pavlo67/workshop/common/libraries/sqllib/sqllib_sqlite"
	"github.com/pavlo67/workshop/common/selectors"

	"github.com/pavlo67/workshop/common/selectors/selectors_sql"
	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/tagger"
)

const tableDefault = "data"

var fieldsToUpdate = []string{"title", "summary", "embedded", "tags", "details"}
var fieldsToUpdateStr = strings.Join(fieldsToUpdate, " = ?, ") + " = ?"

var fieldsToInsert = append(fieldsToUpdate, "source", "source_key", "source_time", "source_data", "url")
var fieldsToInsertStr = strings.Join(fieldsToInsert, ", ")

var fieldsToRead = append(fieldsToInsert, "created_at", "updated_at")
var fieldsToReadStr = strings.Join(fieldsToRead, ", ")

var fieldsToList = append([]string{"id"}, fieldsToRead...)
var fieldsToListStr = strings.Join(fieldsToList, ", ")

var _ data.Operator = &dataSQLite{}
var _ crud.Cleaner = &dataSQLite{}

type dataSQLite struct {
	db    *sql.DB
	table string

	sqlInsert, sqlUpdate, sqlRead, sqlRemove, sqlList, sqlCount string
	stmInsert, stmUpdate, stmRead, stmRemove, stmList, stmCount *sql.Stmt

	taggerOp     tagger.Operator
	interfaceKey joiner.InterfaceKey
}

const onNew = "on dataSQLite.New(): "

func NewData(access config.Access, table string, taggerOp tagger.Operator, interfaceKey joiner.InterfaceKey) (data.Operator, crud.Cleaner, error) {
	db, err := sqllib_sqlite.Connect(access)
	if err != nil {
		return nil, nil, errors.Wrap(err, onNew)
	}

	if table == "" {
		table = tableDefault
	}

	dataOp := dataSQLite{
		db:    db,
		table: table,

		sqlInsert: "INSERT INTO " + table + " (" + fieldsToInsertStr + ") VALUES (" + strings.Repeat(",? ", len(fieldsToInsert))[1:] + ")",
		sqlUpdate: "UPDATE " + table + " SET " + fieldsToUpdateStr + " WHERE id = ?",
		sqlRemove: "DELETE FROM " + table + " where ID = ?",

		sqlRead: "SELECT " + fieldsToReadStr + " FROM " + table + " WHERE id = ?",
		// sqlCount: "SELECT count(*) FROM " + table + " WHERE source = ? AND source_key = ?",
		sqlList: sqlList(table, ""),

		taggerOp:     taggerOp,
		interfaceKey: interfaceKey,
	}

	sqlStmts := []sqllib.SqlStmt{
		{&dataOp.stmInsert, dataOp.sqlInsert},
		{&dataOp.stmUpdate, dataOp.sqlUpdate},
		{&dataOp.stmRemove, dataOp.sqlRemove},

		{&dataOp.stmRead, dataOp.sqlRead},
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

func sqlList(table, condition string) string {
	if strings.TrimSpace(condition) != "" {
		condition = " WHERE " + condition
	}
	return "SELECT " + fieldsToListStr + " FROM " + table + condition + " ORDER BY created_at DESC"
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

		if item.ID != "" {
			values := []interface{}{
				item.Title, item.Summary, embedded, tags, details, item.ID,
			}

			_, err := dataOp.stmUpdate.Exec(values...)
			if err != nil {
				return ids, errors.Wrapf(err, onSave+sqllib.CantExec, dataOp.sqlUpdate, values)
			}

			if dataOp.taggerOp != nil {
				err = dataOp.taggerOp.ReplaceTags(dataOp.interfaceKey, item.ID, item.Tags, nil)
				if err != nil {
					return ids, errors.Wrapf(err, onSave+": can't .ReplaceTags(%#v)", item.Tags)
				}
			}

			ids = append(ids, item.ID)

		} else {
			values := []interface{}{
				item.Title, item.Summary, embedded, tags, details,
				item.Origin.Source, item.Origin.Key, item.Origin.Time, item.Origin.Data, item.URL,
			}

			res, err := dataOp.stmInsert.Exec(values...)
			if err != nil {
				return ids, errors.Wrapf(err, onSave+sqllib.CantExec, dataOp.sqlInsert, values)
			}

			idSQLite, err := res.LastInsertId()
			if err != nil {
				return ids, errors.Wrapf(err, onSave+sqllib.CantGetLastInsertId, dataOp.sqlInsert, values)
			}
			id := common.ID(strconv.FormatInt(idSQLite, 10))

			if dataOp.taggerOp != nil && len(item.Tags) > 0 {
				err = dataOp.taggerOp.SaveTags(dataOp.interfaceKey, id, item.Tags, nil)
				if err != nil {
					return ids, errors.Wrapf(err, onSave+": can't .SaveTags(%#v)", item.Tags)
				}
			}

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
	var embedded, tags, createdAt string
	var sourceTimePtr, updatedAtPtr *string

	err = dataOp.stmRead.QueryRow(idNum).Scan(
		&item.Title, &item.Summary, &embedded, &tags, &item.DetailsRaw,
		&item.Source, &item.Key, &sourceTimePtr, &item.Data, &item.URL,
		&createdAt, &updatedAtPtr,
	)
	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrapf(err, onRead+sqllib.CantScanQueryRow, dataOp.sqlRead, idNum)
	}

	item.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return &item, errors.Wrapf(err, onRead+"can't parse .CreatedAt (%s)", createdAt)
	}

	if updatedAtPtr != nil {
		updatedAt, err := time.Parse(time.RFC3339, *updatedAtPtr)
		if err != nil {
			return &item, errors.Wrapf(err, onRead+"can't parse .UpdatedAt (%s)", *updatedAtPtr)
		}
		item.UpdatedAt = &updatedAt
	}

	if sourceTimePtr != nil {
		sourceTime, err := time.Parse(time.RFC3339, *sourceTimePtr)
		if err != nil {
			return &item, errors.Wrapf(err, onRead+"can't parse .SourceTime (%s)", *sourceTimePtr)
		}
		item.Origin.Time = &sourceTime
	}

	if len(tags) > 0 {
		err = json.Unmarshal([]byte(tags), &item.Tags)
		if err != nil {
			return &item, errors.Wrapf(err, onRead+"can't unmarshal .Tags (%s)", tags)
		}
	}

	if len(embedded) > 0 {
		err = json.Unmarshal([]byte(embedded), &item.Embedded)
		if err != nil {
			return &item, errors.Wrapf(err, onRead+"can't unmarshal .Embedded (%s)", embedded)
		}
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

func (dataOp *dataSQLite) Remove(id common.ID, _ *crud.RemoveOptions) error {
	if len(id) < 1 {
		return errors.New(onRemove + "empty ID")
	}

	idNum, err := strconv.ParseUint(string(id), 10, 64)
	if err != nil {
		return errors.Errorf(onRemove+"wrong ID (%s)", id)
	}

	_, err = dataOp.stmRemove.Exec(idNum)
	if err != nil {
		return errors.Wrapf(err, onRemove+sqllib.CantExec, dataOp.sqlRemove, idNum)
	}

	if dataOp.taggerOp != nil {
		err = dataOp.taggerOp.ReplaceTags(dataOp.interfaceKey, id, nil, nil)
		if err != nil {
			return errors.Wrapf(err, onRemove+": can't .ReplaceTags(%#v)", nil)
		}
	}

	return nil
}

const onList = "on dataSQLite.List()"

func (dataOp *dataSQLite) List(term *selectors.Term, _ *crud.GetOptions) ([]data.Item, error) {
	condition, values, err := selectors_sql.Use(term)

	query := dataOp.sqlList
	stm := dataOp.stmList

	if condition != "" {
		query = sqlList(dataOp.table, condition)
		stm, err = dataOp.db.Prepare(query)
		if err != nil {
			return nil, errors.Wrapf(err, onList+": can't db.Prepare(%s)", query)
		}
	}

	rows, err := stm.Query(values...)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, onList+sqllib.CantQuery, query, values)
	}
	defer rows.Close()

	var items []data.Item

	for rows.Next() {
		var idNum int64
		var item data.Item
		var embedded, tags, createdAt string
		var sourceTimePtr, updatedAtPtr *string

		err := rows.Scan(
			&idNum, &item.Title, &item.Summary, &embedded, &tags, &item.DetailsRaw,
			&item.Origin.Source, &item.Origin.Key, &sourceTimePtr, &item.Origin.Data, &item.URL,
			&createdAt, &updatedAtPtr,
		)
		if err != nil {
			return items, errors.Wrapf(err, onList+sqllib.CantScanQueryRow, query, values)
		}

		item.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return items, errors.Wrapf(err, onList+"can't parse .CreatedAt (%s)", createdAt)
		}

		if updatedAtPtr != nil {
			updatedAt, err := time.Parse(time.RFC3339, *updatedAtPtr)
			if err != nil {
				return items, errors.Wrapf(err, onList+"can't parse .UpdatedAt (%s)", *updatedAtPtr)
			}
			item.UpdatedAt = &updatedAt
		}

		if sourceTimePtr != nil {
			sourceTime, err := time.Parse(time.RFC3339, *sourceTimePtr)
			if err != nil {
				return items, errors.Wrapf(err, onList+"can't parse .SourceTime (%s)", *sourceTimePtr)
			}
			item.Origin.Time = &sourceTime
		}

		if len(tags) > 0 {
			err = json.Unmarshal([]byte(tags), &item.Tags)
			if err != nil {
				return items, errors.Wrapf(err, onList+"can't unmarshal .Tags (%s)", tags)
			}
		}

		if len(embedded) > 0 {
			err = json.Unmarshal([]byte(embedded), &item.Embedded)
			if err != nil {
				return items, errors.Wrapf(err, onList+"can't unmarshal .Embedded (%s)", embedded)
			}
		}

		if err != nil {
			return items, errors.Wrapf(err, onList+sqllib.CantScanQueryRow, query, values)
		}

		item.ID = common.ID(strconv.FormatInt(idNum, 10))
		items = append(items, item)
	}
	err = rows.Err()
	if err != nil {
		return items, errors.Wrapf(err, onList+": "+sqllib.RowsError, query, values)
	}

	return items, nil
}

const onCount = "on dataSQLite.Has(): "

func (dataOp *dataSQLite) Count(*selectors.Term, *crud.GetOptions) ([]crud.Counter, error) {
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

func (dataOp *dataSQLite) Clean(*selectors.Term) error {
	_, err := dataOp.db.Exec("DELETE FROM " + dataOp.table)

	if dataOp.taggerOp != nil {
		// TODO!!!
	}

	return err
}
