package mysqllib

import (
	"database/sql"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"github.com/pavlo67/partes/crud"
	"github.com/pavlo67/partes/crud/selectors"
	"github.com/pavlo67/constructor/auth"
	"github.com/pavlo67/constructor/basis"
	"github.com/pavlo67/constructor/starter/config"

	"github.com/pavlo67/constructor/confidenter/groups"
	"github.com/pavlo67/constructor/confidenter/rights"
)

var _ crud.Operator = &crudMySQL{}

type crudMySQL struct {
	grOp     groups.Operator
	managers rights.Managers

	dbh    *sql.DB
	table  string
	fields []crud.Field

	stmCreate, stmRead, stmUpdate, stmDelete              *sql.Stmt
	sqlCreate, sqlRead, sqlUpdate, sqlDelete, sqlReadList string
}

const onNew = "on crudmysql.NewCRUDOperator()"

func NewCRUDOperator(grOp groups.Operator, mysqlConfig config.ServerAccess, table string, fields []crud.Field, managers rights.Managers) (crud.Operator, crud.Cleaner, error) {
	dbh, err := ConnectToMysql(mysqlConfig)
	if err != nil {
		return nil, nil, errors.Wrap(err, onNew)
	}

	if grOp == nil {
		// DO WARNING, it's may be not ok
	}

	if strings.TrimSpace(table) == "" {
		return nil, nil, errors.New(onNew + ": no table name defined")
	}

	if len(fields) < 1 {
		return nil, nil, errors.New(onNew + ": no table fields defined")
	}

	// TODO: customize special fields - id, created_at, updated_at, r_view, r_owner, managers

	// {Key: "id", Unique: true, AutoUnique: true},
	// {Key: "r_view", Creatable: true, Editable: true, NotEmpty: true},
	// {Key: "r_owner", Creatable: true, Editable: true, NotEmpty: true},
	// {Key: "managers", Creatable: true, Editable: true},
	// {Key: "created_at", NotEmpty: true},
	// {Key: "updated_at"},

	var fieldsToCreateList, fieldsToUpdateList, fieldsToReadList []string
	var wildcardsToCreate []string
	for _, f := range fields {
		//log.Printf("%#v", f)
		fieldsToReadList = append(fieldsToReadList, f.Key)
		if f.Creatable {
			fieldsToCreateList = append(fieldsToCreateList, f.Key)
			wildcardsToCreate = append(wildcardsToCreate, "?")
		}
		if f.Updatable {
			fieldsToUpdateList = append(fieldsToUpdateList, f.Key)
		}
	}

	//log.Fatal(fieldsToUpdateList)

	fieldsToCreate := "`" + strings.Join(fieldsToCreateList, "`, `") + "`"
	fieldsToRead := "`" + strings.Join(fieldsToReadList, "`, `") + "`"
	fieldsToUpdate := "`" + strings.Join(fieldsToUpdateList, "` = ?, `") + "` = ?"

	crudOp := crudMySQL{
		grOp:        grOp,
		managers:    managers,
		dbh:         dbh,
		table:       table,
		fields:      fields,
		sqlCreate:   "insert into `" + table + "` (" + fieldsToCreate + ") values (" + strings.Join(wildcardsToCreate, ",") + ")",
		sqlRead:     "select " + fieldsToRead + " from `" + table + "` where id = ?",
		sqlUpdate:   "update `" + table + "` set " + fieldsToUpdate + " where id = ?",
		sqlDelete:   "delete from `" + table + "` where id = ?",
		sqlReadList: "select SQL_CALC_FOUND_ROWS " + fieldsToRead + " from `" + table + "`",
	}

	sqlStmts := []SqlStmt{
		{&crudOp.stmCreate, crudOp.sqlCreate},
		{&crudOp.stmRead, crudOp.sqlRead},
		{&crudOp.stmUpdate, crudOp.sqlUpdate},
		{&crudOp.stmDelete, crudOp.sqlDelete},
	}

	for _, sqlStmt := range sqlStmts {
		if err = CreateStmt(dbh, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, nil, errors.Wrap(err, onNew)
		}
	}

	return &crudOp, func() error { return crudOp.clean() }, nil
}

func (crudOp *crudMySQL) clean() error {
	_, err := crudOp.dbh.Exec("truncate `" + crudOp.table + "`")
	return err
}

func (crudOp *crudMySQL) Describe() (crud.Description, error) {
	return crud.Description{Fields: crudOp.fields}, nil
}

func (crudOp *crudMySQL) StringMapToNative(data crud.StringMap) (interface{}, error) {
	return crud.StringMapToNativeMap(data, crudOp.fields)
}

func (crudOp *crudMySQL) NativeToStringMap(native interface{}) (crud.StringMap, error) {
	nativeMap, ok := native.(crud.NativeMap)
	if !ok {
		return nil, errors.Wrapf(basis.ErrWrongDataType, "expected crud.NativeMap, actual = %T", native)
	}

	return crud.NativeMapToStringMap(nativeMap, crudOp.fields)
}

const onIDFromNative = "on crudMySQL.IDFromNative()"

func (crudOp *crudMySQL) IDFromNative(native interface{}) (string, error) {
	nativeMap, ok := native.(crud.NativeMap)
	if !ok {
		return "", errors.Wrapf(basis.ErrWrongDataType, onIDFromNative+": expected crud.NativeMap, actual = %T", native)
	}

	id, _ := nativeMap["id"].(string)

	return id, nil
}

const onCreate = "on crudMySQL.Create()"

func (crudOp *crudMySQL) Create(userIS auth.ID, native interface{}) (string, error) {
	nativeMap, ok := native.(crud.NativeMap)
	if !ok {
		return "", errors.Wrapf(basis.ErrWrongDataType, "expected crud.NativeMap, actual = %T", native)
	}

	//if err := groups.OneOfErr(userIS, crudOp.grOp, crudOp.managers[rights.Create]); err != nil {
	//	return "", errors.Wrap(err, onCreate)
	//}

	//rView, rOwner, managers, err := groups.SetRights(userIS, crudOp.grOp, record.RView, record.ROwner, record.Managers)
	//if err != nil {
	//	return "", errors.Wrap(err, onCreate+": can't .SetRights)")
	//}
	//
	//var managersStr []byte
	//if managers != nil {
	//	if managersStr, err = json.Marshal(managers); err != nil {
	//		return "", errors.Wrapf(err, onCreate+": can't json.Marshal($#v)", managers)
	//	}
	//}
	//

	var values []interface{}
	for _, f := range crudOp.fields {
		if f.Creatable {
			values = append(values, nativeMap[f.Key])
		}
	}

	res, err := crudOp.stmCreate.Exec(values...)
	if err != nil {
		return "", errors.Wrapf(err, onCreate+": can't exec SQL: %s, %#v", crudOp.sqlCreate, values)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return "", errors.Wrapf(err, onCreate+": can't get LastInsertId() SQL: %s, %#v", crudOp.sqlCreate, values)
	}

	return strconv.FormatInt(id, 10), nil
}

const onRead = "on crudMySQL.Read()"

func (crudOp *crudMySQL) Read(userIS auth.ID, idStr string) (interface{}, error) {
	if len(idStr) < 1 {
		return nil, errors.Wrap(crud.ErrEmptySelector, onRead)
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return nil, errors.Wrap(crud.ErrBadSelector, onRead+": "+idStr)
	}

	readValuesPtr, err := crud.StringMapToNativePtrList(nil, crudOp.fields)
	if err != nil {
		return nil, errors.Wrap(err, onRead)
	}

	var readValues []interface{}
	for _, rv := range readValuesPtr {
		readValues = append(readValues, rv)
	}

	err = crudOp.stmRead.QueryRow(id).Scan(readValues...)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrapf(err, onRead+": can't exec QueryRow: %s, id = %s", crudOp.sqlRead, id)
	}

	//if !groups.OneOf(userIS, crudOp.grOp, src.RView, src.ROwner) {
	//	return nil, errors.Wrap(basis.ErrNotFound, onRead)
	//}
	//
	//if len(managers) > 0 {
	//	if err = json.Unmarshal(managers, &src.Managers); err != nil {
	//		return nil, errors.Wrapf(err, onRead+": can't json.Unmarshal() for managers field(%s)", managers)
	//	}
	//}

	return crud.NativePtrListToNativeMap(readValuesPtr, crudOp.fields)
}

const onReadList = "on crudMySQL.ReadList()"

func (crudOp *crudMySQL) ReadList(userIS auth.ID, options *content.ListOptions) ([]interface{}, uint64, error) {
	var err error
	var values []interface{}
	var orderAndLimit, condition, conditionCompleted string

	if options != nil {
		condition, values, err = selectors.Mysql(userIS, options.Selector)
		if err != nil {
			return nil, 0, errors.Wrapf(err, ": bad selector ('%#v')", options.Selector)
		}

		conditionCompleted = condition
		if strings.TrimSpace(conditionCompleted) != "" {
			conditionCompleted = " where " + conditionCompleted
		}

		orderAndLimit = OrderAndLimit(options.SortBy, options.Limits)
	}

	var sqlQuery string
	var rows *sql.Rows

	// TODO: check it correct!!!
	if crudOp.grOp != nil {
		sqlQuery, rows, err = groups.QueryAccessible(crudOp.grOp, crudOp.dbh, userIS, crudOp.sqlReadList, condition, orderAndLimit, values)
	} else {
		if strings.TrimSpace(condition) != "" {
			condition = "where " + condition
		}
		sqlQuery = crudOp.sqlReadList + " " + condition + " " + orderAndLimit
		rows, err = crudOp.dbh.Query(sqlQuery, values...)
	}

	if err == sql.ErrNoRows || err == basis.ErrNotFound {
		return nil, 0, nil
	} else if err != nil {
		return nil, 0, errors.Wrapf(err, onReadList+": can't get query (sql='%s', values='%#v')", sqlQuery, values)
	}
	defer rows.Close()

	var items []interface{}

	for rows.Next() {

		readValuesPtr, err := crud.StringMapToNativePtrList(nil, crudOp.fields)
		if err != nil {
			return nil, 0, errors.Wrap(err, onReadList)
		}
		var readValues []interface{}
		for _, rv := range readValuesPtr {
			readValues = append(readValues, rv)
		}

		// var managers []byte
		err = rows.Scan(readValues...)
		if err != nil {
			return items, 0, errors.Wrapf(err, ": can't scan queryRow (sql='%s', values='%#v')", sqlQuery, values)
		}

		//if len(managers) > 0 {
		//	if err = json.Unmarshal(managers, &src.Managers); err != nil {
		//		return items, 0, errors.Wrapf(err, onReadList+": can't json.Unmarshal() for managers field(%s)", managers)
		//	}
		//}

		nativeMap, err := crud.NativePtrListToNativeMap(readValuesPtr, crudOp.fields)
		if err != nil {
			return items, 0, errors.Wrap(err, onReadList)
		}

		items = append(items, nativeMap)
	}
	err = rows.Err()
	if err != nil {
		return items, 0, errors.Wrapf(err, onReadList+": on rows.Err(): sql='%s' (%#v)", sqlQuery, values)
	}

	var allCount uint64
	err = crudOp.dbh.QueryRow("SELECT FOUND_ROWS()").Scan(&allCount)
	if err != nil {
		return nil, 0, errors.Wrapf(err, onReadList+": can't scan ('SELECT FOUND_ROWS()') for sql=%s (%#v)", sqlQuery, values)
	}
	return items, allCount, nil
}

const onUpdate = "on crudMySQL.Update()"

func (crudOp *crudMySQL) Update(userIS auth.ID, native interface{}) (crud.Result, error) {
	nativeMap, ok := native.(crud.NativeMap)
	if !ok {
		return crud.Result{}, errors.Wrapf(basis.ErrWrongDataType, "expected crud.NativeMap, actual = %T", native)
	}

	//record0, err := crudOp.Read(userIS, record.TargetID)
	//if err != nil {
	//	return crud.Result{}, errors.Wrap(err, onUpdate+": can't .Read()")
	//}
	//
	//if !groups.OneOf(userIS, crudOp.grOp, record0.ROwner, record0.Managers[rights.Change]) {
	//	return crud.Result{}, errors.Wrap(basis.ErrNotFound, onUpdate)
	//}
	//
	//rView, rOwner, managers, err := groups.SetRights(userIS, crudOp.grOp, record.RView, record.ROwner, record.Managers)
	//if err != nil {
	//	return crud.Result{}, errors.Wrap(err, onUpdate+": can't .SetRights)")
	//}
	//
	//var managersStr []byte
	//if managers != nil {
	//	if managersStr, err = json.Marshal(managers); err != nil {
	//		return crud.Result{}, errors.Wrapf(err, onUpdate+": can't json.Marshal($#v)", managers)
	//	}
	//}

	var values []interface{}
	for _, f := range crudOp.fields {
		if f.Updatable {
			values = append(values, nativeMap[f.Key])
		}
	}

	id, _ := nativeMap["id"].(string)
	if id == "" {
		return crud.Result{}, errors.Errorf(onUpdate+": no id string in %#v", nativeMap)
	}

	values = append(values, id)

	res, err := crudOp.stmUpdate.Exec(values...)
	if err != nil {
		return crud.Result{}, errors.Wrapf(err, onUpdate+": can't exec SQL: %s, %#v", crudOp.sqlUpdate, values)
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return crud.Result{}, errors.Wrapf(err, onUpdate+": can't get RowsAffected(): %s (%#v)", crudOp.sqlUpdate, values)
	}

	return crud.Result{NumOk: cnt}, nil
}

const onDelete = "on crudMySQL.DeleteList()"

func (crudOp *crudMySQL) Delete(userIS auth.ID, id string) (crud.Result, error) {
	//record0, err := crudOp.Read(userIS, id)
	//if err != nil {
	//	return crud.Result{}, errors.Wrap(err, onDelete+": can't .Read()")
	//}
	//
	//if !groups.OneOf(userIS, crudOp.grOp, record0.ROwner, record0.Managers[rights.Change]) {
	//	return crud.Result{}, errors.Wrap(basis.ErrNotFound, onDelete)
	//}

	res, err := crudOp.stmDelete.Exec(id)
	if err != nil {
		return crud.Result{}, errors.Wrapf(err, onDelete+": can't exec SQL: %s, %s", crudOp.sqlDelete, id)
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return crud.Result{}, errors.Wrapf(err, onDelete+": can't get RowsAffected(): %s, %s", crudOp.sqlDelete, id)
	}
	return crud.Result{cnt}, nil
}

func (crudOp *crudMySQL) Close() error {
	return errors.Wrap(crudOp.dbh.Close(), "on crudMySQL.dbh.Close()")
}
