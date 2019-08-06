package sqllib

import (
	"database/sql"
	"strconv"

	"github.com/pkg/errors"

	"github.com/pavlo67/constructor/starter/config"
)

type Operator interface {
	DB() (*sql.DB, error)
	Connect(cfg config.ServerAccess) error
	CreateSQLQuery(table config.SQLTable) (string, error)
}

// helpers ------------------------------------------------------------

const CantQuery = "can't .Query (sql='%s', values='%#v')"
const CantExec = "can't .Exec (sql='%s', values='%#v')"
const CantGetRowsAffected = "can't .RowsAffected (sql='%s', values='%#v')"
const NoRowOnQuery = "no row on query(sql='%s', values='%#v'"
const CantScanQueryRow = "can't scan query row (sql='%s', values='%#v')"

var ErrNoTable = errors.New("table doesn't exist")

const defaultPageLengthStr = "200"

func OrderAndLimit(sortBy []string, limits []uint64) string {
	var sortStr, limitsStr string
	if len(sortBy) > 0 {
		for _, s := range sortBy {
			if s == "" {
				continue
			}
			desc := ""
			if s[len(s)-1:] == "-" {
				s = s[:len(s)-1]
				desc = " desc"
			} else if s[len(s)-1:] == "+" {
				s = s[:len(s)-1]
			}
			if sortStr != "" {
				sortStr += ", "
			}
			sortStr += "`" + s + "`" + desc
		}
		if sortStr != "" {
			sortStr = " order by " + sortStr
		}
	}
	if len(limits) > 1 {
		// limit[0] can be equal to 0
		var pageLengthStr string
		if limits[1] > 0 {
			pageLengthStr = strconv.FormatUint(limits[1], 10)
		} else {
			pageLengthStr = defaultPageLengthStr
		}
		limitsStr = " limit " + strconv.FormatUint(limits[0], 10) + ", " + pageLengthStr
	} else if len(limits) > 0 {
		if limits[0] > 0 {
			limitsStr = " limit " + strconv.FormatUint(limits[0], 10)
		} else {
			limitsStr = " limit " + defaultPageLengthStr
		}
	}
	return sortStr + limitsStr
}

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
