package filesloader_http

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/common/selectors"

	"github.com/pavlo67/workshop/components/files"

	"github.com/pavlo67/workshop/constructions/filesloader"
)

var _ filesloader.Operator = &filesloaderHTTP{}

type filesloaderHTTP struct {
	pathToStoreDefault string
}

const onNew = "on filesloaderHTTP.New(): "

func New(pathToStoreDefault string) (filesloader.Operator, crud.Cleaner, error) {
	pathToStoreDefault, err := filelib.GetDir(pathToStoreDefault, "./")
	if err != nil {
		return nil, nil, errors.Wrapf(err, onNew+"can't filelib.GetDir('%s', './')", pathToStoreDefault)
	}

	flOp := filesloaderHTTP{
		pathToStoreDefault: pathToStoreDefault,
	}

	return &flOp, nil, nil
}

const onLoad = "on filesloaderHTTP.Load(): "

func (flOp *filesloaderHTTP) Load(pathToLoad, pathToStore string) (*files.Item, error) {
	pathToStore, err := filelib.GetDir(pathToStore, flOp.pathToStoreDefault)
	if err != nil {
		return nil, errors.Wrapf(err, onLoad+"can't filelib.GetDir('%s', '%s')", pathToStore, flOp.pathToStoreDefault)
	}

	return nil, common.ErrNotImplemented
}

const onClean = "on filesloaderHTTP.Clean(): "

func (flOp *filesloaderHTTP) Clean(term *selectors.Term, _ *crud.RemoveOptions) error {
	return nil
}
