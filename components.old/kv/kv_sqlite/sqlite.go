package kv_sqlite

import (
	"database/sql"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/kv"
	"github.com/pavlo67/workshop/libraries/sqllib"
	"github.com/pavlo67/workshop/libraries/sqllib/sqllib_sqlite"
)

type kvSQLite struct {
	db    *sql.DB
	table string

	stmSet, stmGet *sql.Stmt
	sqlSet, sqlGet string
}

var fields = "key, value, saved_at"

const onNew = "on kvSQLite.New()"

func New(sqliteConfig config.ServerAccess, table string) (*kvSQLite, error) {

	sqliteOp, err := sqllib_sqlite.New(sqliteConfig)
	if err != nil {
		return nil, errors.Wrap(err, onNew)
	}

	db := sqliteOp.DB()
	if db == nil {
		return nil, errors.New(onNew + ": no db connector")
	}

	if strings.TrimSpace(table) == "" {
		return nil, errors.New(onNew + ": no table name defined")
	}

	kvOp := kvSQLite{
		db:     db,
		table:  table,
		sqlSet: "REPLACE INTO " + table + " (" + fields + ") VALUES (?,?,strftime('%Y-%m-%d %H-%M-%S','now'))",
		sqlGet: "SELECT " + fields + " FROM " + table + " WHERE key = ?",
	}

	sqlStmts := []sqllib.SqlStmt{
		{&kvOp.stmSet, kvOp.sqlSet},
		{&kvOp.stmGet, kvOp.sqlGet},
	}

	for _, sqlStmt := range sqlStmts {
		if err := sqllib.Prepare(db, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, errors.Wrap(err, onNew)
		}
	}

	return &kvOp, nil
}

const onSet = "on kvSQLite.Set()"

func (kvOp *kvSQLite) Set(key, value string) error {

	values := []interface{}{key, value}
	_, err := kvOp.stmSet.Exec(values...)
	if err != nil {
		return errors.Wrapf(err, onSet+": can't exec SQL: %s, %#v", kvOp.sqlSet, values)
	}

	return nil
}

const onGet = "on kvSQLite.Get()"

func (kvOp *kvSQLite) Get(key string) (*kv.Item, error) {
	var value string
	var storedAt time.Time
	err := kvOp.stmGet.QueryRow(key).Scan(&value, &storedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrapf(err, onGet+": can't exec QueryRow: %s, key = %s", kvOp.sqlGet, key)
	}

	return &kv.Item{
		Key:      key,
		Value:    value,
		StoredAt: storedAt,
	}, nil
}

func (kvOp *kvSQLite) Close() error {
	return errors.Wrap(kvOp.db.Close(), "on kvSQLite.db.Close()")
}

//func (kvOp *kvSQLite) clean() error {
//	_, err := kvOp.dbh.Exec("truncate `" + kvOp.table + "`")
//	return err
//}
