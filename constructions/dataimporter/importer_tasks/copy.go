package importer_tasks

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/actor"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/selectors"
	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/constructions/dataimporter"
	"github.com/pavlo67/workshop/constructions/dataimporter/importer_http_series"
)

// TODO: vary this parameter
const copyLimit = 500

func NewCopyTask(url string, dataOp data.Operator) (actor.Operator, error) {
	url = strings.TrimSpace(url)

	if url == "" {
		return nil, errors.New("on importer_tasks.NewCopyTask(): empty url")
	}

	if dataOp == nil {
		return nil, errors.New("on importer_task.NewCopyTask(): data.Operator == nil")
	}

	impOp, err := importer_http_series.NewSeriesHTTP(url, l)
	if err != nil {
		return nil, errors.Errorf("on importer_tasks.NewCopyTask(): can't importer_series_http.NewSeriesHTTP(%s, '', l)", url)
	}

	// TODO!!! init .lastImportedId

	return &copyTask{url, impOp, "", dataOp, copyLimit}, nil
}

var _ actor.Operator = &copyTask{}

type copyTask struct {
	url            string
	impOp          dataimporter.Operator
	lastImportedID string

	dataOp    data.Operator
	copyLimit int
}

func (it *copyTask) Name() string {
	return "copier from series_http"
}

func (it *copyTask) Run(_ common.Map) (posterior []joiner.Link, info common.Map, err error) {
	if it == nil {
		return nil, nil, errors.New("on copyTask.Run(): it == nil")
	}

	numAll, numProcessed, numNew, err := it.Copy()
	l.Infof("numAll = %d, numProcessed = %d, numNew = %d", numAll, numProcessed, numNew)

	return nil, nil, err
}

func (it *copyTask) Copy() (int, int, int, error) {
	if it == nil {
		return 0, 0, 0, errors.New("on copyTask.Copy(): it == nil")
	}

	// l.Info(it.lastImportedID)

	series, err := it.impOp.Get(it.lastImportedID)
	if err != nil {
		return 0, 0, 0, errors.Errorf("can't impOp.Get(%s) from %s: %s", it.lastImportedID, it.url, err)
	}
	if series == nil {
		return 0, 0, 0, errors.Errorf("no series from impOp.Get(%s) from %s", it.lastImportedID, it.url)
	}

	numAll := len(series.Data)
	var numProcessed, numNew int

	for _, item := range series.Data {
		var cnt uint64

		// l.Info("? ", item.ID)

		numProcessed++

		sourceKey := dataimporter.SourceKey(item.History)
		// TODO!!! check if both are not empty

		term := selectors.Binary(selectors.Eq, "source_key", selectors.Value{sourceKey})

		//itemStr, _ := json.Marshal(item)
		//l.Infof("%s ", itemStr)

		//termStr, _ := json.Marshal(term)
		//l.Infof("%s", termStr)

		cnt, err = it.dataOp.Count(term, nil)
		if err != nil {
			err = errors.Errorf("can't dataOp.CountTags(%#v): %s", term, err)
			return numAll, numProcessed, numNew, err

		} else if cnt > 0 {
			// already exists!
			continue
		}

		importedID := item.ID
		item.ID = ""

		_, err = it.dataOp.Save([]data.Item{item}, nil)
		if err != nil {
			err = errors.Errorf("can't adminOp.Save(%#v): %s", item, err)
			break

		}

		numNew++
		it.lastImportedID = string(importedID)

		// l.Info("--> ", it.lastImportedID)

		if numNew >= copyLimit {
			break
		}
	}

	return numAll, numProcessed, numNew, nil
}
