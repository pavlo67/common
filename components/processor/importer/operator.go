package importer

import (
	"github.com/pavlo67/constructor/components/basis/joiner"

	"time"

	"github.com/pavlo67/constructor/components/basis"
)

const InterfaceKey joiner.InterfaceKey = "importer"

//var ErrNoFount = errors.New("no source is reachable")
//var ErrNoMoreItems = errors.New("no more items.comp")
//var ErrBadItemID = errors.New("bad item id")
//var ErrBadItem = errors.New("bad item")
//var ErrNilItem = errors.New("item is nil")

type Item struct {
	ID      basis.ID   `bson:"_id,omitempty"      json:"id,omitempty"`
	SavedAt *time.Time `bson:"saved_at,omitempty" json:"saved_at,omitempty"`

	OriginKey `       bson:",inline"          json:",inline"`
	Origin    string `bson:"origin,omitempty" json:"origin,omitempty"`

	SourceURL  string     `bson:"source_url,omitempty"  json:"source_url,omitempty"`
	SourceTime *time.Time `bson:"source_time,omitempty" json:"source_time,omitempty"`

	Content  `          bson:"content"            json:"content"`
	Embedded []Content `bson:"embedded,omitempty" json:"embedded,omitempty"`

	Tags []string `bson:"tags,omitempty"   json:"tags,omitempty"`
}

type OriginKey struct {
	SourceID  basis.ID `bson:"source_id,omitempty"  json:"source_id,omitempty"`
	SourceKey string   `bson:"source_key,omitempty" json:"source_key,omitempty"`
}

type Content struct {
	Title   string `bson:"title"                json:"title"`
	Summary string `bson:"summary,omitempty"    json:"summary,omitempty"`
	Details string `bson:"details,omitempty"    json:"details,omitempty"`
	Href    string `bson:"href,omitempty"       json:"href,omitempty"`
}

type Operator interface {
	// Run opens import session with selected data source
	// Init() error

	Get(url string, minKey *string) ([]Item, error)
}

//func Run(impOp Operator, sources []string, newsOp news.Operator) error {
//	var errs basis.Errors
//	var cnt int
//
//	for _, src := range sources {
//		err := impOp.Run(src)
//		if err != nil {
//			errs = append(errs, errors.Wrapf(err, "on impOp.Run(%s)", src))
//			continue
//		}
//
//		for {
//			if cnt%100 == 0 {
//				fmt.Println(cnt)
//			}
//			cnt++
//
//			item, err := impOp.Next()
//			if err != nil {
//				errs = errs.Append(err)
//				log.Printf("on impOp.Next(): %s", err)
//			}
//			if item == nil {
//				continue
//			}
//
//			//if staged {
//			//	err = dsOpStaged.SaveStaged(*item)
//			//} else {
//			//	err = newsOp.Save(*item)
//			//}
//
//			err = newsOp.Save(item)
//
//			if err != nil {
//				errs = errs.Append(err)
//				log.Printf("on dsOp.Save(%#v): %s", *item, err)
//			}
//		}
//	}
//	//err := dsOp.Commit(nil)
//	//if err != nil {
//	//	errs = errs.Append(err)
//	//	log.Printf("on dsOp.Commit(nil): %s", err)
//	//}
//
//	return errs.Err()
//}
