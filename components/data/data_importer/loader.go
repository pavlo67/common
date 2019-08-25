package data_importer

import (
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/basis/common"
	"github.com/pavlo67/workshop/basis/logger"
	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/instruments/importer"
)

func Load(urls []string, impOp importer.Operator, dataOp data.Operator, l logger.Operator) (numAll, numProcessed, numNew int, errs common.Errors) {

	for _, url := range urls {
		l.Info(url)

		//err := impOp.Init(url)
		//if err != nil {
		//	errs = append(errs, errors.Errorf("can't impOp.Run('%s')", url, err))
		//	continue
		//}

		savedAt := time.Now()

		series, err := impOp.Get(url, nil)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "can't impOp.Get('%s', nil)", url))
			continue
		}
		if series == nil {
			errs = append(errs, errors.Errorf("no series from impOp.Get('%s', nil)", url))
			continue
		}

		numAll += len(series.Items)

		for _, item := range series.Items {
			numProcessed++
			cnt, err := dataOp.Has(item.OriginKey, nil)
			if err != nil {
				errs = append(errs, errors.Errorf("can't adminOp.Has(%#v): %s", item.OriginKey, err))
				break
			} else if cnt > 0 {
				// already exists!
				continue
			}

			item.SavedAt = &savedAt
			_, err = dataOp.Save([]data.Item{item}, nil, nil, nil)
			if err != nil {
				errs = append(errs, errors.Errorf("can't adminOp.Save(%#v): %s", item, err))
				break
			} else {
				numNew++
			}
		}

	}

	return numAll, numProcessed, numNew, errs
}
