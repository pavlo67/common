package importer_rss

import (
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/components/importer"
)

var _ importer.Operator = &RSS{}

type RSS struct{}

//var reHTTP = regexp.MustCompile("(?i)^https?://")

const onGet = "on rss.Get(): "

func (r *RSS) Get(feedURL string) (*importer.Series, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		return nil, errors.Wrapf(err, onGet+"can't .ParseURL(%s)", feedURL)
	} else if feed == nil {
		return nil, errors.Errorf(onGet+"no feed obtained with .ParseURL(%s)", feedURL)
	}

	now := time.Now()

	series := importer.Series{URL: feedURL, CreatedAt: now}

	for _, feedItem := range feed.Items {
		item := &Item{
			sourceTime: now,
			sourceURL:  feedURL,
			feedItem:   feedItem,
		}

		dataItem, err := item.GetData()
		if err != nil {
			return &series, errors.Wrapf(err, onGet+"can't .GetData(%#v)", feedItem)
		}
		if dataItem == nil {
			// ??? wtf
			continue
		}

		series.Items = append(series.Items, *dataItem)
	}

	return &series, nil
}
