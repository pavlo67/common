package sqllib

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/errata"
	"github.com/pavlo67/common/common/selectors"
	"github.com/pavlo67/common/common/selectors/logic"
	"github.com/pavlo67/common/common/selectors/selectors_sql"
)

const CantPrepare = "can't .Prepare(%s)"
const CantQuery = "can't .Query('%s', %#v)"
const CantExec = "can't .Exec('%s', %#v)"

const CantGetLastInsertId = "can't .LastInsertId('%s', %#v)"
const CantGetRowsAffected = "can't .RowsAffected('%s', %#v)"
const NoRowOnQuery = "no row on query('%s', %#v)"
const CantScanQueryRow = "can't scan query row ('%s', %#v)"
const RowsError = "error on .Rows ('%s', %#v)"

var ErrNoTable = errata.New("table doesn't exist")

type CorrectWildcards func(query string) string

const onSQLList = "on sqllib.SQLList(): "

func SQLList(table, fields string, options *crud.Options, correctWildcards CorrectWildcards) (string, []interface{}, error) {

	var join, order, limit string
	var values []interface{}

	var term *selectors.Term

	if options == nil {
		term = selectors.In("viewer_key", "")

	} else {
		viewerKey := options.ActorKey
		if options.Term == nil {
			term = selectors.In("viewer_key", "")
		} else {
			term = logic.AND(term, selectors.In("viewer_key", viewerKey))
		}

		if strings.TrimSpace(options.JoinTo.Clause) != "" {
			join = options.JoinTo.Clause
			values = options.JoinTo.Values
		} else if len(options.JoinTo.Values) > 0 {
			return "", nil, fmt.Errorf(onSQLList+"wrong .JoinTo: %#v", options.JoinTo)
		}

		if len(options.OrderBy) > 0 {
			order = " ORDER BY " + strings.Join(options.OrderBy, ", ")
		}

		if options.Offset+options.Limit > 0 {
			if options.Limit > 0 {
				limit += " LIMIT " + strconv.FormatUint(options.Limit, 10)
			}

			if options.Offset > 0 {
				limit += " OFFSET " + strconv.FormatUint(options.Offset, 10)
			}

			// TODO: sqlite & mysql version
		}
	}

	condition, valuesTerm, err := selectors_sql.Use(term)
	if err != nil {
		return "", nil, fmt.Errorf(onSQLList+"wrong selector (%#v): %s", term, err)
	}

	if strings.TrimSpace(condition) != "" {
		condition = " WHERE " + condition
	}

	query := "SELECT " + fields + " FROM " + table + join + condition + order + limit
	if correctWildcards != nil {
		query = correctWildcards(query)
	}

	return query, append(values, valuesTerm...), nil
}

const onSQLCount = "on sqllib.SQLCount(): "

func SQLCount(table string, options *crud.Options, correctWildcards CorrectWildcards) (string, []interface{}, error) {
	var term *selectors.Term
	if options == nil {
		term = selectors.In("viewer_key", "")

	} else if options.Term != nil {
		term = logic.AND(options.Term, selectors.In("viewer_key", options.ActorKey))

	} else {
		term = selectors.In("viewer_key", options.ActorKey)

	}

	condition, values, err := selectors_sql.Use(term)
	if err != nil {
		termStr, _ := json.Marshal(term)
		return "", nil, errata.Wrapf(err, onSQLCount+": can't selectors_sql.Use(%s)", termStr)
	}

	query := "SELECT COUNT(*) FROM " + table
	if strings.TrimSpace(condition) != "" {
		query += " WHERE " + condition
	}

	return query, values, nil
}

type SqlStmt struct {
	Stmt **sql.Stmt
	Sql  string
}

func Prepare(dbh *sql.DB, sqlQuery string, stmt **sql.Stmt) error {
	var err error

	*stmt, err = dbh.Prepare(sqlQuery)
	if err != nil {
		return errata.Wrapf(err, "can't dbh.Prepare(%s)", sqlQuery)
	}

	return nil
}

func Exec(dbh *sql.DB, sqlQuery string, values ...interface{}) (*sql.Result, error) {
	stmt, err := dbh.Prepare(sqlQuery)
	if err != nil {
		return nil, errata.Wrapf(err, CantPrepare, sqlQuery)
	}

	res, err := stmt.Exec(values...)
	if err != nil {
		return nil, errata.Wrapf(err, CantExec, sqlQuery, values)
	}

	return &res, nil
}

func Query(dbh *sql.DB, sqlQuery string, values ...interface{}) (*sql.Rows, error) {
	stmt, err := dbh.Prepare(sqlQuery)
	if err != nil {
		return nil, errata.Wrapf(err, CantPrepare, sqlQuery)
	}

	rows, err := stmt.Query(values...)
	if err != nil {
		return nil, errata.Wrapf(err, CantExec, sqlQuery, values)
	}

	return rows, nil
}

//func QueryStrings(stmt *sql.Stmt, sql string, values ...interface{}) (results []string, err error) {
//	rows, err := stmt.Query(values...)
//	if err != nil {
//		return nil, errors.Wrapf(err, CantExec, sql, values)
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
//	for rows.Right() {
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
