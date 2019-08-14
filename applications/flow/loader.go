package flow

import (
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/constructor/components/basis"
	"github.com/pavlo67/constructor/components/basis/logger"
	"github.com/pavlo67/constructor/components/processor/importer"
)

func Load(urls []string, impOp importer.Operator, adminOp Administrator, l logger.Operator) (numAll, numNewAll int, errs basis.Errors) {

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

		var num, numNew int

		for _, item := range items {
			num++
			ok, err := adminOp.Has(item.OriginKey)
			if err != nil {
				errs = append(errs, errors.Errorf("can't adminOp.Has(%#v): %s", item.OriginKey, err))
			} else if ok {
				// already exists!
				continue
			}

			item.SavedAt = &savedAt
			_, err = adminOp.Save([]importer.Item{item}, nil)
			if err != nil {
				errs = append(errs, errors.Errorf("can't adminOp.Save(%#v): %s", item, err))
			} else {
				numNew++
			}
		}

		numAll += num
		numNewAll += numNew

	}

	return numAll, numNewAll, errs
}
