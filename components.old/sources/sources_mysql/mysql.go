package sources_mysql

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"github.com/pavlo67/partes/crud"
	"github.com/pavlo67/partes/crud/selector"
	"github.com/pavlo67/partes/libs/mysqllib"
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/confidenter/groups"
	"github.com/pavlo67/punctum/confidenter/rights"
	"github.com/pavlo67/punctum/processor/sources"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/workshop/common"
)

type sourcesMySQL struct {
	dbh   *sql.DB
	table string

	grOp     groups.Operator
	managers rights.Managers

	stmCreate, stmRead, stmUpdate, stmDelete              *sql.Stmt
	sqlCreate, sqlRead, sqlUpdate, sqlDelete, sqlReadList string
}

var fields = []string{"url", "title", "import_type", "params", "r_view", "r_owner", "managers"}

const onNew = "on sources_mysql.New()"

func New(grOp groups.Operator, mysqlConfig config.ServerAccess, table string, managers rights.Managers) (*sourcesMySQL, error) {
	dbh, err := mysqllib.ConnectToMysql(mysqlConfig)
	if err != nil {
		return nil, errors.Wrap(err, onNew)
	}

	if grOp == nil {
		l.Warn(onNew + ": no groups.Operator")
	}

	if strings.TrimSpace(table) == "" {
		return nil, errors.New(onNew + ": no table name defined")
	}

	if managers == nil {
		managers = rights.Managers{rights.Create: basis.AnyoneRegistered}
	}

	fieldsToCreate := "`" + strings.Join(fields, "`, `") + "`"
	fieldsToRead := "id, `" + strings.Join(fields, "`, `") + "`, created_at, updated_at"
	fieldsToUpdate := "`" + strings.Join(fields, "` = ?, `") + "` = ?"

	srcOp := sourcesMySQL{
		grOp:        grOp,
		managers:    managers,
		dbh:         dbh,
		table:       table,
		sqlCreate:   "insert into `" + table + "` (" + fieldsToCreate + ") values (?,?,?, ?,?,?, ?)",
		sqlRead:     "select " + fieldsToRead + " from `" + table + "` where id = ?",
		sqlUpdate:   "update `" + table + "` set " + fieldsToUpdate + " where id = ?",
		sqlDelete:   "delete from `" + table + "` where id = ?",
		sqlReadList: "select SQL_CALC_FOUND_ROWS " + fieldsToRead + " from `" + table + "`",
	}

	sqlStmts := []mysqllib.SqlStmt{
		{&srcOp.stmCreate, srcOp.sqlCreate},
		{&srcOp.stmRead, srcOp.sqlRead},
		{&srcOp.stmUpdate, srcOp.sqlUpdate},
		{&srcOp.stmDelete, srcOp.sqlDelete},
	}

	for _, sqlStmt := range sqlStmts {
		if err = mysqllib.CreateStmt(dbh, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, errors.Wrap(err, onNew)
		}
	}

	return &srcOp, nil
}

const onCreate = "on sourcesMySQL.Create()"

func (srcOp *sourcesMySQL) Create(userIS common.ID, source sources.Item) (string, error) {
	if err := groups.OneOfErr(userIS, srcOp.grOp, srcOp.managers[rights.Create]); err != nil {
		return "", errors.Wrap(err, onCreate)
	}

	rView, rOwner, managers, err := groups.SetRights(userIS, srcOp.grOp, source.RView, source.ROwner, source.Managers)
	if err != nil {
		return "", errors.Wrap(err, onCreate+": can't .SetRights)")
	}

	var managersStr []byte
	if managers != nil {
		if managersStr, err = json.Marshal(managers); err != nil {
			return "", errors.Wrapf(err, onCreate+": can't json.marshal($#v)", managers)
		}
	}

	paramsStr, err := json.Marshal(source.Params)
	if err != nil {
		return "", errors.Wrapf(err, onCreate+": can't json.marshal(%#v)", source.Params)
	}

	values := []interface{}{source.URL, source.Title, string(source.Type), paramsStr, string(rView), string(rOwner), managersStr}
	res, err := srcOp.stmCreate.Exec(values...)
	if err != nil {
		return "", errors.Wrapf(err, onCreate+": can't exec SQL: %s, %#v", srcOp.sqlCreate, values)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return "", errors.Wrapf(err, onCreate+": can't get LastInsertId() SQL: %s, %#v", srcOp.sqlCreate, values)
	}

	return strconv.FormatInt(id, 10), nil
}

func (srcOp *sourcesMySQL) clean() error {
	_, err := srcOp.dbh.Exec("truncate `" + srcOp.table + "`")
	return err
}

const onRead = "on sourcesMySQL.Read()"

func (srcOp *sourcesMySQL) Read(userIS common.ID, idStr string) (*sources.Item, error) {

	if len(idStr) < 1 {
		return nil, errors.Wrap(crud.ErrEmptySelector, onRead)
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return nil, errors.Wrap(crud.ErrBadSelector, onRead+": "+idStr)
	}

	var src sources.Item
	var managers []byte
	err = srcOp.stmRead.QueryRow(id).Scan(&src.ID, &src.URL, &src.Title, &src.Type, &src.ParamsRaw, &src.RView, &src.ROwner, &managers, &src.CreatedAt, &src.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrapf(err, onRead+": can't exec QueryRow: %s, id = %s", srcOp.sqlRead, id)
	}

	if !groups.OneOf(userIS, srcOp.grOp, src.RView, src.ROwner) {
		return nil, errors.Wrap(basis.ErrNotFound, onRead)
	}

	if len(managers) > 0 {
		if err = json.Unmarshal(managers, &src.Managers); err != nil {
			return nil, errors.Wrapf(err, onRead+": can't json.unmarshal() for managers field(%s)", managers)
		}
	}

	return &src, nil
}

const onReadList = "on sourcesMySQL.ReadList()"

func (srcOp *sourcesMySQL) ReadList(userIS common.ID, options *content.ListOptions) ([]sources.Item, uint64, error) {
	var err error
	var values []interface{}
	var orderAndLimit, condition, conditionCompleted string

	if options != nil {
		condition, values, err = selector.Mysql(userIS, options.Selector)
		if err != nil {
			return nil, 0, errors.Wrapf(err, ": bad selector ('%#v')", options.Selector)
		}

		conditionCompleted = condition
		if strings.TrimSpace(conditionCompleted) != "" {
			conditionCompleted = " where " + conditionCompleted
		}

		orderAndLimit = mysqllib.OrderAndLimit(options.SortBy, options.Limits)
	}

	sqlQuery, rows, err := groups.QueryAccessible(srcOp.grOp, srcOp.dbh, userIS, srcOp.sqlReadList, condition, orderAndLimit, values)
	defer rows.Close()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, nil
		}
		return nil, 0, errors.Wrapf(err, onReadList+": can't get query (sql='%s', values='%#v')", sqlQuery, values)
	}

	items := []sources.Item{}
	for rows.Next() {
		var src sources.Item
		var managers []byte
		err = rows.Scan(&src.ID, &src.URL, &src.Title, &src.Type, &src.ParamsRaw, &src.RView, &src.ROwner, &managers, &src.CreatedAt, &src.UpdatedAt)
		if err != nil {
			return items, 0, errors.Wrapf(err, ": can't scan queryRow (sql='%s', values='%#v')", sqlQuery, values)
		}

		if len(managers) > 0 {
			if err = json.Unmarshal(managers, &src.Managers); err != nil {
				return items, 0, errors.Wrapf(err, onReadList+": can't json.unmarshal() for managers field(%s)", managers)
			}
		}

		items = append(items, src)
	}
	err = rows.Err()
	if err != nil {
		return items, 0, errors.Wrapf(err, onReadList+": on rows.Err(): sql='%s' (%#v)", sqlQuery, values)
	}

	var allCount uint64
	err = srcOp.dbh.QueryRow("SELECT FOUND_ROWS()").Scan(&allCount)
	if err != nil {
		return nil, 0, errors.Wrapf(err, onReadList+": can't scan ('SELECT FOUND_ROWS()') for sql=%s (%#v)", sqlQuery, values)
	}
	return items, allCount, nil
}

const onUpdate = "on sourcesMySQL.Update()"

func (srcOp *sourcesMySQL) Update(userIS common.ID, source sources.Item) (crud.Result, error) {
	source0, err := srcOp.Read(userIS, source.ID)
	if err != nil {
		return crud.Result{}, errors.Wrap(err, onUpdate+": can't .Read()")
	}

	if !groups.OneOf(userIS, srcOp.grOp, source0.ROwner, source0.Managers[rights.Change]) {
		return crud.Result{}, errors.Wrap(basis.ErrNotFound, onUpdate)
	}

	rView, rOwner, managers, err := groups.SetRights(userIS, srcOp.grOp, source.RView, source.ROwner, source.Managers)
	if err != nil {
		return crud.Result{}, errors.Wrap(err, onUpdate+": can't .SetRights)")
	}

	var managersStr []byte
	if managers != nil {
		if managersStr, err = json.Marshal(managers); err != nil {
			return crud.Result{}, errors.Wrapf(err, onUpdate+": can't json.marshal($#v)", managers)
		}
	}

	paramsStr, err := json.Marshal(source.Params)
	if err != nil {
		return crud.Result{}, errors.Wrapf(err, onUpdate+": can't json.marshal(%#v)", source.Params)
	}

	values := []interface{}{source.URL, source.Title, string(source.Type), paramsStr, string(rView), string(rOwner), managersStr, source.ID}
	res, err := srcOp.stmUpdate.Exec(values...)
	if err != nil {
		return crud.Result{}, errors.Wrapf(err, onUpdate+": can't exec SQL: %s, %#v", srcOp.sqlUpdate, values)
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return crud.Result{}, errors.Wrapf(err, onUpdate+": can't get RowsAffected(): %s (%#v)", srcOp.sqlUpdate, values)
	}

	return crud.Result{NumOk: cnt}, nil
}

const onDelete = "on sourcesMySQL.DeleteList()"

func (srcOp *sourcesMySQL) Delete(userIS common.ID, id string) (crud.Result, error) {
	source0, err := srcOp.Read(userIS, id)
	if err != nil {
		return crud.Result{}, errors.Wrap(err, onDelete+": can't .Read()")
	}

	if !groups.OneOf(userIS, srcOp.grOp, source0.ROwner, source0.Managers[rights.Change]) {
		return crud.Result{}, errors.Wrap(basis.ErrNotFound, onDelete)
	}

	res, err := srcOp.stmDelete.Exec(id)
	if err != nil {
		return crud.Result{}, errors.Wrapf(err, onDelete+": can't exec SQL: %s, %s", srcOp.sqlDelete, id)
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return crud.Result{}, errors.Wrapf(err, onDelete+": can't get RowsAffected(): %s, %s", srcOp.sqlDelete, id)
	}
	return crud.Result{cnt}, nil
}

func (srcOp *sourcesMySQL) Close() error {
	return errors.Wrap(srcOp.dbh.Close(), "on sourcesMySQL.dbh.Close()")
}
