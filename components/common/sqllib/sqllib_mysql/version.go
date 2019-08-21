package sqllib_mysql

import (
	"regexp"
)

var reVersion = regexp.MustCompile("^\\d+\\.\\d+")

//func MySQLVersion(dbh *sql.DB) (string, error) {
//	var stmt *sql.Stmt
//	sqlQuery := "select VERSION()"
//	if err := Exec(dbh, sqlQuery, &stmt); err != nil {
//		return "", err
//	}
//	rows, err := stmt.Query()
//	if err != nil {
//		return "", errors.Wrapf(err, "can't query SQL: %s", sqlQuery)
//	}
//	if rows != nil {
//		defer rows.Close()
//	}
//
//	var v []byte
//	for rows.Next() {
//		err = rows.Scan(&v)
//		if err != nil {
//			return "", errors.Wrapf(err, "can't scan SQL: %s", sqlQuery)
//		}
//
//		version := reVersion.Find(v)
//		if version != nil {
//			parts := strings.Split(string(version), ".")
//			if len(parts) != 2 {
//				return "", errors.Errorf("can't get version from '%s'", version)
//			}
//			major := parts[0]
//			if len(major) < 2 {
//				major = "0" + major
//			}
//
//			minor := parts[1]
//			if len(minor) < 2 {
//				minor = "0" + minor
//			}
//
//			return major + "." + minor, nil
//		}
//	}
//
//	return "", errors.New("there is no mysql version???")
//
//}
//
