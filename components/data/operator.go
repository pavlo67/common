package data

import (
	"time"

	"github.com/pavlo67/workshop/basis/common"
	"github.com/pavlo67/workshop/basis/crud"
	"github.com/pavlo67/workshop/basis/joiner"
	"github.com/pavlo67/workshop/basis/selectors"
	"github.com/pavlo67/workshop/components/instruments/indexer"
	"github.com/pavlo67/workshop/components/marks"
)

const InterfaceKey joiner.InterfaceKey = "data"

type Item struct {
	ID      common.ID  `bson:"_id,omitempty"      json:"id,omitempty"`
	SavedAt *time.Time `bson:"saved_at,omitempty" json:"saved_at,omitempty"`

	OriginKey `       bson:",inline"          json:",inline"`
	Origin    string `bson:"origin,omitempty" json:"origin,omitempty"`

	SourceURL  string     `bson:"source_url,omitempty"  json:"source_url,omitempty"`
	SourceTime *time.Time `bson:"source_time,omitempty" json:"source_time,omitempty"`

	Content  `          bson:"content"            json:"content"`
	Embedded []Content `bson:"embedded,omitempty" json:"embedded,omitempty"`

	Tags  []string       `bson:"tags,omitempty"  json:"tags,omitempty"`
	Index map[string]int `bson:"index,omitempty" json:"index,omitempty"`
}

type OriginKey struct {
	SourceID  common.ID `bson:"source_id,omitempty"  json:"source_id,omitempty"`
	SourceKey string    `bson:"source_key,omitempty" json:"source_key,omitempty"`
}

type Content struct {
	Type    crud.Type `bson:"type"              json:"type"`
	Title   string    `bson:"title"             json:"title"`
	Summary string    `bson:"summary,omitempty" json:"summary,omitempty"`
	Details string    `bson:"details,omitempty" json:"details,omitempty"`
	Href    string    `bson:"href,omitempty"    json:"href,omitempty"`
}

type Operator interface {
	Has(OriginKey, *crud.GetOptions) (uint, error)
	Read(common.ID, *crud.GetOptions) (*Item, error)

	Save([]Item, marks.Operator, indexer.Operator, *crud.SaveOptions) ([]common.ID, error)
	Remove(selectors.Term, marks.Operator, indexer.Operator, *crud.RemoveOptions) error

	List(selectors.Term, indexer.Operator, *crud.GetOptions) ([]crud.Brief, error)
	Count(selectors.Term, indexer.Operator, *crud.GetOptions) ([]crud.Part, error)
	Reindex(selectors.Term, indexer.Operator, *crud.GetOptions) error
}

//func (item *Item) PartesTexti() ([]textus.Pars, error) {
//	if item == nil || item.Content == nil {
//		return nil, basis.ErrNull
//	}
//
//	return []textus.Pars{
//		{
//			Fons:            item.Origin,
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

//func (src *Origin) Key(keyAdd string) string {
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
