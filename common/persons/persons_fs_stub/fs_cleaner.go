package persons_fs_stub

import (
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/errata"
	"github.com/pavlo67/common/common/libraries/filelib"
)

var _ crud.Cleaner = &personsFSStub{}

const onClean = "on personsFSStub.Clean()"

func (pfs *personsFSStub) Clean(*crud.Options) error {
	if err := filelib.ClearDir(pfs.path); err != nil {
		return errata.Wrap(err, onClean)
	}
	return nil
}
