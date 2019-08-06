package flow_sqlite

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/constructor/apps/content"
	"github.com/pavlo67/constructor/apps/flow"
	"github.com/pavlo67/constructor/apps/links"
	"github.com/pavlo67/constructor/basis"
	"github.com/pavlo67/constructor/basis/sqllib"
	"github.com/pavlo67/constructor/structura"
)

var _ flow.Operator = &flowSQLite{}

type flowSQLite struct {
	// grOp     groups.Operator

	dbh *sql.DB

	stmRead, stmListAll, stmListByTag, stmListBySourceID, stmTags, stmSources, stmHas, stmSave, stmRemove *sql.Stmt
	sqlRead, sqlListAll, sqlListByTag, sqlListBySourceID, sqlTags, sqlSources, sqlHas, sqlSave, sqlRemove string
}

const onNew = "on flowSQLite.New()"

//func NewCRUDOperator(grOp groups.Operator, mysqlConfig config.ServerAccess, table string, fields []crud.Field, managers rights.Managers) (crud.Operator, crud.Cleaner, error) {
//	dbh, err := ConnectToMysql(mysqlConfig)
//	if err != nil {
//		return nil, nil, errors.Wrap(err, onNew)
//	}
//
//	if grOp == nil {
//		// DO WARNING, it's may be not ok
//	}
//
//	if strings.TrimSpace(table) == "" {
//		return nil, nil, errors.New(onNew + ": no table name defined")
//	}
//
//	if len(fields) < 1 {
//		return nil, nil, errors.New(onNew + ": no table fields defined")
//	}
//
//	// TODO: customize special fields - id, created_at, updated_at, r_view, r_owner, managers
//
//	// {Key: "id", Unique: true, AutoUnique: true},
//	// {Key: "r_view", Creatable: true, Editable: true, NotEmpty: true},
//	// {Key: "r_owner", Creatable: true, Editable: true, NotEmpty: true},
//	// {Key: "managers", Creatable: true, Editable: true},
//	// {Key: "created_at", NotEmpty: true},
//	// {Key: "updated_at"},
//
//	var fieldsToCreateList, fieldsToUpdateList, fieldsToReadList []string
//	var wildcardsToCreate []string
//	for _, f := range fields {
//		//log.Printf("%#v", f)
//		fieldsToReadList = append(fieldsToReadList, f.Key)
//		if f.Creatable {
//			fieldsToCreateList = append(fieldsToCreateList, f.Key)
//			wildcardsToCreate = append(wildcardsToCreate, "?")
//		}
//		if f.Updatable {
//			fieldsToUpdateList = append(fieldsToUpdateList, f.Key)
//		}
//	}
//
//	//log.Fatal(fieldsToUpdateList)
//
//	fieldsToCreate := "`" + strings.Join(fieldsToCreateList, "`, `") + "`"
//	fieldsToRead := "`" + strings.Join(fieldsToReadList, "`, `") + "`"
//	fieldsToUpdate := "`" + strings.Join(fieldsToUpdateList, "` = ?, `") + "` = ?"
//
//	crudOp := flowSQLite{
//		grOp:        grOp,
//		managers:    managers,
//		dbh:         dbh,
//		table:       table,
//		fields:      fields,
//		sqlCreate:   "insert into `" + table + "` (" + fieldsToCreate + ") values (" + strings.Join(wildcardsToCreate, ",") + ")",
//		sqlRead:     "select " + fieldsToRead + " from `" + table + "` where id = ?",
//		sqlUpdate:   "update `" + table + "` set " + fieldsToUpdate + " where id = ?",
//		sqlDelete:   "delete from `" + table + "` where id = ?",
//		sqlReadList: "select SQL_CALC_FOUND_ROWS " + fieldsToRead + " from `" + table + "`",
//	}
//
//	sqlStmts := []SqlStmt{
//		{&crudOp.stmCreate, crudOp.sqlCreate},
//		{&crudOp.stmRead, crudOp.sqlRead},
//		{&crudOp.stmUpdate, crudOp.sqlUpdate},
//		{&crudOp.stmDelete, crudOp.sqlDelete},
//	}
//
//	for _, sqlStmt := range sqlStmts {
//		if err = Exec(dbh, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
//			return nil, nil, errors.Wrap(err, onNew)
//		}
//	}
//
//	return &crudOp, func() error { return crudOp.clean() }, nil
//}
//
//func (crudOp *flowSQLite) clean() error {
//	_, err := crudOp.dbh.Exec("truncate `" + crudOp.table + "`")
//	return err
//}

const onSources = "on flowSQLite.Sources(): "

func (flowOp *flowSQLite) Sources(_ *structura.GetOptions) ([]basis.ID, error) {
	rows, err := flowOp.stmSources.Query()
	if err != nil {
		return nil, errors.Errorf(onSources+sqllib.CantQuery, flowOp.sqlSources, nil)
	}
	defer rows.Close()

	var sourceIDs []basis.ID
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			return nil, errors.Errorf(onSources+sqllib.CantScanQueryRow, flowOp.sqlSources, nil)
		}

		sourceIDs = append(sourceIDs, basis.ID(id))
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Errorf(onSources+sqllib.CantScanQueryRow, flowOp.sqlSources, nil)

	}

	return sourceIDs, nil
}

func (flowOp *flowSQLite) Tags(*structura.GetOptions) ([]links.Tag, error) {
	return nil, basis.ErrNotImplemented
}

func (flowOp *flowSQLite) ListAll(before *time.Time, options *structura.GetOptions) ([]content.Brief, error) {
	return nil, basis.ErrNotImplemented
}

func (flowOp *flowSQLite) ListBySourceID(sourceID basis.ID, before *time.Time, options *structura.GetOptions) ([]content.Brief, error) {
	return nil, basis.ErrNotImplemented
}

func (flowOp *flowSQLite) ListByTag(tag string, before *time.Time, options *structura.GetOptions) ([]content.Brief, error) {
	return nil, basis.ErrNotImplemented
}

func (flowOp *flowSQLite) Read(basis.ID, *structura.GetOptions) (*flow.Item, error) {
	return nil, basis.ErrNotImplemented
}

func (flowOp *flowSQLite) Has(*flow.Source) (bool, error) {
	return false, basis.ErrNotImplemented
}

func (flowOp *flowSQLite) Save([]flow.Item, *structura.SaveOptions) ([]basis.ID, error) {
	return nil, basis.ErrNotImplemented
}

func (flowOp *flowSQLite) Remove(sourceIDs []basis.ID, before *time.Time, options *structura.RemoveOptions) error {
	return basis.ErrNotImplemented
}

func (flowOp *flowSQLite) Close() error {
	return basis.ErrNotImplemented
}

func (flowOp *flowSQLite) Clean() error {
	return basis.ErrNotImplemented

}

//const onCreate = "on flowSQLite.Create()"
//
//func (crudOp *flowSQLite) Create(userIS auth.ID, native interface{}) (string, error) {
//	nativeMap, ok := native.(crud.NativeMap)
//	if !ok {
//		return "", errors.Wrapf(basis.ErrWrongDataType, "expected crud.NativeMap, actual = %T", native)
//	}
//
//	//if err := groups.OneOfErr(userIS, crudOp.grOp, crudOp.managers[rights.Create]); err != nil {
//	//	return "", errors.Wrap(err, onCreate)
//	//}
//
//	//rView, rOwner, managers, err := groups.SetRights(userIS, crudOp.grOp, record.RView, record.ROwner, record.Managers)
//	//if err != nil {
//	//	return "", errors.Wrap(err, onCreate+": can't .SetRights)")
//	//}
//	//
//	//var managersStr []byte
//	//if managers != nil {
//	//	if managersStr, err = json.Marshal(managers); err != nil {
//	//		return "", errors.Wrapf(err, onCreate+": can't json.Marshal($#v)", managers)
//	//	}
//	//}
//	//
//
//	var values []interface{}
//	for _, f := range crudOp.fields {
//		if f.Creatable {
//			values = append(values, nativeMap[f.Key])
//		}
//	}
//
//	res, err := crudOp.stmCreate.Exec(values...)
//	if err != nil {
//		return "", errors.Wrapf(err, onCreate+": can't exec SQL: %s, %#v", crudOp.sqlCreate, values)
//	}
//
//	id, err := res.LastInsertId()
//	if err != nil {
//		return "", errors.Wrapf(err, onCreate+": can't get LastInsertId() SQL: %s, %#v", crudOp.sqlCreate, values)
//	}
//
//	return strconv.FormatInt(id, 10), nil
//}
//
//const onRead = "on flowSQLite.Read()"
//
//func (crudOp *flowSQLite) Read(userIS auth.ID, idStr string) (interface{}, error) {
//	if len(idStr) < 1 {
//		return nil, errors.Wrap(crud.ErrEmptySelector, onRead)
//	}
//
//	id, err := strconv.ParseUint(idStr, 10, 64)
//	if err != nil {
//		return nil, errors.Wrap(crud.ErrBadSelector, onRead+": "+idStr)
//	}
//
//	readValuesPtr, err := crud.StringMapToNativePtrList(nil, crudOp.fields)
//	if err != nil {
//		return nil, errors.Wrap(err, onRead)
//	}
//
//	var readValues []interface{}
//	for _, rv := range readValuesPtr {
//		readValues = append(readValues, rv)
//	}
//
//	err = crudOp.stmRead.QueryRow(id).Scan(readValues...)
//	if err == sql.ErrNoRows {
//		return nil, nil
//	}
//	if err != nil {
//		return nil, errors.Wrapf(err, onRead+": can't exec QueryRow: %s, id = %s", crudOp.sqlRead, id)
//	}
//
//	//if !groups.OneOf(userIS, crudOp.grOp, src.RView, src.ROwner) {
//	//	return nil, errors.Wrap(basis.ErrNotFound, onRead)
//	//}
//	//
//	//if len(managers) > 0 {
//	//	if err = json.Unmarshal(managers, &src.Managers); err != nil {
//	//		return nil, errors.Wrapf(err, onRead+": can't json.Unmarshal() for managers field(%s)", managers)
//	//	}
//	//}
//
//	return crud.NativePtrListToNativeMap(readValuesPtr, crudOp.fields)
//}
//
//const onReadList = "on flowSQLite.ReadList()"
//
//func (crudOp *flowSQLite) ReadList(userIS auth.ID, options *content.ListOptions) ([]interface{}, uint64, error) {
//	var err error
//	var values []interface{}
//	var orderAndLimit, condition, conditionCompleted string
//
//	if options != nil {
//		condition, values, err = selectors.Mysql(userIS, options.Selector)
//		if err != nil {
//			return nil, 0, errors.Wrapf(err, ": bad selector ('%#v')", options.Selector)
//		}
//
//		conditionCompleted = condition
//		if strings.TrimSpace(conditionCompleted) != "" {
//			conditionCompleted = " where " + conditionCompleted
//		}
//
//		orderAndLimit = OrderAndLimit(options.SortBy, options.Limits)
//	}
//
//	var sqlQuery string
//	var rows *sql.Rows
//
//	// TODO: check it correct!!!
//	if crudOp.grOp != nil {
//		sqlQuery, rows, err = groups.QueryAccessible(crudOp.grOp, crudOp.dbh, userIS, crudOp.sqlReadList, condition, orderAndLimit, values)
//	} else {
//		if strings.TrimSpace(condition) != "" {
//			condition = "where " + condition
//		}
//		sqlQuery = crudOp.sqlReadList + " " + condition + " " + orderAndLimit
//		rows, err = crudOp.dbh.Query(sqlQuery, values...)
//	}
//
//	if err == sql.ErrNoRows || err == basis.ErrNotFound {
//		return nil, 0, nil
//	} else if err != nil {
//		return nil, 0, errors.Wrapf(err, onReadList+": can't get query (sql='%s', values='%#v')", sqlQuery, values)
//	}
//	defer rows.Close()
//
//	var items []interface{}
//
//	for rows.Next() {
//
//		readValuesPtr, err := crud.StringMapToNativePtrList(nil, crudOp.fields)
//		if err != nil {
//			return nil, 0, errors.Wrap(err, onReadList)
//		}
//		var readValues []interface{}
//		for _, rv := range readValuesPtr {
//			readValues = append(readValues, rv)
//		}
//
//		// var managers []byte
//		err = rows.Scan(readValues...)
//		if err != nil {
//			return items, 0, errors.Wrapf(err, ": can't scan queryRow (sql='%s', values='%#v')", sqlQuery, values)
//		}
//
//		//if len(managers) > 0 {
//		//	if err = json.Unmarshal(managers, &src.Managers); err != nil {
//		//		return items, 0, errors.Wrapf(err, onReadList+": can't json.Unmarshal() for managers field(%s)", managers)
//		//	}
//		//}
//
//		nativeMap, err := crud.NativePtrListToNativeMap(readValuesPtr, crudOp.fields)
//		if err != nil {
//			return items, 0, errors.Wrap(err, onReadList)
//		}
//
//		items = append(items, nativeMap)
//	}
//	err = rows.Err()
//	if err != nil {
//		return items, 0, errors.Wrapf(err, onReadList+": on rows.Err(): sql='%s' (%#v)", sqlQuery, values)
//	}
//
//	var allCount uint64
//	err = crudOp.dbh.QueryRow("SELECT FOUND_ROWS()").Scan(&allCount)
//	if err != nil {
//		return nil, 0, errors.Wrapf(err, onReadList+": can't scan ('SELECT FOUND_ROWS()') for sql=%s (%#v)", sqlQuery, values)
//	}
//	return items, allCount, nil
//}
//
//func New(mysqlConfig config.ServerAccess, table, tableStaged string, contentTemplate interface{}) (*datastoreMySQL, error) {
//
//	if strings.TrimSpace(table) == "" {
//		return nil, errors.New(onNew + ": no table name is defined")
//	}
//
//	if contentTemplate == nil {
//		l.Warn(onNew + ": no contentTemplate is defined")
//	}
//
//	dbh, err := mysqllib.ConnectToMysql(mysqlConfig)
//	if err != nil {
//		return nil, errors.Wrap(err, onNew)
//	}
//
//	tableStaged = strings.TrimSpace(tableStaged)
//
//	fieldsToAdd := "`" + strings.Join(fields, "`, `") + "`"
//	fieldsToRead := "id," + fieldsToAdd + ", stored_at"
//
//	dsOp := datastoreMySQL{
//		dbh:         dbh,
//		table:       table,
//		tableStaged: tableStaged,
//
//		contentTemplate: contentTemplate,
//
//		sqlAdd:       "insert ignore into `" + table + "` (" + fieldsToAdd + ") values (?,?,?, ?,?,?, ?,?,?)",
//		sqlReadList:  "select SQL_CALC_FOUND_ROWS " + fieldsToRead + " from `" + table + "`",
//		sqlDelete:    "delete from `" + table + "`",
//		sqlLastKey:   "select max(content_key) from `" + table + "` where content_type = ?",
//		sqlKeyExists: "select id from `" + table + "` where content_type = ? and content_key = ?",
//
//		sqlAddStaged:    "insert ignore into `" + tableStaged + "` (" + fieldsToAdd + ") values (?,?,?, ?,?,?, ?,?,?)",
//		sqlMarkStaged:   "update `" + tableStaged + "` set `status` = ? ",
//		sqlCommitStaged: "insert into `" + table + "` (" + fieldsToAdd + ") select (" + fieldsToAdd + ") from `" + tableStaged + "` where `status` = ?",
//		sqlGetStaged:    "select SQL_CALC_FOUND_ROWS " + fieldsToRead + " from `" + table + "`",
//		sqlDeleteStaged: "delete from `" + tableStaged + "`",
//	}
//
//	sqlStmts := []mysqllib.SqlStmt{
//		{&dsOp.stmAdd, dsOp.sqlAdd},
//		{&dsOp.stmLastKey, dsOp.sqlLastKey},
//		{&dsOp.stmKeyExists, dsOp.sqlKeyExists},
//	}
//
//	if tableStaged != "" {
//		sqlStmts = append(
//			sqlStmts,
//			mysqllib.SqlStmt{&dsOp.stmAddStaged, dsOp.sqlAddStaged},
//			mysqllib.SqlStmt{&dsOp.stmCommitStaged, dsOp.sqlCommitStaged},
//		)
//	}
//
//	for _, sqlStmt := range sqlStmts {
//		if err = mysqllib.CreateStmt(dbh, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
//			return nil, errors.Wrap(err, onNew)
//		}
//	}
//
//	return &dsOp, nil
//}
//
//func (dsOp *datastoreMySQL) Clean() error {
//	_, err1 := dsOp.dbh.Exec("truncate `" + dsOp.table + "`")
//
//	var err2 error
//	if dsOp.tableStaged != "" {
//		_, err2 = dsOp.dbh.Exec("truncate `" + dsOp.tableStaged + "`")
//	}
//
//	return basis.MultiError(err1, err2).Err()
//}
//
//const onAdd = "on datastoreMySQL.Save()"
//
//func (dsOp *datastoreMySQL) Save(item *flow.Item) error {
//	if item == nil {
//		return errors.Wrap(basis.ErrNullItem, onAdd)
//	}
//
//	var content []byte
//	if item.Content != nil {
//		switch v := item.Content.(type) {
//		case string:
//			content = []byte(v)
//		case *string:
//			if v != nil {
//				content = []byte(*v)
//			}
//		case []byte:
//			content = v
//		case *[]byte:
//			if v != nil {
//				content = *v
//			}
//		default:
//			var err error
//			if content, err = json.Marshal(item.Content); err != nil {
//				return errors.Wrapf(err, onAdd+": can't json.Marshal($#v)", item.Content)
//			}
//		}
//	}
//
//	values := []interface{}{
//		item.Source.URL,
//		item.Source.Key,
//		item.Source.Time,
//		item.Original,
//		string(item.ContentType),
//		item.ContentKey,
//		content,
//		item.Status,
//		item.History,
//	}
//
//	_, err := dsOp.stmAdd.Exec(values...)
//	if err != nil {
//		return errors.Wrapf(err, onAdd+": can't exec SQL: %s, %#v", dsOp.sqlAdd, values)
//	}
//
//	return nil
//}
//
//const onKeyExists = "on datastoreMySQL.KeyExists()"
//
//func (dsOp *datastoreMySQL) KeyExists(class flow.Type, key string) (bool, error) {
//	values := []interface{}{string(class), key}
//	rows, err := dsOp.stmKeyExists.Query(values...)
//	if err == sql.ErrNoRows {
//		return false, nil
//	} else if err != nil {
//		return false, errors.Wrapf(err, onKeyExists+": can't query (sql='%s', values='%#v')", dsOp.sqlKeyExists, values)
//	}
//	defer rows.Close()
//	if rows.Next() {
//		return true, nil
//	}
//	err = rows.Err()
//	if err != nil {
//		return false, errors.Wrapf(err, onKeyExists+": on rows.Err() (sql='%s', values='%#v')", dsOp.sqlKeyExists, values)
//	}
//
//	return false, nil
//}
//
//const onLastKey = "on datastoreMySQL.LastKey()"
//
//func (dsOp *datastoreMySQL) LastKey(class flow.Type, options *crud.ReadOptions) (string, error) {
//
//	// TODO: use options!!!
//
//	values := []interface{}{string(class)}
//	rows, err := dsOp.stmLastKey.Query(values...)
//	if err == sql.ErrNoRows {
//		return "", nil
//	} else if err != nil {
//		return "", errors.Wrapf(err, onLastKey+": can't query (sql='%s', values='%#v')", dsOp.sqlLastKey, values)
//	}
//	defer rows.Close()
//	if rows.Next() {
//		var lastKey string
//		err = rows.Scan(&lastKey)
//		if err != nil {
//			return "", errors.Wrapf(err, onLastKey+": can't scan query row (sql='%s', values='%#v')", dsOp.sqlLastKey, values)
//		}
//		return lastKey, nil
//	}
//	err = rows.Err()
//	if err != nil {
//		return "", errors.Wrapf(err, onLastKey+": on rows.Err() (sql='%s', values='%#v')", dsOp.sqlLastKey, values)
//	}
//
//	return "", nil
//}
//
//const onGet = "on datastoreMySQL.ReadList()"
//
//func (dsOp *datastoreMySQL) ReadList(options *crud.ReadOptions) ([]flow.Item, uint64, error) {
//	var err error
//	var values []interface{}
//	var orderAndLimit, condition, conditionCompleted string
//
//	if options != nil {
//		condition, values, err = selectors.Mysql("", options.Selector)
//		if err != nil {
//			return nil, 0, errors.Wrapf(err, onGet+": bad selector ('%#v')", options.Selector)
//		}
//
//		conditionCompleted = condition
//		if strings.TrimSpace(conditionCompleted) != "" {
//			conditionCompleted = " where " + conditionCompleted
//		}
//
//		orderAndLimit = mysqllib.OrderAndLimit(options.SortBy, options.Limits)
//	}
//	if strings.TrimSpace(condition) != "" {
//		condition = "where " + condition
//	}
//
//	// log.Fatal(condition, values)
//
//	sqlQuery := dsOp.sqlReadList + " " + condition + " " + orderAndLimit
//	rows, err := dsOp.dbh.Query(sqlQuery, values...)
//	if err == sql.ErrNoRows || err == basis.ErrNotFound {
//		return nil, 0, nil
//	} else if err != nil {
//		return nil, 0, errors.Wrapf(err, onGet+": can't get query (sql='%s', values='%#v')", sqlQuery, values)
//	}
//	defer rows.Close()
//
//	var items []flow.Item
//	for rows.Next() {
//		var item flow.Item
//		var contentBytes []byte
//		err = rows.Scan(
//			&item.ID,
//			&item.Source.URL,
//			&item.Source.Key,
//			&item.Source.Time,
//			&item.Original,
//			&item.ContentType,
//			&item.ContentKey,
//			&contentBytes,
//			&item.Status,
//			&item.History,
//			&item.StoredAt)
//		if err != nil {
//			return items, 0, errors.Wrapf(err, ": can't scan queryRow (sql='%s', values='%#v')", sqlQuery, values)
//		}
//
//		if len(contentBytes) > 0 {
//			if dsOp.contentTemplate == nil {
//				item.ContentStr = string(contentBytes)
//			} else {
//				content := reflect.New(reflect.ValueOf(dsOp.contentTemplate).Elem().Type()).Interface()
//				if err = json.Unmarshal(contentBytes, item.Content); err != nil {
//					l.Error(onGet+": can't json.Unmarshal() for content field(%s): %s", contentBytes, err)
//					item.ContentStr = string(contentBytes)
//				} else {
//					item.Content = content
//				}
//			}
//		}
//
//		items = append(items, item)
//	}
//	err = rows.Err()
//	if err != nil {
//		return items, 0, errors.Wrapf(err, onGet+": on rows.Err(): sql='%s' (%#v)", sqlQuery, values)
//	}
//
//	var allCount uint64
//	err = dsOp.dbh.QueryRow("SELECT FOUND_ROWS()").Scan(&allCount)
//	if err != nil {
//		return nil, 0, errors.Wrapf(err, onGet+": can't scan ('SELECT FOUND_ROWS()') for sql=%s (%#v)", sqlQuery, values)
//	}
//	return items, allCount, nil
//}
//
//const onDelete = "on datastoreMySQL.Delete()"
//
//func (dsOp *datastoreMySQL) Delete(options *crud.ReadOptions) (crud.Result, error) {
//	var err error
//	var values []interface{}
//	var orderAndLimit, condition, conditionCompleted string
//
//	if options != nil {
//		condition, values, err = selectors.Mysql("", options.Selector)
//		if err != nil {
//			return crud.Result{}, errors.Wrapf(err, onDelete+": bad selector ('%#v')", options.Selector)
//		}
//
//		conditionCompleted = condition
//		if strings.TrimSpace(conditionCompleted) != "" {
//			conditionCompleted = " where " + conditionCompleted
//		}
//
//		orderAndLimit = mysqllib.OrderAndLimit(options.SortBy, options.Limits)
//	}
//
//	if strings.TrimSpace(condition) != "" {
//		condition = "where " + condition
//	}
//
//	sqlQuery := dsOp.sqlDelete + " " + condition + " " + orderAndLimit
//	res, err := dsOp.dbh.Exec(sqlQuery, values...)
//	if err != nil {
//		return crud.Result{}, errors.Wrapf(err, onDelete+": can't exec SQL: %s, %s", sqlQuery, values)
//	}
//	cnt, err := res.RowsAffected()
//	if err != nil {
//		return crud.Result{}, errors.Wrapf(err, onDelete+": can't get RowsAffected(): %s, %s", sqlQuery, values)
//	}
//	return crud.Result{cnt}, nil
//}
//
//func (dsOp *datastoreMySQL) Close() error {
//	return errors.Wrap(dsOp.dbh.Close(), "on datastoreMySQL.dbh.Close()")
//}
//
////const onCommit = "on datastoreMySQL.Commit()"
////
////func (dsOp *datastoreMySQL) Commit(options *crud.ReadOptions) error {
////	sql := "update `" + dsOp.tableStaged + "` set status = ?"
////	values := []interface{}{"confirmed " + time.Now().Format(time.RFC3339)}
////	_, err := dsOp.dbh.Exec(sql, values...)
////	if err != nil {
////		return errors.Wrapf(err, onCommit+": can't exec SQL: %s, %#v", sql, values)
////	}
////
////	_, err = dsOp.stmCommit.Exec(values...)
////	if err != nil {
////		return errors.Wrapf(err, onCommit+": can't exec SQL: %s, %#v", dsOp.sqlCommit, values)
////	}
////
////	return nil
////}
////
