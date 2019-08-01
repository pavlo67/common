package news

import (
	"time"

	"github.com/pavlo67/constructor/processor/flow"
	"github.com/pavlo67/constructor/starter/joiner"
)

const InterfaceKey joiner.InterfaceKey = "news"

type Item struct {
	ID          string `bson:"_id,omitempty" json:"id,omitempty"`
	flow.Source `bson:",inline"       json:",inline"`
	Content     `bson:"content"       json:"content"`
	SavedAt     *time.Time `bson:"saved_at"      json:"saved_at"`
}

type Content struct {
	Title    string     `bson:"title,omitempty"    json:"title,omitempty"`
	Summary  string     `bson:"summary,omitempty"  json:"summary,omitempty"`
	Text     string     `bson:"text,omitempty"     json:"text,omitempty"`
	Tags     []string   `bson:"tags,omitempty"     json:"tags,omitempty"`
	Embedded []Embedded `bson:"embedded,omitempty" json:"embedded,omitempty"`
	Href     string     `bson:"href,omitempty"     json:"href,omitempty"`
	Time     *time.Time `bson:"time,omitempty"     json:"time,omitempty"`
}

type Embedded struct {
	SourceURL string `bson:"source_url,omitempty" json:"source_url,omitempty"`
	Href      string `bson:"href,omitempty"       json:"href,omitempty"`
	Title     string `bson:"title,omitempty"      json:"title,omitempty"`
}

type Operator interface {
	Has(*flow.Source) (bool, error)
	Save(item *Item) error
	ReadList(*content.ListOptions) ([]Item, *uint64, error)
	DeleteList(*content.ListOptions) error
	Close() error
}

//func (item *Item) PartesTexti() ([]textus.Pars, error) {
//	if item == nil || item.Content == nil {
//		return nil, basis.ErrNull
//	}
//
//	return []textus.Pars{
//		{
//			Fons:            item.Source,
//			Origo:           item.Original,
//			ClavisContentus: item.ContentKey,
//			Contentus: &textus.Contentus{
//				Titulus:    item.Content.Title,
//				Index:      item.Content.Summary,
//				Textus:     item.Content.Text,
//				Appendices: map[string][]string{"tags": item.Content.Tags},
//			},
//		},
//	}, nil
//
//}
