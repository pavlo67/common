package flow_sqlite

//import (
//	"database/sql"
//	"encoding/json"
//	"reflect"
//	"strings"
//
//	_ "github.com/go-sql-driver/mysql"
//	"github.com/pavlo67/constructor/basis/mysqllib"
//	"github.com/pavlo67/constructor/starter/config"
//	"github.com/pkg/errors"
//	"github.com/pavlo67/constructor/basis"
//	"github.com/pavlo67/constructor/apps/flow"
//	"github.com/pavlo67/partes/crud"
//)
//
//type datastoreMySQL struct {
//	dbh         *sql.DB
//	table       string
//	tableStaged string
//
//	contentTemplate interface{}
//
//	stmAdd, stmLastKey, stmKeyExists                         *sql.Stmt
//	sqlAdd, sqlLastKey, sqlKeyExists, sqlReadList, sqlDelete string
//
//	stmAddStaged, stmCommitStaged                                               *sql.Stmt
//	sqlAddStaged, sqlCommitStaged, sqlMarkStaged, sqlGetStaged, sqlDeleteStaged string
//}
//
//var fields = []string{
//	"source_url", "source_key", "source_time", "original",
//	"content_type", "content_key", "content",
//	"status", "history",
//}
//
//const onNew = "on datamysql.New()"
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
