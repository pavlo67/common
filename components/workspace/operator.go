package workspace

import (
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/selectors"
	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/tagger"
	"github.com/pavlo67/workshop/components/text"
)

const InterfaceKey joiner.InterfaceKey = "workspace"

type Tagger = tagger.Operator

type Operator interface {
	data.Operator
	Tagger
	ListWithTag(*selectors.Term, tagger.Tag, *crud.GetOptions) ([]data.Item, error)
	ListWithText(*selectors.Term, text.ToSearch, *crud.GetOptions) ([]data.Item, error)
}
