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
	"github.com/pavlo67/workshop/common/identity"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/libraries/sqllib"
	"github.com/pavlo67/workshop/common/libraries/sqllib/sqllib_pg"
	"github.com/pavlo67/workshop/common/selectors"
	"github.com/pavlo67/workshop/common/selectors/selectors_sql"

	"github.com/pavlo67/workshop/components/tasks"
)

var fieldsToInsert = []string{"worker_type", "params", "status", "results", "history"}
var fieldsToInsertStr = strings.Join(fieldsToInsert, ",")

var fieldsToRead = append(fieldsToInsert, "created_at")
var fieldsToReadStr = strings.Join(fieldsToRead, ",")

var fieldsToList = append([]string{"id"}, fieldsToRead...)
var fieldsToListStr = strings.Join(fieldsToList, ",")

var fieldsToStartFinish = []string{"status", "results"}
var fieldsToStartFinishStr = sqllib_pg.WildcardsForUpdate(fieldsToStartFinish)
var fieldsToReadToStartFinishStr = strings.Join(fieldsToStartFinish, ",")

var _ tasks.Operator = &tasksPostgres{}
var _ crud.Cleaner = &tasksPostgres{}

type tasksPostgres struct {
	db    *sql.DB
	table string

	sqlInsert, sqlRead, sqlList, sqlReadToStart, sqlStartFinish, sqlReadToStartFinish, sqlClean string
	stmInsert, stmRead, stmList, stmReadToStart, stmStartFinish, stmReadToStartFinish           *sql.Stmt

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

		sqlReadToStartFinish: "SELECT " + fieldsToReadToStartFinishStr + " FROM " + table + " WHERE id = $1",
		sqlStartFinish:       "UPDATE " + table + " SET " + fieldsToStartFinishStr + " WHERE id = $" + strconv.Itoa(len(fieldsToStartFinish)+1),

		//sqlRemove: "DELETE FROM " + table + " where Key = $1",
		sqlClean: "DELETE FROM " + table,

		interfaceKey: interfaceKey,
	}

	sqlStmts := []sqllib.SqlStmt{
		{&tasksOp.stmInsert, tasksOp.sqlInsert},
		{&tasksOp.stmRead, tasksOp.sqlRead},
		{&tasksOp.stmList, tasksOp.sqlList},
		//	{&tasksOp.stmRemove, tasksOp.sqlRemove},

		{&tasksOp.stmReadToStartFinish, tasksOp.sqlReadToStartFinish},
		{&tasksOp.stmStartFinish, tasksOp.sqlStartFinish},
	}

	for _, sqlStmt := range sqlStmts {
		if err := sqllib.Prepare(db, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, nil, errors.Wrap(err, onNew)
		}
	}

	return &tasksOp, &tasksOp, nil
}

const onSave = "on tasksPostgres.Save(): "

func (tasksOp *tasksPostgres) Save(item tasks.Item, options *crud.SaveOptions) (common.ID, error) {

	var actor *identity.Key
	if options != nil {
		actor = options.Actor
	}

	item.History = append(item.History, crud.Action{
		Actor:  actor,
		Key:    crud.SavedAction,
		DoneAt: time.Now(),
	})

	history, err := json.Marshal(item.History)
	if err != nil {
		return "", errors.Wrapf(err, onSave+"can't marshal .History(%#v)", item)
	}

	var results []byte
	if len(item.Results) > 0 {
		results, err = json.Marshal(item.Results)
		if err != nil {
			return "", errors.Wrapf(err, onSave+"can't marshal .Results(%#v)", item)
		}

	}

	// TODO!!! be careful: old items can't be changed

	values := []interface{}{item.TypeKey, item.Content, item.Status, results, history}

	var lastInsertId uint64

	err = tasksOp.stmInsert.QueryRow(values...).Scan(&lastInsertId)
	if err != nil {
		return "", errors.Wrapf(err, onSave+sqllib.CantExec, tasksOp.sqlInsert, values)
	}

	return common.ID(strconv.FormatUint(lastInsertId, 10)), nil
}

const onRead = "on tasksPostgres.Read(): "

func (tasksOp *tasksPostgres) Read(id common.ID, _ *crud.GetOptions) (*tasks.Item, error) {
	if len(id) < 1 {
		return nil, errors.New(onRead + "empty Key")
	}

	idNum, err := strconv.ParseUint(string(id), 10, 64)
	if err != nil {
		return nil, errors.Errorf(onRead+"wrong Key (%s)", id)
	}

	item := tasks.Item{ID: id}
	var history, results []byte
	var createdAtStr string

	err = tasksOp.stmRead.QueryRow(idNum).Scan(
		&item.TypeKey, &item.Content, &item.Status, &results, &history, &createdAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrapf(err, onRead+sqllib.CantScanQueryRow, tasksOp.sqlRead, idNum)
	}

	if len(history) > 0 {
		err = json.Unmarshal(history, &item.History)
		if err != nil {
			return &item, errors.Wrapf(err, onRead+"can't unmarshal .History (%s)", history)
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
		var history, results []byte
		var createdAtStr string

		err := rows.Scan(
			&idNum, &item.TypeKey, &item.Content, &item.Status, &results, &history, &createdAtStr,
		)
		if err != nil {
			return items, errors.Wrapf(err, onList+sqllib.CantScanQueryRow, query, values)
		}

		if len(history) > 0 {
			err = json.Unmarshal(history, &item.History)
			if err != nil {
				return items, errors.Wrapf(err, onList+"can't unmarshal .History (%s)", history)
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

		items = append(items, item)
	}
	err = rows.Err()
	if err != nil {
		return items, errors.Wrapf(err, onList+": "+sqllib.RowsError, query, values)
	}

	return items, nil
}

const onReadToStartFinish = "on tasksPostgres.readToStartFinish(): "

func (tasksOp *tasksPostgres) readToStartFinish(idNum uint64) (tasks.Status, []tasks.Result, error) {
	var status tasks.Status
	var resultsBytes []byte
	var results []tasks.Result

	err := tasksOp.stmReadToStartFinish.QueryRow(idNum).Scan(&status, &resultsBytes)
	if err != nil {
		return "", nil, errors.Wrapf(err, onReadToStartFinish+sqllib.CantScanQueryRow, tasksOp.sqlReadToStartFinish, idNum)
	}

	if len(resultsBytes) > 0 {
		err = json.Unmarshal(resultsBytes, &results)
		if err != nil {
			return "", nil, errors.Wrapf(err, onReadToStartFinish+"can't unmarshal .Results (%s)", resultsBytes)
		}
	}

	return status, results, nil
}

const onStart = "on tasksPostgres.Start(): "

func (tasksOp *tasksPostgres) Start(id common.ID, _ *crud.SaveOptions) error {
	if len(id) < 1 {
		return errors.New(onStart + "empty Key")
	}

	idNum, err := strconv.ParseUint(string(id), 10, 64)
	if err != nil {
		return errors.Errorf(onStart+"wrong Key (%s)", id)
	}

	status, results, err := tasksOp.readToStartFinish(idNum)
	if err != nil {
		return errors.Wrap(err, onStart)
	}

	if status != "" {
		return errors.Wrapf(err, onStart+"can't start task (%s) due to non-empty status (%s), previous results = %#v)", id, status, results)
	}

	// setting results ---------------------------------------------------------------------------------

	now := time.Now()
	results = append(results, tasks.Result{
		Timing: tasks.Timing{StartedAt: &now},
	})

	resultsBytes, err := json.Marshal(results)
	if err != nil {
		return errors.Wrapf(err, onStart+"can't marshal .Results (%#v)", results)
	}

	// saving the updates -----------------------------------------------------------------------------

	values := []interface{}{tasks.StatusStarted, resultsBytes, idNum}
	_, err = tasksOp.stmStartFinish.Exec(values...)
	if err != nil {
		return errors.Wrapf(err, onStart+sqllib.CantExec, tasksOp.sqlStartFinish, values)
	}

	return nil
}

const onFinish = "on tasksPostgres.Finish(): "

func (tasksOp *tasksPostgres) Finish(id common.ID, result tasks.Result, _ *crud.SaveOptions) error {
	if len(id) < 1 {
		return errors.New(onFinish + "empty Key")
	}

	idNum, err := strconv.ParseUint(string(id), 10, 64)
	if err != nil {
		return errors.Errorf(onFinish+"wrong Key (%s)", id)
	}

	_, results, err := tasksOp.readToStartFinish(idNum)
	if err != nil {
		return errors.Wrap(err, onStart)
	}

	// TODO: check status

	// setting results ---------------------------------------------------------------------------------

	if len(results) > 0 && results[len(results)-1].FinishedAt == nil {
		result.StartedAt = results[len(results)-1].StartedAt
	}
	if result.FinishedAt == nil {
		now := time.Now()
		result.FinishedAt = &now
	}

	results = append(results, result)
	resultsBytes, err := json.Marshal(results)
	if err != nil {
		return errors.Wrapf(err, onFinish+"can't marshal .Results(%#v)", results)
	}

	// saving the updates -----------------------------------------------------------------------------

	values := []interface{}{"", resultsBytes, idNum}
	_, err = tasksOp.stmStartFinish.Exec(values...)
	if err != nil {
		return errors.Wrapf(err, onFinish+sqllib.CantExec, tasksOp.sqlStartFinish, values)
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
