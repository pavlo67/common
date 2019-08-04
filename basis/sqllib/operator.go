package sqllib

import (
	"database/sql"

	"github.com/pavlo67/constructor/starter/config"
)

type Operator interface {
	DB() (*sql.DB, error)
	Connect(cfg config.ServerAccess) error
	CreateSQLQuery(table config.SQLTable) (string, error)
}
