package news_router

import (
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/processor/importer/importer_rss"
	"github.com/pavlo67/punctum/processor/news"
)

func Load(urls []string, newsOp news.Operator) (int, int, basis.Errors) {
	err := impOp.Init(sourceURL)
	if err != nil {
		return 0, 0, basis.Errors{errors.Errorf("can't impOp.Init('%s')", sourceURL, err)}
	}

	rssOp := &importer_rss.RSS{}

	savedAt := time.Now()

	for _, rssURL := range urls {
		l.Info(rssURL)

		num, numNew, errs := Load(rssOp, urls, newsOp)

		numAll += num
		numNewAll += numNew
		errsAll = append(errsAll, errs...)
	}

	var num, numNew int
	var errs basis.Errors

	for {
		item, err := impOp.Next()
		if err != nil {
			l.Infof("can't get next item: %v", err)
			continue
		}
		if item == nil {
			break
		}

		num++
		ok, err := newsOp.Has(&item.Source)
		if err != nil {
			errs = append(errs, err)
		} else if ok {
			// already exists!
			continue
		}

		item.SavedAt = &savedAt
		err = newsOp.Save(item)
		if err != nil {
			errs = append(errs, err)
		} else {
			numNew++
		}
	}

	return num, numNew, errs

}
