package persons_fs_stub

import (
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/libraries/filelib"
	"github.com/pavlo67/common/common/selectors"
)

var _ crud.Cleaner = &personsFSStub{}

const onClean = "on personsFSStub.Clean()"

func (pfs *personsFSStub) Clean(*selectors.Term, *crud.Options) error {
	if err := filelib.ClearDir(pfs.path); err != nil {
		return errors.Wrap(err, onClean)
	}
	return nil
}
