package tasks_postgres

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

	"github.com/pavlo67/workshop/components/tasks"
)

var fieldsToInsert = []string{"type", "params"}
var fieldsToInsertStr = strings.Join(fieldsToInsert, ", ")

var fieldsToRead = append(fieldsToInsert, "status", "results", "created_at", "updated_at")
var fieldsToReadStr = strings.Join(fieldsToRead, ", ")

var fieldsToList = append([]string{"id"}, fieldsToRead...)
var fieldsToListStr = strings.Join(fieldsToList, ", ")

var _ tasks.Operator = &tasksSQLite{}
var _ crud.Cleaner = &tasksSQLite{}

type tasksSQLite struct {
	db    *sql.DB
	table string

	sqlInsert, sqlRead, sqlList string
	stmInsert, stmRead, stmList *sql.Stmt

	interfaceKey joiner.InterfaceKey
}

const onNew = "on tasksSQLite.New(): "

func New(access config.Access, table string, interfaceKey joiner.InterfaceKey) (tasks.Operator, crud.Cleaner, error) {
	db, err := sqllib_sqlite.Connect(access)
	if err != nil {
		return nil, nil, errors.Wrap(err, onNew)
	}

	if table == "" {
		table = tasks.CollectionDefault
	}

	tasksOp := tasksSQLite{
		db:    db,
		table: table,

		sqlInsert: "INSERT INTO " + table + " (" + fieldsToInsertStr + ") VALUES (" + strings.Repeat(",? ", len(fieldsToInsert))[1:] + ")",
		sqlRead:   "SELECT " + fieldsToReadStr + " FROM " + table + " WHERE id = ?",
		sqlList:   sqllib.SQLList(table, fieldsToListStr, "", &crud.GetOptions{OrderBy: []string{"created_at DESC"}}),

		//sqlUpdate: "UPDATE " + table + " SET " + fieldsToUpdateStr + " WHERE id = ?",
		//sqlRemove: "DELETE FROM " + table + " where ID = ?",
		//sqlClean: "DELETE FROM " + table,

		interfaceKey: interfaceKey,
	}

	sqlStmts := []sqllib.SqlStmt{
		{&tasksOp.stmInsert, tasksOp.sqlInsert},
		{&tasksOp.stmRead, tasksOp.sqlRead},
		{&tasksOp.stmList, tasksOp.sqlList},

		//	{&tasksOp.stmUpdate, tasksOp.sqlUpdate},
		//	{&tasksOp.stmRemove, tasksOp.sqlRemove},
		//	{&tasksOp.stmClean, tasksOp.sqlClean},
	}

	for _, sqlStmt := range sqlStmts {
		if err := sqllib.Prepare(db, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, nil, errors.Wrap(err, onNew)
		}
	}

	return &tasksOp, &tasksOp, nil
}

const onSave = "on tasksSQLite.Save(): "

func (tasksOp *tasksSQLite) Save(task tasks.Task, _ *crud.SaveOptions) (common.ID, error) {
	var paramsBytes []byte

	if task.Params != nil {
		var err error
		paramsBytes, err = json.Marshal(task.Params)
		if err != nil {
			return "", errors.Wrapf(err, onSave+"can't .Marshal(%#v)", task.Params)
		}
	}

	values := []interface{}{task.WorkerType, string(paramsBytes)}

	res, err := tasksOp.stmInsert.Exec(values...)
	if err != nil {
		return "", errors.Wrapf(err, onSave+sqllib.CantExec, tasksOp.sqlInsert, values)
	}

	idSQLite, err := res.LastInsertId()
	if err != nil {
		return "", errors.Wrapf(err, onSave+sqllib.CantGetLastInsertId, tasksOp.sqlInsert, values)
	}

	return common.ID(strconv.FormatInt(idSQLite, 10)), nil
}

const onRead = "on tasksSQLite.Read(): "

func (tasksOp *tasksSQLite) Read(id common.ID, _ *crud.GetOptions) (*tasks.Item, error) {
	if len(id) < 1 {
		return nil, errors.New(onRead + "empty ID")
	}

	idNum, err := strconv.ParseUint(string(id), 10, 64)
	if err != nil {
		return nil, errors.Errorf(onRead+"wrong ID (%s)", id)
	}

	item := tasks.Item{ID: id}
	var status, results, params, createdAt string
	var updatedAtPtr *string

	err = tasksOp.stmRead.QueryRow(idNum).Scan(
		&item.WorkerType, &params, &status, &results, &createdAt, &updatedAtPtr,
	)
	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrapf(err, onRead+sqllib.CantScanQueryRow, tasksOp.sqlRead, idNum)
	}

	if len(params) > 0 {
		err = json.Unmarshal([]byte(params), &item.Params)
		if err != nil {
			return &item, errors.Wrapf(err, onRead+"can't unmarshal .Params (%s)", params)
		}
	}

	if len(status) > 0 {
		err = json.Unmarshal([]byte(status), &item.Status)
		if err != nil {
			return &item, errors.Wrapf(err, onRead+"can't unmarshal .Status (%s)", status)
		}
	}

	if len(results) > 0 {
		err = json.Unmarshal([]byte(results), &item.Results)
		if err != nil {
			return &item, errors.Wrapf(err, onRead+"can't unmarshal .Results (%s)", results)
		}
	}

	item.History.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return &item, errors.Wrapf(err, onRead+"can't parse .CreatedAt (%s)", createdAt)
	}

	if updatedAtPtr != nil {
		updatedAt, err := time.Parse(time.RFC3339, *updatedAtPtr)
		if err != nil {
			return &item, errors.Wrapf(err, onRead+"can't parse .UpdatedAt (%s)", *updatedAtPtr)
		}
		item.History.UpdatedAt = &updatedAt
	}

	return &item, nil
}

const onRemove = "on tasksSQLite.Remove()"

func (tasksOp *tasksSQLite) Remove(common.ID, *crud.RemoveOptions) error {
	return common.ErrNotImplemented
}

const onList = "on tasksSQLite.List()"

func (tasksOp *tasksSQLite) List(term *selectors.Term, options *crud.GetOptions) ([]tasks.Item, error) {
	condition, values, err := selectors_sql.Use(term)
	if err != nil {
		return nil, errors.Errorf(onList+"wrong selector (%#v): %s", term, err)
	}

	query := tasksOp.sqlList
	stm := tasksOp.stmList

	if condition != "" || options != nil {
		query = sqllib.SQLList(tasksOp.table, fieldsToListStr, condition, options)
		stm, err = tasksOp.db.Prepare(query)
		if err != nil {
			return nil, errors.Wrapf(err, onList+": can't db.Prepare(%s)", query)
		}
	}

	// l.Infof("%s / %#v\n%s", condition, values, query)

	rows, err := stm.Query(values...)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, onList+sqllib.CantQuery, query, values)
	}
	defer rows.Close()

	var items []tasks.Item

	for rows.Next() {
		var idNum int64
		var item tasks.Item
		var status, results, params, createdAt string
		var updatedAtPtr *string

		err := rows.Scan(
			&idNum, &item.WorkerType, &params, &status, &results, &createdAt, &updatedAtPtr,
		)
		if err != nil {
			return items, errors.Wrapf(err, onList+sqllib.CantScanQueryRow, query, values)
		}

		if len(params) > 0 {
			err = json.Unmarshal([]byte(params), &item.Params)
			if err != nil {
				return items, errors.Wrapf(err, onRead+"can't unmarshal .Params (%s)", params)
			}
		}

		if len(status) > 0 {
			err = json.Unmarshal([]byte(status), &item.Status)
			if err != nil {
				return items, errors.Wrapf(err, onRead+"can't unmarshal .Status (%s)", status)
			}
		}

		if len(results) > 0 {
			err = json.Unmarshal([]byte(results), &item.Results)
			if err != nil {
				return items, errors.Wrapf(err, onRead+"can't unmarshal .Results (%s)", results)
			}
		}

		item.History.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return items, errors.Wrapf(err, onRead+"can't parse .CreatedAt (%s)", createdAt)
		}

		if updatedAtPtr != nil {
			updatedAt, err := time.Parse(time.RFC3339, *updatedAtPtr)
			if err != nil {
				return items, errors.Wrapf(err, onRead+"can't parse .UpdatedAt (%s)", *updatedAtPtr)
			}
			item.History.UpdatedAt = &updatedAt
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

const onSetResult = "on tasksSQLite.SetResult(): "

func (tasksOp *tasksSQLite) SetResult(common.ID, tasks.Result, *crud.SaveOptions) error {
	return nil
}

func (tasksOp *tasksSQLite) Close() error {
	return errors.Wrap(tasksOp.db.Close(), "on tasksSQLite.Close()")
}

const onClean = "on tasksSQLite.Clean(): "

func (tasksOp *tasksSQLite) Clean(term *selectors.Term, _ *crud.RemoveOptions) error {
	//var termTags *selectors.Term
	//
	//condition, values, err := selectors_sql.Use(term)
	//
	//if strings.TrimSpace(condition) != "" {
	//	ids, err := tasksOp.ids(condition, values)
	//
	//	query := tasksOp.sqlClean + " WHERE " + condition
	//	_, err = tasksOp.db.Exec(query, values...)
	//	if err != nil {
	//		return errors.Wrapf(err, onClean+sqllib.CantExec, query, values)
	//	}
	//
	//	termTags = logic.AND(selectors.In("key", tasksOp.interfaceKey), selectors.In("id", ids...))
	//
	//} else {
	//	_, err = tasksOp.stmClean.Exec()
	//	if err != nil {
	//		return errors.Wrapf(err, onClean+sqllib.CantExec, tasksOp.sqlClean, nil)
	//	}
	//
	//	termTags = selectors.In("key", tasksOp.interfaceKey) // TODO!!! correct field key
	//}
	//
	//if tasksOp.taggerCleaner != nil {
	//	err = tasksOp.taggerCleaner.Clean(termTags, nil)
	//	if err != nil {
	//		return errors.Wrap(err, onClean)
	//	}
	//}
	//
	//return err

	return nil
}
