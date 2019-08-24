package data_importer

import (
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/basis/common"
	"github.com/pavlo67/workshop/basis/logger"
	"github.com/pavlo67/workshop/basis/instruments/importer"
)

func Load(urls []string, impOp importer.Operator, adminOp Administrator, l logger.Operator) (numAll, numProcessed, numNew int, errs common.Errors) {

	for _, url := range urls {
		l.Info(url)

		//err := impOp.Init(url)
		//if err != nil {
		//	errs = append(errs, errors.Errorf("can't impOp.Run('%s')", url, err))
		//	continue
		//}

		savedAt := time.Now()

		items, err := impOp.Get(url, nil)
		if err != nil {
			errs = append(errs, errors.Errorf("can't impOp.Get('%s', nil)", url, err))
			continue
		}

		numAll += len(items)

		for _, item := range items {
			numProcessed++
			ok, err := adminOp.Has(item.OriginKey)
			if err != nil {
				errs = append(errs, errors.Errorf("can't adminOp.Has(%#v): %s", item.OriginKey, err))
				break
			} else if ok {
				// already exists!
				continue
			}

			item.SavedAt = &savedAt
			_, err = adminOp.Save([]importer.Item{item}, nil)
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
