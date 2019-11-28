package importer_rss

import (
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/flow"
	"github.com/pavlo67/workshop/components/importer"
	"github.com/pavlo67/workshop/components/tagger"
)

var _ importer.Operator = &RSS{}

type RSS struct{}

//var reHTTP = regexp.MustCompile("(?i)^https?://")

func (r *RSS) Get(feedURL string) (*importer.Series, error) {

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		return nil, errors.Wrapf(err, "can't .ParseURL(%s)", feedURL)
	} else if feed == nil {
		return nil, errors.Errorf("no feed obtained with .ParseURL(%s)", feedURL)
	}

	now := time.Now()

	series := importer.Series{URL: feedURL}

	for _, item := range feed.Items {
		originalID := item.GUID
		if originalID == "" {
			originalID = item.Link
		}

		sourceTime := now
		if item.PublishedParsed != nil {
			sourceTime = *item.PublishedParsed
		}

		var embedded []data.Item
		if strings.TrimSpace(item.Link) != "" {
			embedded = append(embedded, data.Item{
				URL: importer.ImportedHREF + item.Link,
			})
		}

		if item.Image != nil {
			embedded = append(embedded, data.Item{
				URL:   importer.ImportedHREF + item.Image.URL,
				Title: item.Image.Title,
			})
		}

		if len(item.Enclosures) > 0 {
			for _, p := range item.Enclosures {
				embedded = append(embedded, data.Item{
					URL:   importer.ImportedHREF + p.URL,
					Title: p.Type + ": " + p.Length,
				})
			}
		}

		var tags []tagger.Tag
		for _, c := range item.Categories {
			tags = append(tags, tagger.Tag(c))
		}

		series.Items = append(series.Items, data.Item{
			Title:    item.Title,
			Summary:  item.Description,
			Embedded: embedded,
			Tags:     tags,
			Details:  item.Content,
			Status:   crud.Status{CreatedAt: now},
			Origin: flow.Origin{
				Source: feedURL,
				Key:    originalID,
				Time:   &sourceTime,
				Data:   &item,
			},
		})
	}

	return &series, nil
}

// language := feed.Language

// type Census struct {
// 	Label        string                `json:"title,omitempty"`
// 	Description    string                `json:"description,omitempty"`
// 	Contentus        string                `json:"content,omitempty"`
// 	Link        string                `json:"link,omitempty"`
// 	Updated        string                `json:"updated,omitempty"`
// 	UpdatedParsed    *time.Time            `json:"updatedParsed,omitempty"`
// 	Published    string                `json:"published,omitempty"`
// 	PublishedParsed    *time.Time            `json:"publishedParsed,omitempty"`
// 	Author        *Person                `json:"author,omitempty"`
// 	GUID        string                `json:"guid,omitempty"`
// 	Image        *Image                `json:"image,omitempty"`
// 	Categories    []string            `json:"categories,omitempty"`
// 	Enclosures    []*Enclosure            `json:"enclosures,omitempty"`
// 	DublinCoreExt    *ext.DublinCoreExtension    `json:"dcExt,omitempty"`
// 	ITunesExt    *ext.ITunesItemExtension    `json:"itunesExt,omitempty"`
// 	Extensions    ext.Extensions            `json:"extensions,omitempty"`
// 	Custom        map[string]string        `json:"custom,omitempty"`
// }
