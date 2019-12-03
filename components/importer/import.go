package importer

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/selectors"

	"github.com/pavlo67/workshop/components/data"
)

func Load(url string, impOp Operator, dataOp data.Operator, l logger.Operator) (numAll, numProcessed, numNew int, err error) {
	l.Info(url)

	series, err := impOp.Get(url)
	if err != nil {
		return numAll, numProcessed, numNew, errors.Wrapf(err, "can't impOp.Get('%s', nil)", url)
	}
	if series == nil {
		return numAll, numProcessed, numNew, errors.Errorf("no series from impOp.Get('%s', nil)", url)
	}

	numAll += len(series.Items)

	for _, item := range series.Items {
		numProcessed++

		// TODO!!!
		term := selectors.TermBinary(
			selectors.And,
			selectors.TermBinary(selectors.Eq, "source", url),
			selectors.TermBinary(selectors.Eq, "source_key", item.Key),
		)

		cnt, err := dataOp.Count(term, nil)
		if err != nil {
			return numAll, numProcessed, numNew, errors.Errorf("can't dataOp.Count(%#v): %s", term, err)
		} else if cnt > 0 {
			// already exists!
			continue
		}

		_, err = dataOp.Save([]data.Item{item}, nil)
		if err != nil {
			return numAll, numProcessed, numNew, errors.Errorf("can't adminOp.Save(%#v): %s", item, err)
		} else {
			numNew++
		}
	}

	return numAll, numProcessed, numNew, nil
}
