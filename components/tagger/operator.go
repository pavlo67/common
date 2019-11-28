package tagger

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/selectors"
)

const InterfaceKey joiner.InterfaceKey = "tagger"

type Tag string

type Operator interface {
	SaveTags(joiner.InterfaceKey, common.ID, []Tag, *crud.SaveOptions) error
	ReplaceTags(joiner.InterfaceKey, common.ID, []Tag, *crud.SaveOptions) error
	ListTags(joiner.InterfaceKey, common.ID, *crud.GetOptions) ([]Tag, error)
	CountTagged(*joiner.InterfaceKey, *crud.GetOptions) (crud.Counter, error)

	RemoveTags(joiner.InterfaceKey, common.ID, []Tag, *crud.SaveOptions) error
	CleanTags(joiner.InterfaceKey, *selectors.Term, *crud.SaveOptions) error

	IndexWithTag(Tag, *crud.GetOptions) (crud.Index, error)
}
