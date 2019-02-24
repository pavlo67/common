package convertor

import (
	"github.com/pavlo67/punctum/notebook/notes"
	"github.com/pavlo67/punctum/processor.old/news"
)

type Factory interface {
	Init(url string, params interface{}) (Operator, error)
}

// Operator is an abstraction of external entity that can be converted into different our ones.
type Operator interface {
	// Original returns the original entity "as is".
	Original() (interface{}, error)

	// Object returns an object the original entity is converted into.
	Object() (*notes.Item, error)

	// FlowItem returns an flow.Census the original entity is converted into.
	FlowItem() (*news.Item, error)

	// Files returns a filer.comp list the original entity is converted into.
	Files() ([]files.File, error)
}
