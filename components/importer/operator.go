package importer

import (
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/transport"

	"time"

	"github.com/pavlo67/workshop/components/data"
)

const InterfaceKey joiner.InterfaceKey = "importer"

//var ErrNoFount = errors.New("no source is reachable")
//var ErrNoMoreItems = errors.New("no more items.comp")
//var ErrBadItemID = errors.New("bad item id")
//var ErrBadItem = errors.New("bad item")
//var ErrNilItem = errors.New("item is nil")

type DataSeries struct {
	URL       string
	CreatedAt time.Time

	Type  transport.DataType
	Data  []data.Item
	MaxID string
}

type Operator interface {
	// Prepare opens import session with selected data source
	// Init() error

	Get(url string) (*DataSeries, error)
}
