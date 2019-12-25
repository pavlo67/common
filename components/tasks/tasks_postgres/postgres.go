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
	"github.com/pavlo67/workshop/common/libraries/sqllib/sqllib_postgres"
	"github.com/pavlo67/workshop/common/selectors"
	"github.com/pavlo67/workshop/common/selectors/selectors_sql"

	"github.com/pavlo67/workshop/components/tasks"
)

var fieldsToInsert = []string{"worker_type", "params", "status", "results"}
var fieldsToInsertStr = strings.Join(fieldsToInsert, ",")

var fieldsToRead = append(fieldsToInsert, "created_at", "updated_at")
var fieldsToReadStr = strings.Join(fieldsToRead, ",")

var fieldsToList = append([]string{"id"}, fieldsToRead...)
var fieldsToListStr = strings.Join(fieldsToList, ",")

var fieldsToSetResults = []string{"status", "results", "updated_at"}
var fieldsToReadToSetStr = strings.Join(fieldsToSetResults, ",")
var fieldsToSetResultsStr = sqllib_postgres.WildcardsForUpdate(fieldsToSetResults)

var _ tasks.Operator = &tasksPostgres{}
var _ crud.Cleaner = &tasksPostgres{}

type tasksPostgres struct {
	db    *sql.DB
	table string

	sqlInsert, sqlRead, sqlList, sqlReadToSet, sqlSetResults, sqlClean string
	stmInsert, stmRead, stmList, stmReadToSet, stmSetResults           *sql.Stmt

	interfaceKey joiner.InterfaceKey
}

const onNew = "on tasksPostgres.New(): "

func New(access config.Access, table string, interfaceKey joiner.InterfaceKey) (tasks.Operator, crud.Cleaner, error) {
	db, err := sqllib_postgres.Connect(access)
	if err != nil {
		return nil, nil, errors.Wrap(err, onNew)
	}

	if table == "" {
		table = tasks.CollectionDefault
	}

	tasksOp := tasksPostgres{
		db:    db,
		table: table,

		sqlInsert: "INSERT INTO " + table + " (" + fieldsToInsertStr + ") VALUES (" + sqllib_postgres.WildcardsForInsert(fieldsToInsert) + ") RETURNING id",
		sqlRead:   "SELECT " + fieldsToReadStr + " FROM " + table + " WHERE id = $1",
		sqlList:   sqllib.SQLList(table, fieldsToListStr, "", &crud.GetOptions{OrderBy: []string{"created_at"}}),

		sqlReadToSet:  "SELECT " + fieldsToReadToSetStr + " FROM " + table + " WHERE id = $1",
		sqlSetResults: "UPDATE " + table + " SET " + fieldsToSetResultsStr + " WHERE id = $" + strconv.Itoa(len(fieldsToSetResults)+1),

		//sqlRemove: "DELETE FROM " + table + " where ID = $1",
		sqlClean: "DELETE FROM " + table,

		interfaceKey: interfaceKey,
	}

	sqlStmts := []sqllib.SqlStmt{
		{&tasksOp.stmInsert, tasksOp.sqlInsert},
		{&tasksOp.stmRead, tasksOp.sqlRead},
		{&tasksOp.stmList, tasksOp.sqlList},
		{&tasksOp.stmSetResults, tasksOp.sqlSetResults},
		{&tasksOp.stmReadToSet, tasksOp.sqlReadToSet},

		//	{&tasksOp.stmRemove, tasksOp.sqlRemove},
	}

	for _, sqlStmt := range sqlStmts {
		if err := sqllib.Prepare(db, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, nil, errors.Wrap(err, onNew)
		}
	}

	return &tasksOp, &tasksOp, nil
}

const onSave = "on tasksPostgres.Save(): "

func (tasksOp *tasksPostgres) Save(task tasks.Task, _ *crud.SaveOptions) (common.ID, error) {
	var paramsBytes []byte

	if task.Params != nil {
		var err error
		paramsBytes, err = json.Marshal(task.Params)
		if err != nil {
			return "", errors.Wrapf(err, onSave+"can't .Marshal(%#v)", task.Params)
		}
	}

	values := []interface{}{task.WorkerType, string(paramsBytes), "", ""}

	var lastInsertId uint64

	err := tasksOp.stmInsert.QueryRow(values...).Scan(&lastInsertId)
	if err != nil {
		return "", errors.Wrapf(err, onSave+sqllib.CantExec, tasksOp.sqlInsert, values)
	}

	return common.ID(strconv.FormatUint(lastInsertId, 10)), nil
}

const onRead = "on tasksPostgres.Read(): "

func (tasksOp *tasksPostgres) Read(id common.ID, _ *crud.GetOptions) (*tasks.Item, error) {
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

const onRemove = "on tasksPostgres.Remove()"

func (tasksOp *tasksPostgres) Remove(common.ID, *crud.RemoveOptions) error {
	return common.ErrNotImplemented
}

const onList = "on tasksPostgres.List()"

func (tasksOp *tasksPostgres) List(term *selectors.Term, options *crud.GetOptions) ([]tasks.Item, error) {
	condition, values, err := selectors_sql.Use(term)
	if err != nil {
		return nil, errors.Errorf(onList+"wrong selector (%#v): %s", term, err)
	}

	query := tasksOp.sqlList
	stm := tasksOp.stmList

	if condition != "" || options != nil {
		query = sqllib_postgres.CorrectWildcards(sqllib.SQLList(tasksOp.table, fieldsToListStr, condition, options))

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

const onSetResult = "on tasksPostgres.SetResult(): "

func (tasksOp *tasksPostgres) SetResult(id common.ID, result tasks.Result, _ *crud.SaveOptions) error {
	if len(id) < 1 {
		return errors.New(onSetResult + "empty ID")
	}

	idNum, err := strconv.ParseUint(string(id), 10, 64)
	if err != nil {
		return errors.Errorf(onSetResult+"wrong ID (%s)", id)
	}

	var statusStr, resultsStr string
	var updatedAtPtr *string

	err = tasksOp.stmReadToSet.QueryRow(idNum).Scan(&statusStr, &resultsStr, &updatedAtPtr)
	if err != nil {
		return errors.Wrapf(err, onSetResult+sqllib.CantScanQueryRow, tasksOp.sqlReadToSet, idNum)
	}

	//var status tasks.Status
	//if len(statusStr) > 0 {
	//	err = json.Unmarshal([]byte(statusStr), &status)
	//	if err != nil {
	//		return  errors.Wrapf(err, onSetResult+"can't unmarshal .Status (%s)", statusStr)
	//	}
	//}
	// TODO!!!
	statusBytes := []byte(statusStr)

	var results []tasks.Result
	if len(resultsStr) > 0 {
		err = json.Unmarshal([]byte(resultsStr), &results)
		if err != nil {
			return errors.Wrapf(err, onSetResult+"can't unmarshal .Results (%s)", resultsStr)
		}
	}
	results = append(results, result)
	resultsBytes, err := json.Marshal(results)
	if err != nil {
		return errors.Wrapf(err, onSetResult+"can't .Marshal(%#v)", results)
	}

	values := []interface{}{
		statusBytes, resultsBytes, time.Now().Format(time.RFC3339), id,
	}

	_, err = tasksOp.stmSetResults.Exec(values...)
	if err != nil {
		return errors.Wrapf(err, onSetResult+sqllib.CantExec, tasksOp.sqlSetResults, values)
	}

	return nil
}

func (tasksOp *tasksPostgres) Close() error {
	return errors.Wrap(tasksOp.db.Close(), "on tasksPostgres.Close()")
}

const onClean = "on tasksPostgres.Clean(): "

func (tasksOp *tasksPostgres) Clean(term *selectors.Term, _ *crud.RemoveOptions) error {
	condition, values, err := selectors_sql.Use(term)
	if err != nil {
		return errors.Errorf(onClean+"wrong selector (%#v): %s", term, err)
	}

	query := tasksOp.sqlClean
	if strings.TrimSpace(condition) != "" {
		query += " WHERE " + sqllib_postgres.CorrectWildcards(condition)
	}

	_, err = tasksOp.db.Exec(query, values...)
	if err != nil {
		return errors.Wrapf(err, onClean+sqllib.CantExec, query, values)
	}

	return nil
}
