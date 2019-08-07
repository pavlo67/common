package sqllib

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/pavlo67/constructor/starter/config"
)

type Operator interface {
	DB() (*sql.DB, error)

	TableExistsSQL(tableName string) (string, error)
	CreateSQL(table config.SQLTable) (string, error)
}

// helpers ------------------------------------------------------------

func Close(op Operator) error {
	db, err := op.DB()
	if err != nil {
		return err
	}

	return db.Close()
}

const onTableExists = "on TableExists()"

func TableExists(op Operator, tableName string) (bool, error) {
	sqlQuery, err := op.TableExistsSQL(tableName)
	if err != nil {
		return false, errors.Wrap(err, onTableExists)
	}

	db, err := op.DB()
	if err != nil {
		return false, errors.Wrap(err, onTableExists)
	}

	rows, err := Query(db, sqlQuery)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return false, errors.Wrap(err, onTableExists)
	}

	var t string
	for rows.Next() {
		err = rows.Scan(&t)
		if err != nil {
			return false, errors.Wrapf(err, CantScanQueryRow, sqlQuery, nil)
		}
		return true, nil
	}
	err = rows.Err()
	if err != nil {
		return false, errors.Wrapf(err, CantScanQueryRow, sqlQuery, nil)
	}

	return false, nil
}

func DropTable(dbh *sql.DB, table string) error {
	sqlQuery := "DROP TABLE IF EXISTS`" + table + "`"

	_, err := Exec(dbh, sqlQuery)

	if err != nil {
		return err
	}
	return nil
}
