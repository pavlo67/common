package sqllib_sqlite

import (
	"database/sql"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/constructor/basis"
	"github.com/pavlo67/constructor/basis/sqllib"
	"github.com/pavlo67/constructor/starter/config"
)

var _ sqllib.Operator = &SQLite{}

type SQLite struct {
	db *sql.DB
}

func (sqlOp *SQLite) DB() (*sql.DB, error) {
	if sqlOp == nil {
		return nil, basis.ErrNull
	}

	if sqlOp.db == nil {
		return nil, errors.New("no SQLite connection")
	}

	return sqlOp.db, nil
}

func (sqlOp *SQLite) Connect(cfg config.ServerAccess) error {
	if sqlOp == nil {
		return basis.ErrNull
	}

	if strings.TrimSpace(cfg.Path) == "" {
		return errors.New("no path to SQLite database is defined")
	}

	var err error
	sqlOp.db, err = sql.Open("sqlite3", cfg.Path)
	if err != nil {
		return errors.Wrapf(err, "wrong db connect (cfg = %#v)", cfg)
	}

	err = sqlOp.db.Ping()
	if err != nil {
		return errors.Wrapf(err, "wrong .Ping on db connect (cfg = %#v)", cfg)
	}
	return nil
}

var reIndexLimit = regexp.MustCompile(`\(\d+\)`)

func (_ SQLite) CreateSQLQuery(table config.SQLTable) (string, error) {

	//reText := regexp.MustCompile("text")
	//reTime := regexp.MustCompile("(datetime|timestamp)")

	sqlQuery := "create table `" + table.Name + "` ( "
	firsFieldAdded := false

	for _, f := range table.Fields {
		if firsFieldAdded {
			sqlQuery += ", \n"
		}
		sqlQuery += "`" + f.Name + "` " + f.Type + " "
		if !f.Null {
			sqlQuery += " NOT NULL "
		} else {
			sqlQuery += " NULL "
		}

		if f.Default != "" {
			sqlQuery += " DEFAULT '" + f.Default + "' "
		}

		if f.Extra != "" {
			sqlQuery += f.Extra + " "
		}
		firsFieldAdded = true
	}
	for _, i := range table.Indexes {
		sqlQuery += ", \n"
		if strings.ToUpper(i.Type) == "PRIMARY" {
			sqlQuery += "PRIMARY KEY  ("
			liF := 0
			for _, f := range i.Fields {
				if liF > 0 {
					sqlQuery += ", "
				}
				sqlQuery += "`" + f + "`"
				liF++
			}
			sqlQuery += ")"
		}
		if strings.ToUpper(i.Type) == "UNIQUE" || i.Type == "" {
			if i.Type == strings.ToUpper("UNIQUE") {
				sqlQuery += "UNIQUE KEY `" + i.Name + "` ("
			} else {
				sqlQuery += "KEY `" + i.Name + "` USING BTREE ("
			}
			liF := 0
			for _, f := range i.Fields {
				if liF > 0 {
					sqlQuery += ", "
				}
				if reIndexLimit.MatchString(f) {
					sqlQuery += f
				} else {
					sqlQuery += "`" + f + "`"
				}
				liF++
			}
			sqlQuery += ")"
		}
	}
	sqlQuery += ");"

	return sqlQuery, nil
}
