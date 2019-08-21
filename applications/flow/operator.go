package flow

import (
	"time"

	"github.com/pavlo67/constructor/components/common"
	"github.com/pavlo67/constructor/components/common/joiner"
	"github.com/pavlo67/constructor/components/structura/content"

	"github.com/pavlo67/constructor/components/processor/importer"
)

const InterfaceKey joiner.InterfaceKey = "flow"

type Part struct {
	Ket   string
	Count uint64
}

type Operator interface {
	ListAll(before *time.Time, options *content.GetOptions) ([]content.Brief, error)
	ListByURL(url string, before *time.Time, options *content.GetOptions) ([]content.Brief, error)
	ListByTag(tag string, before *time.Time, options *content.GetOptions) ([]content.Brief, error)
	Read(common.ID, *content.GetOptions) (*importer.Item, error)
	URLs(*content.GetOptions) ([]Part, error)
	Tags(*content.GetOptions) ([]Part, error)
}

type Administrator interface {
	Has(importer.OriginKey) (bool, error)
	Save([]importer.Item, *content.SaveOptions) ([]common.ID, error)
	// Remove(sourceIDs []basis.ID, before *time.Time, options *content.RemoveOptions) error
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
