package sqllib_sqlite

import (
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/basis/common/sqllib"
	"github.com/pavlo67/workshop/basis/config"
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

func (sqlOp *SQLite) DB() *sql.DB {
	if sqlOp == nil {
		return nil
	}

	return sqlOp.db
}
