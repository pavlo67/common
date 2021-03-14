package persons_fs

import (
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/filelib"
	"github.com/pavlo67/common/common/selectors"
	"github.com/pkg/errors"
)

var _ db.Cleaner = &personsFSStub{}

const onClean = "on personsFSStub.Clean()"

func (personsOp *personsFSStub) Clean(*selectors.Term) error {
	if err := filelib.ClearDir(personsOp.path); err != nil {
		return errors.Wrap(err, onClean)
	}
	return nil
}
