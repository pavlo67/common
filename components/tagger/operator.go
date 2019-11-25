package tagger

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/components/crud"
)

const InterfaceKey joiner.InterfaceKey = "tagger"

type Tag string

type TagInfo struct {
	Tag
	Count uint64
}

type Tagged map[joiner.InterfaceKey][]common.ID

type Operator interface {
	Save(joiner.InterfaceKey, common.ID, []Tag, *crud.SaveOptions) error
	Remove(joiner.InterfaceKey, common.ID, []Tag, *crud.SaveOptions) error
	Replace(joiner.InterfaceKey, common.ID, []Tag, *crud.SaveOptions) error

	Tags(joiner.InterfaceKey, common.ID, *crud.GetOptions) ([]Tag, error)
	ListTagged(Tag, *crud.GetOptions) (Tagged, error)
}
