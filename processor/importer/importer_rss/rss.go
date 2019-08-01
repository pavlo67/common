package importer_rss

import (
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/pavlo67/constructor/basis"
	"github.com/pavlo67/constructor/processor/flow"
	"github.com/pavlo67/constructor/processor/importer"
	"github.com/pavlo67/constructor/processor/news"
)

const SourceTypeRSS flow.SourceType = "rss"

var _ importer.Operator = &RSS{}

type RSS struct {
	feedURL   string
	language  string
	items     []*gofeed.Item
	itemIndex int
}

func (r *RSS) Init(feedURL string) error {
	r.feedURL = feedURL
	r.items = nil

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(r.feedURL)
	if err != nil {
		return err
	} else if feed == nil {
		return basis.ErrNull
	}

	r.itemIndex = -1
	r.language = feed.Language
	r.items = feed.Items

	return nil
}

//var reHTTP = regexp.MustCompile("(?i)^https?://")

func (r *RSS) Next() (*news.Item, error) {
	r.itemIndex++

	if r.itemIndex >= len(r.items) {
		return nil, nil
	}

	item := r.items[r.itemIndex]

	originalID := item.GUID
	if originalID == "" {
		originalID = item.Link
	}

	// contentKey := r.feedURL + "#" + originalID
	//proto := reHTTP.FindString(item.GUID)
	//if len(proto) == 0 {
	//	return item.GUID
	//}

	sourceTime := time.Now()
	if item.PublishedParsed != nil {
		sourceTime = *item.PublishedParsed
	}

	// Original:   interface{}(item),

	var embedded []news.Embedded

	if item.Image != nil {
		embedded = append(embedded, news.Embedded{
			SourceURL: item.Image.URL,
			Title:     item.Image.Title,
		})
	}

	if len(item.Enclosures) > 0 {
		for _, p := range item.Enclosures {
			embedded = append(embedded, news.Embedded{
				SourceURL: p.URL,
				Title:     p.Type + ": " + p.Length,
			})
		}
	}

	return &news.Item{
		Source: flow.Source{
			URL:      r.feedURL,
			SourceID: originalID,
		},
		Content: news.Content{
			Time:     &sourceTime,
			Title:    item.Title,
			Summary:  item.Description,
			Text:     item.Content,
			Embedded: embedded,
			Tags:     item.Categories,
			Href:     item.Link,
		},
	}, nil
}

func (r *RSS) Close() error {
	return nil
}

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

//// Object forms an items.Object from the imported entity
// func (entity Entity) Object() (obj *things.Object, err error) {
//	if entity.item == nil {
//		return nil, importer.ErrNilItem
//	}
//
//	item := entity.item
//
//	//language := ""
//	//if entity.rss != nil {
//	//	language = entity.rss.language
//	//}
//
//	createdLinks := []things.Link{}
//	if item.Author != nil {
//		//email, err := url.Parse(item.Author.Email)
//		//if err != nil {
//		//	email = nil
//		//}
//
//		createdLinks = append(createdLinks, things.Link{
//			Type: "author",
//			//Nick:    []items.Text{{Text: item.Author.Nick, Language: language}},
//			//Whereto: email,
//			Name: item.Author.Name,
//			To:   item.Author.Email,
//		})
//	}
//	if item.Link != "" {
//		//URL, err := url.Parse(item.Link)
//		//if err != nil {
//		//	URL = nil
//		//}
//
//		createdLinks = append(createdLinks, things.Link{
//			Type: "url",
//			//Nick:    []items.Text{{Text: item.Link, Language: language}},
//			//Whereto: URL,
//			Name: item.Link,
//			To:   item.Link,
//		})
//	}
//	if item.Image != nil {
//		//URL, err := url.Parse(item.Image.URL)
//		//if err != nil {
//		//	URL = nil
//		//}
//		createdLinks = append(createdLinks, things.Link{
//			Type: "image",
//			//Nick:    []items.Text{{Text: item.Label, Language: language}},
//			//Whereto: URL,
//			Name: item.Title,
//			To:   item.Image.URL,
//		})
//	}
//	for _, category := range item.Categories {
//		createdLinks = append(createdLinks, things.Link{
//			Type: links.TypeTag,
//			//Nick: []items.Text{{Text: category, Language: language}},
//			Name: category,
//		})
//	}
//
//	return &things.Object{
//		//Nick:    []items.Text{{Text: item.Label, Language: language}},
//		//Summary: []items.Text{{Text: item.Description, Language: language}},
//		Name:    item.Title,
//		Contentus: item.Description + " " + item.Contentus,
//		Links:   createdLinks,
//	}, nil
//
// }
