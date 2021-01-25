package persons_fs_stub

import (
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/selectors"
)

var _ crud.Cleaner = &personsFSStub{}

func (pfs *personsFSStub) Clean(*selectors.Term, *crud.Options) error {

	return errors.NotImplemented
}
