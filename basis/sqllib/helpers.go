package sqllib

import (
	"database/sql"

	"github.com/pkg/errors"
)

type SqlStmt struct {
	Stmt **sql.Stmt
	Sql  string
}

func Exec(dbh *sql.DB, sqlQuery string) (*sql.Result, error) {
	stmt, err := dbh.Prepare(sqlQuery)
	if err != nil {
		return nil, errors.Wrapf(err, "can't prepare (sqlQuery: %v)", sqlQuery)
	}

	res, err := stmt.Exec()
	if err != nil {
		return nil, errors.Wrapf(err, "can't exec SQL: %s", sqlQuery)
	}

	return &res, nil
}

func Query(dbh *sql.DB, sqlQuery string) (*sql.Rows, error) {
	stmt, err := dbh.Prepare(sqlQuery)
	if err != nil {
		return nil, errors.Wrapf(err, "can't prepare (sqlQuery: %v)", sqlQuery)
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, errors.Wrapf(err, "can't query SQL: %s", sqlQuery)
	}

	return rows, nil
}

func TableExists(dbh *sql.DB, table string) error {
	sqlQuery := "SHOW TABLES LIKE '" + table + "'"

	rows, err := Query(dbh, sqlQuery)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return err
	}

	var t string
	for rows.Next() {
		err = rows.Scan(&t)
		if err != nil {
			return errors.Wrapf(err, "can't scan SQL: %s", sqlQuery)
		}
		return nil
	}

	return ErrNoTable
}

func DropTable(dbh *sql.DB, table string) error {
	sqlQuery := "DROP TABLE IF EXISTS`" + table + "`"

	_, err := Exec(dbh, sqlQuery)

	if err != nil {
		return err
	}
	return nil
}

//func QueryStrings(stmt *sql.Stmt, sql string, values ...interface{}) (results []string, err error) {
//	rows, err := stmt.Query(values...)
//	if err != nil {
//		return nil, errors.Wrapf(err, basis.CantExecQuery, sql, values)
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var r string
//		if err := rows.Scan(&r); err != nil {
//			return results, errors.Wrapf(err, CantScanQueryRow, sql, values)
//		}
//
//		results = append(results, r)
//	}
//
//	err = rows.Err()
//	if err != nil {
//		return results, errors.Wrapf(err, CantScanQueryRow, sql, values)
//	}
//
//	return results, nil
//}

//func QueryIDs(stmt *sql.Stmt, sql string, values ...interface{}) (ids []uint64, err error) {
//	rows, err := stmt.Query(values...)
//	if err != nil {
//		return nil, errors.Wrapf(err, basis.CantExecQuery, sql, values)
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var id uint64
//		if err := rows.Scan(&id); err != nil {
//			return ids, errors.Wrapf(err, CantScanQueryRow, sql, values)
//		}
//
//		ids = append(ids, id)
//	}
//
//	err = rows.Err()
//	if err != nil {
//		return ids, errors.Wrapf(err, CantScanQueryRow, sql, values)
//	}
//
//	return ids, nil
//}
