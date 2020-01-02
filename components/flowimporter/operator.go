package flowimporter

import (
	"time"

	"github.com/pavlo67/workshop/common/joiner"

	"github.com/pavlo67/workshop/components/data"
)

const InterfaceKey joiner.InterfaceKey = "importer"

type DataSeries struct {
	URL       string
	CreatedAt time.Time
	Data      []data.Item
}

type Operator interface {
	// Prepare opens import session with selected data source
	// Init() error

	Get(key string) (*DataSeries, error)
}
