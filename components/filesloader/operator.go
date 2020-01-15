package filesloader

import (
	"github.com/pavlo67/workshop/common/joiner"

	"github.com/pavlo67/workshop/components/files"
)

const InterfaceKey joiner.InterfaceKey = "filesloader"

type Operator interface {
	Load(urlToLoad, pathToStore string, priority Priority) (*files.Item, error)
}
