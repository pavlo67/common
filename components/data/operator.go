package data

import (
	"time"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/selectors"
	"github.com/pavlo67/workshop/components/instruments/indexer"
	"github.com/pavlo67/workshop/components/marks"
)

const InterfaceKey joiner.InterfaceKey = "data"

type Item struct {
	Brief   `            bson:",inline"           json:",inline"`
	Details interface{} `bson:"details,omitempty" json:"details,omitempty"`

	Tags  []string       `bson:"tags,omitempty"  json:"tags,omitempty"`
	Index map[string]int `bson:"index,omitempty" json:"index,omitempty"`

	Origin     `           bson:"origin,omitempty"      json:"origin,omitempty"`
	OriginTime *time.Time `bson:"origin_time,omitempty" json:"origin_time,omitempty"`
	OriginData string     `bson:"origin_data,omitempty" json:"origin_data,omitempty"`
}

type Origin struct {
	ID  common.ID `bson:"id,omitempty"  json:"id,omitempty"`
	Key string    `bson:"key,omitempty" json:"key,omitempty"`
}

type Brief struct {
	crud.Brief `             bson:",inline"            json:",inline"`
	Embedded   []crud.Brief `bson:"embedded,omitempty" json:"embedded,omitempty"`
	SavedAt    time.Time    `bson:"saved_at,omitempty" json:"saved_at,omitempty"`
}

type Operator interface {
	Has(Origin, *crud.GetOptions) (uint, error)
	Read(common.ID, *crud.GetOptions) (*Item, error)

	Save([]Item, marks.Operator, indexer.Operator, *crud.SaveOptions) ([]common.ID, error)
	Remove(*selectors.Term, marks.Operator, indexer.Operator, *crud.RemoveOptions) error

	List(*selectors.Term, indexer.Operator, *crud.GetOptions) ([]Brief, error)
	Count(*selectors.Term, indexer.Operator, *crud.GetOptions) ([]crud.Part, error)
	Reindex(*selectors.Term, indexer.Operator, *crud.GetOptions) error
}

//func (item *Item) PartesTexti() ([]textus.Pars, error) {
//	if item == nil || item.Content == nil {
//		return nil, basis.ErrNull
//	}
//
//	return []textus.Pars{
//		{
//			Fons:            item.OriginData,
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

//func (src *OriginData) Key(keyAdd string) string {
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
//	sourceID := strings.TrimSpace(src.ID)
//	if sourceID == "" {
//		return url
//	}
//
//	return url + "#" + sourceID
//}
