package persons_fs

import (
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/filelib"
	"github.com/pkg/errors"
)

var _ crud.Cleaner = &personsFSStub{}

const onClean = "on personsFSStub.Clean()"

func (pfs *personsFSStub) Clean(*crud.Options) error {
	if err := filelib.ClearDir(pfs.path); err != nil {
		return errors.Wrap(err, onClean)
	}
	return nil
}
