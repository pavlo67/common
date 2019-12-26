package filesloader

import (
	"github.com/pavlo67/workshop/common/joiner"

	"github.com/pavlo67/workshop/components/files"
)

const InterfaceKey joiner.InterfaceKey = "filesloader"

type Operator interface {
	Load(pathToLoad, pathToStore string) (*files.Item, error)
}
