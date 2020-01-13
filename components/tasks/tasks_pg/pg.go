package tasks_pg

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
	"github.com/pavlo67/workshop/common/libraries/sqllib/sqllib_pg"
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

var fieldsToStart = []string{"status", "updated_at"}
var fieldsToStartStr = sqllib_pg.WildcardsForUpdate(fieldsToStart)

// var fieldsToReadToStartStr = strings.Join(fieldsToStart[:len(fieldsToStart)-1], ",")

var fieldsToFinish = []string{"status", "results", "updated_at"}
var fieldsToFinishStr = sqllib_pg.WildcardsForUpdate(fieldsToFinish)
var fieldsToReadToFinishStr = strings.Join(fieldsToFinish[:len(fieldsToFinish)-1], ",")

var _ tasks.Operator = &tasksPostgres{}
var _ crud.Cleaner = &tasksPostgres{}

type tasksPostgres struct {
	db    *sql.DB
	table string

	sqlInsert, sqlRead, sqlList, sqlReadToStart, sqlStart, sqlReadToFinish, sqlFinish, sqlClean string
	stmInsert, stmRead, stmList, stmReadToStart, stmStart, stmReadToFinish, stmFinish           *sql.Stmt

	interfaceKey joiner.InterfaceKey
}

const onNew = "on tasksPostgres.New(): "

func New(access config.Access, table string, interfaceKey joiner.InterfaceKey) (tasks.Operator, crud.Cleaner, error) {
	db, err := sqllib_pg.Connect(access)
	if err != nil {
		return nil, nil, errors.Wrap(err, onNew)
	}

	if table == "" {
		table = tasks.CollectionDefault
	}

	tasksOp := tasksPostgres{
		db:    db,
		table: table,

		sqlInsert: "INSERT INTO " + table + " (" + fieldsToInsertStr + ") VALUES (" + sqllib_pg.WildcardsForInsert(fieldsToInsert) + ") RETURNING id",
		sqlRead:   "SELECT " + fieldsToReadStr + " FROM " + table + " WHERE id = $1",
		sqlList:   sqllib.SQLList(table, fieldsToListStr, "", &crud.GetOptions{OrderBy: []string{"created_at"}}),

		// sqlReadToStart: "SELECT " + fieldsToReadToStartStr + " FROM " + table + " WHERE id = $1",
		sqlStart: "UPDATE " + table + " SET " + fieldsToStartStr + " WHERE id = $" + strconv.Itoa(len(fieldsToStart)+1),

		sqlReadToFinish: "SELECT " + fieldsToReadToFinishStr + " FROM " + table + " WHERE id = $1",
		sqlFinish:       "UPDATE " + table + " SET " + fieldsToFinishStr + " WHERE id = $" + strconv.Itoa(len(fieldsToFinish)+1),

		//sqlRemove: "DELETE FROM " + table + " where ID = $1",
		sqlClean: "DELETE FROM " + table,

		interfaceKey: interfaceKey,
	}

	sqlStmts := []sqllib.SqlStmt{
		{&tasksOp.stmInsert, tasksOp.sqlInsert},
		{&tasksOp.stmRead, tasksOp.sqlRead},
		{&tasksOp.stmList, tasksOp.sqlList},
		//	{&tasksOp.stmRemove, tasksOp.sqlRemove},

		// {&tasksOp.stmReadToStart, tasksOp.sqlReadToStart},
		{&tasksOp.stmStart, tasksOp.sqlStart},

		{&tasksOp.stmReadToFinish, tasksOp.sqlReadToFinish},
		{&tasksOp.stmFinish, tasksOp.sqlFinish},
	}

	for _, sqlStmt := range sqlStmts {
		if err := sqllib.Prepare(db, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, nil, errors.Wrap(err, onNew)
		}
	}

	return &tasksOp, &tasksOp, nil
}

const onSave = "on tasksPostgres.Save(): "

func (tasksOp *tasksPostgres) Save(task crud.Data, _ *crud.SaveOptions) (common.ID, error) {
	var content interface{}
	if task.Content != nil {
		content = task.Content
	} else {
		content = ""
	}

	values := []interface{}{task.TypeKey, content, "", ""}

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
	var status, results []byte
	var createdAtStr string
	var updatedAtPtr *string

	err = tasksOp.stmRead.QueryRow(idNum).Scan(
		&item.TypeKey, &item.Content, &status, &results, &createdAtStr, &updatedAtPtr,
	)
	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrapf(err, onRead+sqllib.CantScanQueryRow, tasksOp.sqlRead, idNum)
	}

	if len(status) > 0 {
		err = json.Unmarshal(status, &item.Status)
		if err != nil {
			return &item, errors.Wrapf(err, onRead+"can't unmarshal .History (%s)", status)
		}
	}

	if len(results) > 0 {
		err = json.Unmarshal(results, &item.Results)
		if err != nil {
			return &item, errors.Wrapf(err, onRead+"can't unmarshal .Results (%s)", results)
		}
	}

	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		// TODO??? return &item, errors.Wrapf(err, onRead+"can't parse .CreatedAt (%s)", createdAtStr)
	} else {
		item.History = item.History.SaveAction(crud.Action{Key: crud.CreatedAction, DoneAt: createdAt, Related: &joiner.Link{InterfaceKey: tasks.InterfaceKey, ID: id}})
	}

	if updatedAtPtr != nil {
		updatedAt, err := time.Parse(time.RFC3339, *updatedAtPtr)
		if err != nil {
			// TODO??? return &item, errors.Wrapf(err, onRead+"can't parse .UpdatedAt (%s)", *updatedAtPtr)
		}
		item.History = item.History.SaveAction(crud.Action{Key: crud.UpdatedAction, DoneAt: updatedAt, Related: &joiner.Link{InterfaceKey: tasks.InterfaceKey, ID: id}})
	}

	return &item, nil
}

const onRemove = "on tasksPostgres.Remove()"

func (tasksOp *tasksPostgres) Remove(common.ID, *crud.RemoveOptions) error {
	return common.ErrNotImplemented
}

const onList = "on tasksPostgres.ListTags()"

func (tasksOp *tasksPostgres) List(term *selectors.Term, options *crud.GetOptions) ([]tasks.Item, error) {
	condition, values, err := selectors_sql.Use(term)
	if err != nil {
		return nil, errors.Errorf(onList+"wrong selector (%#v): %s", term, err)
	}

	query := tasksOp.sqlList
	stm := tasksOp.stmList

	if condition != "" || options != nil {
		query = sqllib_pg.CorrectWildcards(sqllib.SQLList(tasksOp.table, fieldsToListStr, condition, options))

		stm, err = tasksOp.db.Prepare(query)
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

	var items []tasks.Item

	for rows.Next() {
		var idNum int64
		var item tasks.Item
		var status, results []byte
		var createdAtStr string
		var updatedAtPtr *string

		err := rows.Scan(
			&idNum, &item.TypeKey, &item.Content, &status, &results, &createdAtStr, &updatedAtPtr,
		)
		if err != nil {
			return items, errors.Wrapf(err, onList+sqllib.CantScanQueryRow, query, values)
		}

		if len(status) > 0 {
			err = json.Unmarshal(status, &item.Status)
			if err != nil {
				return items, errors.Wrapf(err, onList+"can't unmarshal .History (%s)", status)
			}
		}

		if len(results) > 0 {
			err = json.Unmarshal(results, &item.Results)
			if err != nil {
				return items, errors.Wrapf(err, onList+"can't unmarshal .Results (%s)", results)
			}
		}

		item.ID = common.ID(strconv.FormatInt(idNum, 10))

		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			// TODO???  return &item, errors.Wrapf(err, onRead+"can't parse .CreatedAt (%s)", createdAtStr)
		} else {
			item.History = item.History.SaveAction(crud.Action{Key: crud.CreatedAction, DoneAt: createdAt, Related: &joiner.Link{InterfaceKey: tasks.InterfaceKey, ID: item.ID}})
		}

		if updatedAtPtr != nil {
			updatedAt, err := time.Parse(time.RFC3339, *updatedAtPtr)
			if err != nil {
				// TODO??? return &item, errors.Wrapf(err, onRead+"can't parse .UpdatedAt (%s)", *updatedAtPtr)
			}
			item.History = item.History.SaveAction(crud.Action{Key: crud.UpdatedAction, DoneAt: updatedAt, Related: &joiner.Link{InterfaceKey: tasks.InterfaceKey, ID: item.ID}})
		}

		items = append(items, item)
	}
	err = rows.Err()
	if err != nil {
		return items, errors.Wrapf(err, onList+": "+sqllib.RowsError, query, values)
	}

	return items, nil
}

const onStart = "on tasksPostgres.Start(): "

func (tasksOp *tasksPostgres) Start(id common.ID, _ *crud.SaveOptions) error {
	if len(id) < 1 {
		return errors.New(onStart + "empty ID")
	}

	idNum, err := strconv.ParseUint(string(id), 10, 64)
	if err != nil {
		return errors.Errorf(onStart+"wrong ID (%s)", id)
	}

	//err = tasksOp.stmReadToStart.QueryRow(idNum).Scan(&statusStr)
	//if err != nil {
	//	return errors.Wrapf(err, onStart+sqllib.CantScanQueryRow, tasksOp.sqlReadToStart, idNum)
	//}

	// setting status ---------------------------------------------------------------------------------

	startedAt := time.Now()
	status := tasks.Status{tasks.Timing{StartedAt: &startedAt}}
	statusBytes, err := json.Marshal(status)
	if err != nil {
		return errors.Wrapf(err, onStart+"can't marshal .History (%#v)", status)
	}

	// saving the updates -----------------------------------------------------------------------------

	values := []interface{}{statusBytes, startedAt.Format(time.RFC3339), idNum}
	_, err = tasksOp.stmStart.Exec(values...)
	if err != nil {
		return errors.Wrapf(err, onStart+sqllib.CantExec, tasksOp.sqlStart, values)
	}

	return nil
}

const onFinish = "on tasksPostgres.Finish(): "

func (tasksOp *tasksPostgres) Finish(id common.ID, result tasks.Result, _ *crud.SaveOptions) error {
	if len(id) < 1 {
		return errors.New(onFinish + "empty ID")
	}

	idNum, err := strconv.ParseUint(string(id), 10, 64)
	if err != nil {
		return errors.Errorf(onFinish+"wrong ID (%s)", id)
	}

	var statusStr, resultsStr string

	err = tasksOp.stmReadToFinish.QueryRow(idNum).Scan(&statusStr, &resultsStr)
	if err != nil {
		return errors.Wrapf(err, onFinish+sqllib.CantScanQueryRow, tasksOp.sqlReadToFinish, idNum)
	}

	// clearing status --------------------------------------------------------------------------------

	var status tasks.Status
	if len(statusStr) > 0 {
		err = json.Unmarshal([]byte(statusStr), &status)
		if err != nil {
			return errors.Wrapf(err, onFinish+"can't unmarshal .History (%s)", statusStr)
		}
	}
	if result.StartedAt == nil && status.StartedAt != nil {
		result.StartedAt = status.StartedAt
	}
	statusStr = ""

	// adding the result ------------------------------------------------------------------------------

	var results []tasks.Result
	if len(resultsStr) > 0 {
		err = json.Unmarshal([]byte(resultsStr), &results)
		if err != nil {
			return errors.Wrapf(err, onFinish+"can't unmarshal .Results (%s)", resultsStr)
		}
	}

	if result.FinishedAt == nil {
		now := time.Now()
		result.FinishedAt = &now
	}

	results = append(results, result)
	resultsBytes, err := json.Marshal(results)
	if err != nil {
		return errors.Wrapf(err, onFinish+"can't .Marshal(%#v)", results)
	}

	// saving the updates -----------------------------------------------------------------------------

	values := []interface{}{statusStr, resultsBytes, time.Now().Format(time.RFC3339), idNum}
	_, err = tasksOp.stmFinish.Exec(values...)
	if err != nil {
		return errors.Wrapf(err, onFinish+sqllib.CantExec, tasksOp.sqlFinish, values)
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
		query += " WHERE " + sqllib_pg.CorrectWildcards(condition)
	}

	_, err = tasksOp.db.Exec(query, values...)
	if err != nil {
		return errors.Wrapf(err, onClean+sqllib.CantExec, query, values)
	}

	return nil
}

func (tasksOp *tasksPostgres) SelectToClean(*crud.RemoveOptions) (*selectors.Term, error) {
	return nil, common.ErrNotImplemented
}
