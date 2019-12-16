package data_tagged

import (
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/selectors"
	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/hypertext"
	"github.com/pavlo67/workshop/components/tagger"
)

const InterfaceKey joiner.InterfaceKey = "data_tagged"

type Tagger = tagger.Operator // to use data.Operator and tagger.Operator simultaneously in Operator interface

type Operator interface {
	data.Operator
	Tagger
	ListWithTag(*selectors.Term, string, *crud.GetOptions) ([]data.Item, error)
	ListWithText(*selectors.Term, hypertext.ToSearch, *crud.GetOptions) ([]data.Item, error)
}
