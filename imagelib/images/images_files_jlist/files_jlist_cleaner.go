package images_files_jlist

import (
	"os"

	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/filelib"
)

var _ db.Cleaner = &imagesFilesJList{}

const onClean = "on imagesFilesJList.Clean()"

func (op imagesFilesJList) Clean() error {
	if os.Getenv("ENV") != "test" {
		return errors.New("image files can't be cleaned in non-test mode")
	}

	d, err := os.Open(op.basePath)
	if err != nil {
		return nil
	}
	defer d.Close()

	if err = filelib.ClearDir(op.basePath); err != nil {
		return errors.Wrap(err, onClean)
	}

	return nil
}
