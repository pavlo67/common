package sqllib_sqlite

import (
	"database/sql"

	"github.com/pkg/errors"

	"strings"

	"github.com/pavlo67/constructor/components/basis"
	"github.com/pavlo67/constructor/components/basis/config"
	"github.com/pavlo67/constructor/components/basis/sqllib"
)

var _ sqllib.Operator = &SQLite{}

func New(cfg config.ServerAccess) (sqllib.Operator, error) {
	if strings.TrimSpace(cfg.Path) == "" {
		return nil, errors.New("no path to SQLite database is defined")
	}

	db, err := sql.Open("sqlite3", cfg.Path)
	if err != nil {
		return nil, errors.Wrapf(err, "wrong db connect (cfg = %#v)", cfg)
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.Wrapf(err, "wrong .Ping on db connect (cfg = %#v)", cfg)
	}

	return &SQLite{db}, nil
}

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
