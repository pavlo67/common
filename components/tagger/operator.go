package tagger

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"
)

const InterfaceKey joiner.InterfaceKey = "tagger"

type Tag string

type Operator interface {
	SaveTags(joiner.InterfaceKey, common.ID, []Tag, *crud.SaveOptions) error
	RemoveTags(joiner.InterfaceKey, common.ID, []Tag, *crud.SaveOptions) error
	ReplaceTags(joiner.InterfaceKey, common.ID, []Tag, *crud.SaveOptions) error
	ListTags(joiner.InterfaceKey, common.ID, *crud.GetOptions) ([]Tag, error)

	CountTagged(*joiner.InterfaceKey, *crud.GetOptions) (crud.Counter, error)
	IndexWithTag(Tag, *crud.GetOptions) (crud.Index, error)
}
