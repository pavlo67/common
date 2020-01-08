package flowimporter_task

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/actor"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/selectors"

	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/dataimporter"
	"github.com/pavlo67/workshop/components/dataimporter/importer_http_rss"
	"github.com/pavlo67/workshop/components/datatagged"
	"github.com/pavlo67/workshop/components/sources"
)

func New(dataOp datatagged.Operator, sourcesOp sources.Operator) (actor.Operator, error) {

	if dataOp == nil {
		return nil, errors.New("on flowimporter_task.New(): data.Operator == nil")
	}
	if sourcesOp == nil {
		return nil, errors.New("on flowimporter_task.New(): sources.Operator == nil")
	}

	return &loadTask{dataOp, sourcesOp}, nil
}

var _ actor.Operator = &loadTask{}

type loadTask struct {
	dataOp    data.Operator
	sourcesOp sources.Operator
}

func (it *loadTask) Name() string {
	return "loader"
}

const onRun = "on loadTask.Run(): "

func (it *loadTask) Run(_ common.Map) (posterior []joiner.Link, info common.Map, err error) {
	if it == nil {
		return nil, nil, errors.New("on importer_task.Run(): loadTask == nil")
	}

	sourceItems, err := it.sourcesOp.List(nil, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, onRun+"can't it.sourcesOp.List(nil, nil)")
	}

	for i, sourceItem := range sourceItems {
		url := strings.TrimSpace(sourceItem.URL)
		if url == "" {
			l.Errorf("no url in sources.Item(%#v)", sourceItem)
			continue
		}

		l.Infof("%d f %d: %s", i+1, len(sourceItems), url)

		numAll, numProcessed, numNew, err := Load(url, it.dataOp)
		l.Infof("numAll = %d, numProcessed = %d, numNew = %d", numAll, numProcessed, numNew)

		if err != nil {
			l.Error(err)
		}
	}

	// TODO!!! return posterior
	return nil, nil, nil
}

func Load(url string, dataOp data.Operator) (int, int, int, error) {
	impOp, err := importer_http_rss.New(l)
	if err != nil {
		return 0, 0, 0, errors.Errorf("can't importer_rss.New(%#v): %s", l, err)
	}

	series, err := impOp.Get(url)
	if err != nil {
		return 0, 0, 0, errors.Errorf("can't impOp.Get('%s', nil): %s", url, err)
	}
	if series == nil {
		return 0, 0, 0, errors.Errorf("no series from impOp.Get('%s', nil)", url)
	}

	numAll := len(series.Data)
	var numProcessed, numNew int

	for _, item := range series.Data {
		var cnt uint64

		numProcessed++

		sourceKey := dataimporter.SourceKey(item.History)
		// TODO!!! check if both are not empty

		term := selectors.Binary(selectors.Eq, "source_key", selectors.Value{sourceKey})

		//itemStr, _ := json.Marshal(item)
		//l.Infof("%s ", itemStr)

		//termStr, _ := json.Marshal(term)
		//l.Infof("%s", termStr)

		cnt, err = dataOp.Count(term, nil)
		if err != nil {
			break

		} else if cnt > 0 {
			// already exists!
			continue
		}

		item.ID = ""

		_, err = dataOp.Save([]data.Item{item}, nil)
		if err != nil {
			break

		} else {
			numNew++
		}
	}

	return numAll, numProcessed, numNew, nil
}