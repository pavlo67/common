package importer_actor

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/selectors"

	"encoding/json"

	"github.com/pavlo67/workshop/common/selectors/logic"
	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/importer/importer_rss"
)

func Load(dataOp data.Operator, l logger.Operator) error {
	url := "https://rss.unian.net/site/news_ukr.rss"

	l.Info(url)

	impOp := &importer_rss.RSS{}

	series, err := impOp.Get(url)
	if err != nil {
		return errors.Errorf("can't impOp.Get('%s', nil): %s", url, err)
	}
	if series == nil {
		return errors.Errorf("no series from impOp.Get('%s', nil)", url)
	}

	numAll := len(series.Items)
	var numProcessed, numNew int

	for _, item := range series.Items {
		var cnt uint64

		numProcessed++

		// term1 := selectors.Binary(selectors.Eq, "source", selectors.ValueUnary{url})

		// TODO!!!
		term := logic.AND(
			selectors.Binary(selectors.Eq, "source", selectors.Value{url}),
			selectors.Binary(selectors.Eq, "source_key", selectors.Value{item.Key}),
		)

		//itemStr, _ := json.Marshal(item)
		//l.Infof("%s ", itemStr)

		termStr, _ := json.Marshal(term)
		l.Infof("%s", termStr)

		cnt, err = dataOp.Count(term, nil)
		if err != nil {
			err = errors.Errorf("can't dataOp.Count(%#v): %s", term, err)
			break

		} else if cnt > 0 {
			// already exists!
			continue
		}

		_, err = dataOp.Save([]data.Item{item}, nil)
		if err != nil {
			err = errors.Errorf("can't adminOp.Save(%#v): %s", item, err)
			break

		} else {
			numNew++
		}
	}

	l.Infof("numAll = %d, numProcessed = %d, numNew = %d", numAll, numProcessed, numNew)

	return err
}
