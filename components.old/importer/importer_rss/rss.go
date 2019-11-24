package importer_rss

import (
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/instruments/importer"
)

// const SourceTypeRSS f.SourceType = "rss"

var _ importer.Operator = &RSS{}

type RSS struct {
	sourceID common.ID
}

//var reHTTP = regexp.MustCompile("(?i)^https?://")

func (r *RSS) Get(feedURL string, minKey *string) (*importer.Series, error) {

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		return nil, errors.Wrapf(err, "can't .ParseURL(%s)", feedURL)
	} else if feed == nil {
		return nil, errors.Errorf("no feed obtained with .ParseURL(%s)", feedURL)
	}

	// language := feed.Language

	series := importer.Series{
		URL: feedURL,
	}

	for _, item := range feed.Items {
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

		var embedded []data.Content
		if item.Image != nil {
			embedded = append(embedded, data.Content{
				Href:  item.Image.URL,
				Title: item.Image.Title,
			})
		}
		if len(item.Enclosures) > 0 {
			for _, p := range item.Enclosures {
				embedded = append(embedded, data.Content{
					Href:  p.URL,
					Title: p.Type + ": " + p.Length,
				})
			}
		}

		series.Items = append(series.Items, data.Item{
			SourceURL:  feedURL,
			SourceTime: &sourceTime,

			Origin: data.Origin{
				ID:  r.sourceID,
				Key: originalID,
			},

			// OriginData

			Content: data.Content{
				Title:   item.Title,
				Summary: item.Description,
				Details: item.Content,
				Href:    item.Link,
			},

			Embedded: embedded,
			Tags:     item.Categories,
		})

	}

	return &series, nil
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
//			Title: item.Author.Title,
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
//			Title: item.Link,
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
//			Title: item.Title,
//			To:   item.Image.URL,
//		})
//	}
//	for _, category := range item.Categories {
//		createdLinks = append(createdLinks, things.Link{
//			Type: links.TypeTag,
//			//Nick: []items.Text{{Text: category, Language: language}},
//			Title: category,
//		})
//	}
//
//	return &things.Object{
//		//Nick:    []items.Text{{Text: item.Label, Language: language}},
//		//Summary: []items.Text{{Text: item.Description, Language: language}},
//		Title:    item.Title,
//		Contentus: item.Description + " " + item.Contentus,
//		Tags:   createdLinks,
//	}, nil
//
// }
