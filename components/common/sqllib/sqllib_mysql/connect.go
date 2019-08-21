package sqllib_mysql

import (
	"database/sql"
	"strconv"

	"github.com/pkg/errors"

	"github.com/pavlo67/constructor/components/common/config"
)

func Connect(cfg config.ServerAccess) (db *sql.DB, err error) {
	port := ":3306"
	if cfg.Port != 0 {
		port = ":" + strconv.Itoa(cfg.Port)
	}

	addr := cfg.User + ":" + cfg.Pass +
		"@tcp(" + cfg.Host + port + ")/" + cfg.Path + "?parseTime=true"

	dbh, err := sql.Open("mysql", addr)
	if err != nil {
		return nil, errors.Wrapf(err, "wrong db connect (credentials='%v')", cfg)
	}

	err = dbh.Ping()
	if err != nil {
		return nil, errors.Wrapf(err, "wrong .Ping on db connect (credentials='%v')", cfg)
	}
	return dbh, nil
}
