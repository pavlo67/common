package router_news

import (
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/processor/importer/importer_rss"
	"github.com/pavlo67/punctum/processor/news"
)

func Load(urls []string, newsOp news.Operator) (numAll, numNewAll int, errs basis.Errors) {
	impOp := &importer_rss.RSS{}

	for _, url := range urls {
		l.Info(url)

		err := impOp.Init(url)
		if err != nil {
			errs = append(errs, errors.Errorf("can't impOp.Run('%s')", url, err))
			continue
		}

		savedAt := time.Now()

		var num, numNew int

		for {
			item, err := impOp.Next()
			if err != nil {
				errs = append(errs, errors.Errorf("can't get next item: %v", err))
				continue
			}
			if item == nil {
				break
			}

			num++
			ok, err := newsOp.Has(&item.Source)
			if err != nil {
				errs = append(errs, errors.Errorf("can't newsOp.Has(%s): %s", item.Source, err))
			} else if ok {
				// already exists!
				continue
			}

			item.SavedAt = &savedAt
			err = newsOp.Save(item)
			if err != nil {
				errs = append(errs, errors.Errorf("can't newsOp.Save(%#v): %s", item, err))
			} else {
				numNew++
			}
		}

		numAll += num
		numNewAll += numNew

	}

	return numAll, numNewAll, errs
}
