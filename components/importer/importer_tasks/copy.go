package importer_tasks

import (
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/scheduler"
	"github.com/pavlo67/workshop/common/selectors"
	"github.com/pavlo67/workshop/common/selectors/logic"

	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/importer"
	"github.com/pavlo67/workshop/components/importer/importer_series_http"
)

// TODO: vary this parameter
const copyLimit = 500

func NewCopyTask(url string, dataOp data.Operator) (scheduler.Task, error) {
	url = strings.TrimSpace(url)

	if url == "" {
		return nil, errors.New("on importer_tasks.NewCopyTask(): empty url")
	}

	if dataOp == nil {
		return nil, errors.New("on importer_task.NewCopyTask(): data.Operator == nil")
	}

	impOp, err := importer_series_http.NewSeriesHTTP(url, l)
	if err != nil {
		return nil, errors.Errorf("on importer_tasks.NewCopyTask(): can't importer_series_http.NewSeriesHTTP(%s, '', l)", url)
	}

	// TODO!!! init .lastImportedId

	return &copyTask{url, impOp, "", dataOp, copyLimit}, nil
}

var _ scheduler.Task = &copyTask{}

type copyTask struct {
	url            string
	impOp          importer.Operator
	lastImportedID string

	dataOp    data.Operator
	copyLimit int
}

func (it *copyTask) Name() string {
	return "copier from series_http"
}

func (it *copyTask) Run(timeSheduled time.Time) error {
	if it == nil {
		return errors.New("on copyTask.Run(): it == nil")
	}

	numAll, numProcessed, numNew, err := it.Copy()
	l.Infof("numAll = %d, numProcessed = %d, numNew = %d", numAll, numProcessed, numNew)

	return err
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

		term := logic.AND(
			selectors.Binary(selectors.Eq, "source", selectors.Value{item.Origin.Source}),
			selectors.Binary(selectors.Eq, "source_key", selectors.Value{item.Origin.Key}),
		)

		//itemStr, _ := json.Marshal(item)
		//l.Infof("%s ", itemStr)

		//termStr, _ := json.Marshal(term)
		//l.Infof("%s", termStr)

		cnt, err = it.dataOp.Count(term, nil)
		if err != nil {
			err = errors.Errorf("can't dataOp.Count(%#v): %s", term, err)
			return numAll, numProcessed, numNew, err

		} else if cnt > 0 {
			// already exists!
			continue
		}

		importedID := item.ID
		item.ID = ""

		// TODO: remove this kostyl!!!
		item.Origin.Time = &item.Status.CreatedAt

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
