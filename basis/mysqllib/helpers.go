package mysqllib

import (
	"database/sql"
	"regexp"
	"strconv"
	"strings"

	"github.com/pavlo67/associatio/basis"
	"github.com/pavlo67/associatio/starter/config"
	"github.com/pkg/errors"
)

func CreateStmt(dbh *sql.DB, sql string, stmt **sql.Stmt) (err error) {
	*stmt, err = dbh.Prepare(sql)
	if err != nil {
		return errors.Wrapf(err, "can't prepare (sql: %v)", sql)
	}
	return nil
}

type SqlStmt struct {
	Stmt **sql.Stmt
	Sql  string
}

const MySQLVersionLimit = "05.04"
const MySQLVersionLimit2 = "08.00"
const MySQLVersionFullTextLimit = "05.06"

func ConnectToMysql(mysqlConfig config.ServerAccess) (db *sql.DB, err error) {
	port := ":3306"
	if mysqlConfig.Port != 0 {
		port = ":" + strconv.Itoa(mysqlConfig.Port)
	}

	addr := mysqlConfig.User + ":" + mysqlConfig.Pass +
		"@tcp(" + mysqlConfig.Host + port + ")/" + mysqlConfig.Path + "?parseTime=true"

	dbh, err := sql.Open("mysql", addr)
	if err != nil {
		return nil, errors.Wrapf(err, "wrong db connect (credentials='%v')", mysqlConfig)
	}

	err = dbh.Ping()
	if err != nil {
		return nil, errors.Wrapf(err, "wrong .Ping on db connect (credentials='%v')", mysqlConfig)
	}
	return dbh, nil
}

func TableExists(dbh *sql.DB, table string) error {
	var stmt *sql.Stmt
	sqlQuery := "SHOW TABLES LIKE '" + table + "'"
	if err := CreateStmt(dbh, sqlQuery, &stmt); err != nil {
		return err
	}
	rows, err := stmt.Query()
	if err != nil {
		return errors.Wrapf(err, "can't query SQL: %s", sqlQuery)
	}
	if rows != nil {
		defer rows.Close()
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
	var stmt *sql.Stmt
	sqlQuery := "DROP TABLE IF EXISTS`" + table + "`"
	if err := CreateStmt(dbh, sqlQuery, &stmt); err != nil {
		return err
	}
	_, err := stmt.Exec()
	if err != nil {
		return errors.Wrapf(err, "can't exec SQL: %s", sqlQuery)
	}
	return nil
}

var reVersion = regexp.MustCompile("^\\d+\\.\\d+")

func MySQLVersion(dbh *sql.DB) (string, error) {
	var stmt *sql.Stmt
	sqlQuery := "select VERSION()"
	if err := CreateStmt(dbh, sqlQuery, &stmt); err != nil {
		return "", err
	}
	rows, err := stmt.Query()
	if err != nil {
		return "", errors.Wrapf(err, "can't query SQL: %s", sqlQuery)
	}
	if rows != nil {
		defer rows.Close()
	}

	var v []byte
	for rows.Next() {
		err = rows.Scan(&v)
		if err != nil {
			return "", errors.Wrapf(err, "can't scan SQL: %s", sqlQuery)
		}

		version := reVersion.Find(v)
		if version != nil {
			parts := strings.Split(string(version), ".")
			if len(parts) != 2 {
				return "", errors.Errorf("can't get version from '%s'", version)
			}
			major := parts[0]
			if len(major) < 2 {
				major = "0" + major
			}

			minor := parts[1]
			if len(minor) < 2 {
				minor = "0" + minor
			}

			return major + "." + minor, nil
		}
	}

	return "", errors.New("there is no mysql version???")

}

func QueryStrings(stmt *sql.Stmt, sql string, values ...interface{}) (results []string, err error) {
	rows, err := stmt.Query(values...)
	if err != nil {
		return nil, errors.Wrapf(err, basis.CantExecQuery, sql, values)
	}
	defer rows.Close()

	for rows.Next() {
		var r string
		if err := rows.Scan(&r); err != nil {
			return results, errors.Wrapf(err, CantScanQueryRow, sql, values)
		}

		results = append(results, r)
	}

	err = rows.Err()
	if err != nil {
		return results, errors.Wrapf(err, CantScanQueryRow, sql, values)
	}

	return results, nil
}

func QueryIDs(stmt *sql.Stmt, sql string, values ...interface{}) (ids []uint64, err error) {
	rows, err := stmt.Query(values...)
	if err != nil {
		return nil, errors.Wrapf(err, basis.CantExecQuery, sql, values)
	}
	defer rows.Close()

	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			return ids, errors.Wrapf(err, CantScanQueryRow, sql, values)
		}

		ids = append(ids, id)
	}

	err = rows.Err()
	if err != nil {
		return ids, errors.Wrapf(err, CantScanQueryRow, sql, values)
	}

	return ids, nil
}
