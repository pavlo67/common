package importer_http_rss

import (
	"time"

	"github.com/mmcdole/gofeed"

	"encoding/json"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/flow"
	"github.com/pavlo67/workshop/common/types"
	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/tags"
)

var _ data.Convertor = &Item{}

type Item struct {
	sourceTime time.Time
	sourceURL  string
	feedItem   *gofeed.Item
}

func (item *Item) GetData() (*data.Item, error) {
	if item == nil || item.feedItem == nil {
		return nil, nil
	}

	feedItem := item.feedItem

	originalID := feedItem.GUID
	if originalID == "" {
		originalID = feedItem.Link
	}

	sourceTime := item.sourceTime

	// TODO: save in special field
	// if feedItem.PublishedParsed != nil {
	//	sourceTime = *feedItem.PublishedParsed
	// }

	status := crud.History{CreatedAt: sourceTime}

	var embedded []data.Item

	if feedItem.Image != nil {
		embedded = append(embedded, data.Item{
			TypeKey: types.KeyHRefImage,
			URL:     feedItem.Image.URL,
			Title:   feedItem.Image.Title,
			History: status,
		})
	}

	if len(feedItem.Enclosures) > 0 {
		for _, p := range feedItem.Enclosures {
			embedded = append(embedded, data.Item{
				TypeKey: types.KeyHRef,
				URL:     p.URL,
				Title:   p.Type + ": " + p.Length,
				History: status,
			})
		}
	}

	var items []tags.Item
	for _, c := range feedItem.Categories {
		items = append(items, tags.Item{Label: c})
	}

	origin, _ := json.Marshal(feedItem)

	return &data.Item{
		URL:      feedItem.Link,
		TypeKey:  types.KeyString,
		Title:    feedItem.Title,
		Summary:  feedItem.Description,
		Embedded: embedded,
		Tags:     items,
		Details:  feedItem.Content,
		History:  status,
		Origin: flow.Origin{
			Source: item.sourceURL,
			Key:    originalID,
			Time:   &sourceTime,
			Data:   string(origin),
		},
	}, nil

}

func (*Item) SaveData(data.Item) error {
	return common.ErrNotImplemented
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
