package datatagged

import (
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/selectors"
	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/hypertext"
	"github.com/pavlo67/workshop/components/tagger"
)

const InterfaceKey joiner.InterfaceKey = "data_tagged"

type Tagger = tagger.Operator // to use data.Operator and tagger.Operator simultaneously in ActorKey interface

type Operator interface {
	data.Operator
	Tagger
	ListTagged(string, *selectors.Term, *crud.GetOptions) ([]data.Item, error)
	ListWithText(hypertext.ToSearch, *selectors.Term, *crud.GetOptions) ([]data.Item, error)
}
