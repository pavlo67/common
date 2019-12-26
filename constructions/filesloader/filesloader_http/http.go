package filesloader_http

import (
	"sort"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/common/libraries/httplib"
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

type toPrepare struct {
	url      string
	fileType string
	fileName string
	priority int
}

const onLoad = "on filesloaderHTTP.Load(): "

func (flOp *filesloaderHTTP) Load(urlToLoad, pathToStore string, priority filesloader.Priority) (*files.Item, error) {
	pathToStore, err := filelib.GetDir(pathToStore, flOp.pathToStoreDefault)
	if err != nil {
		return nil, errors.Wrapf(err, onLoad+"can't filelib.GetDir('%s', '%s')", pathToStore, flOp.pathToStoreDefault)
	}

	if priority == nil {
		priority = filesloader.PriorityDefault(urlToLoad, false)
	}

	var fileIndex int

	fileName, fileType, err := httplib.DownloadFile(urlToLoad, pathToStore, fileIndex, 0644)
	// TODO!!! postpone errors
	if err != nil {
		return nil, err
	}
	fileIndex++

	filesToPrepare := []toPrepare{{urlToLoad, fileType, fileName, 1}}

	for len(filesToPrepare) > 0 {
		fileToPrepare := filesToPrepare[0]
		filesToPrepare = filesToPrepare[1:]

		var posterior []toPrepare

		posterior, fileIndex, err = flOp.PreparePosterior(fileToPrepare, pathToStore, fileIndex, priority)
		// TODO!!! postpone errors
		if err != nil {
			return nil, err
		}

		if len(posterior) > 0 {
			filesToPrepare = append(filesToPrepare, posterior...)
			sort.Slice(filesToPrepare, func(i, j int) bool { return filesToPrepare[j].priority < filesToPrepare[i].priority })
		}
	}

	return &files.Item{urlToLoad}, nil
}

const onPreparePosterior = "on filesloaderHTTP.PreparePosterior(): "

func (flOp *filesloaderHTTP) PreparePosterior(fileToPrepare toPrepare, pathToStore string, fileIndex int, priority filesloader.Priority) ([]toPrepare, int, error) {
	return nil, fileIndex, nil
}

const onClean = "on filesloaderHTTP.Clean(): "

func (flOp *filesloaderHTTP) Clean(term *selectors.Term, _ *crud.RemoveOptions) error {
	return nil
}
