package files_fs

import (
	"os"

	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/selectors"
)

var _ db.Cleaner = &filesFS{}

const onClean = "on filesFS.Clean()"

func (filesOp *filesFS) Clean(term *selectors.Term) error {
	if err := os.RemoveAll(filesOp.basePath); err != nil {
		return errors.Wrapf(err, onClean+": removing %s", filesOp.basePath)
	}

	return nil
}
