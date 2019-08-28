package kvmysql

import (
	"database/sql"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"github.com/pavlo67/partes/libs/mysqllib"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"
)

type kvMySQL struct {
	dbh   *sql.DB
	table string

	stmSet, stmGet *sql.Stmt
	sqlSet, sqlGet string
}

var fields = []string{"data_type", "key", "value"}

const onNew = "on kvmysql.New()"

func New(mysqlConfig config.ServerAccess, table string) (*kvMySQL, error) {
	dbh, err := mysqllib.ConnectToMysql(mysqlConfig)
	if err != nil {
		return nil, errors.Wrap(err, onNew)
	}

	if strings.TrimSpace(table) == "" {
		return nil, errors.New(onNew + ": no table name defined")
	}

	fieldsToSet := "`" + strings.Join(fields, "`, `") + "`"

	kvOp := kvMySQL{
		dbh:    dbh,
		table:  table,
		sqlSet: "replace into `" + table + "` (" + fieldsToSet + ") values (?,?,?)",
		sqlGet: "select `value`, stored_at from `" + table + "` where data_type = ? and `key` = ?",
	}

	sqlStmts := []mysqllib.SqlStmt{
		{&kvOp.stmSet, kvOp.sqlSet},
		{&kvOp.stmGet, kvOp.sqlGet},
	}

	for _, sqlStmt := range sqlStmts {
		if err = mysqllib.CreateStmt(dbh, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, errors.Wrap(err, onNew)
		}
	}

	return &kvOp, nil
}

const onSet = "on kvMySQL.Set()"

func (kvOp *kvMySQL) Set(dataType joiner.InterfaceKey, key, value string) error {
	values := []interface{}{string(dataType), key, value}
	_, err := kvOp.stmSet.Exec(values...)
	if err != nil {
		return errors.Wrapf(err, onSet+": can't exec SQL: %s, %#v", kvOp.sqlSet, values)
	}

	return nil
}

//func (kvOp *kvMySQL) clean() error {
//	_, err := kvOp.dbh.Exec("truncate `" + kvOp.table + "`")
//	return err
//}

const onGet = "on kvMySQL.ReadList()"

func (kvOp *kvMySQL) Get(dataType joiner.InterfaceKey, key string) (*string, *time.Time, error) {
	var value string
	var storedAt time.Time
	err := kvOp.stmGet.QueryRow(string(dataType), key).Scan(&value, &storedAt)
	if err == sql.ErrNoRows {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, errors.Wrapf(err, onGet+": can't exec QueryRow: %s, dataType = %s, key = %s", kvOp.sqlGet, dataType, key)
	}

	return &value, &storedAt, nil
}

func (kvOp *kvMySQL) Close() error {
	return errors.Wrap(kvOp.dbh.Close(), "on kvMySQL.dbh.Close()")
}
