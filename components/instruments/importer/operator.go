package importer

import (
	"github.com/pavlo67/workshop/basis/joiner"

	"github.com/pavlo67/workshop/components/data"
)

const InterfaceKey joiner.InterfaceKey = "importer"

//var ErrNoFount = errors.New("no source is reachable")
//var ErrNoMoreItems = errors.New("no more items.comp")
//var ErrBadItemID = errors.New("bad item id")
//var ErrBadItem = errors.New("bad item")
//var ErrNilItem = errors.New("item is nil")

type Series struct {
	URL string
	// MinKey *string
	Items []data.Item
}

type Operator interface {
	// Run opens import session with selected data source
	// Init() error

	Get(url string, minKey *string) (*Series, error)
}
