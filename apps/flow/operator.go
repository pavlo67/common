package flow

import (
	"time"

	"github.com/pavlo67/constructor/apps/content"
	"github.com/pavlo67/constructor/apps/links"
	"github.com/pavlo67/constructor/auth"
	"github.com/pavlo67/constructor/basis"
	"github.com/pavlo67/constructor/basis/selectors"
	"github.com/pavlo67/constructor/starter/joiner"
	"github.com/pavlo67/constructor/structura"
)

const InterfaceKey joiner.InterfaceKey = "flow"

type Item struct {
	ID      basis.ID   `bson:"_id,omitempty"      json:"id,omitempty"`
	SavedAt *time.Time `bson:"saved_at,omitempty" json:"saved_at,omitempty"`

	Source     `           bson:",inline"               json:",inline"`
	SourceTime *time.Time `bson:"source_time,omitempty" json:"source_time,omitempty"`

	Content  `          bson:"content"            json:"content"`
	Embedded []Content `bson:"embedded,omitempty" json:"embedded,omitempty"`

	Tags  []links.Tag `bson:"tags,omitempty"     json:"tags,omitempty"`
	RView auth.ID     `bson:"r_view,omitempty"  json:"r_view,omitempty"`
}

type Source struct {
	SourceID basis.ID `bson:"source_id,omitempty"       json:"source_id,omitempty"`
	Original string   `bson:"original,omitempty"  json:"original,omitempty"`
}

type Content struct {
	SourceURL string      `bson:"source_url,omitempty" json:"source_url,omitempty"`
	Title     string      `bson:"title"                json:"title"`
	Summary   string      `bson:"summary,omitempty"    json:"summary,omitempty"`
	Text      string      `bson:"text,omitempty"       json:"text,omitempty"`
	Tags      []links.Tag `bson:"tags,omitempty"       json:"tags,omitempty"`
	Href      string      `bson:"href,omitempty"       json:"href,omitempty"`
}

type Operator interface {
	Read(basis.ID, *structura.GetOptions) (Item, error)
	List(*selectors.Term, *structura.GetOptions) ([]content.Brief, error)
	Tags(*selectors.Term, *structura.GetOptions) ([]links.Tag, error)
	Sources(*selectors.Term, *structura.GetOptions) ([]basis.ID, error)
	Close() error
}

type Administrator interface {
	Has(*Source) (bool, error)
	Save([]Item, *structura.SaveOptions) ([]basis.ID, error)
	Remove(*selectors.Term, *structura.RemoveOptions) error
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

//func (src *Source) Key(keyAdd string) string {
//	if src == nil {
//		return ""
//	}
//
//	url := strings.TrimSpace(src.URL)
//
//	pos := strings.Index(url, "#")
//	if pos >= 0 {
//		url = url[:pos]
//	}
//
//	if url == "" {
//		return ""
//	}
//
//	if len(keyAdd) > 0 {
//		url += "#" + keyAdd
//	}
//
//	sourceID := strings.TrimSpace(src.SourceID)
//	if sourceID == "" {
//		return url
//	}
//
//	return url + "#" + sourceID
//}
