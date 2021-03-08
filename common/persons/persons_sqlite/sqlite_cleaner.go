package persons_sqlite

import (
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/sqllib"
)

var _ crud.Cleaner = &personsSQLite{}

const onClean = "on personsSQLite.Clean(): "

func (personsOp *personsSQLite) Clean(_ *crud.Options) error {
	if _, err := personsOp.stmClean.Exec(); err != nil {
		return errors.Wrapf(err, onClean+sqllib.CantExec, personsOp.sqlClean, nil)
	}

	return nil
}
