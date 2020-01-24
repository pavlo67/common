package tagger

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"
)

const InterfaceKey joiner.InterfaceKey = "tagger"
const CleanerInterfaceKey joiner.InterfaceKey = "tag_cleaner"

type Tag struct {
	Label  string     `bson:",omitempty" json:",omitempty"`
	Params common.Map `bson:",omitempty" json:",omitempty"`
}

type TagCount struct {
	Label     string `bson:",omitempty" json:",omitempty"`
	Immediate uint64 `bson:",omitempty" json:",omitempty"`
	Full      uint64 `bson:",omitempty" json:",omitempty"`
}

type Tagged struct {
	ID     common.ID  `bson:",omitempty" json:",omitempty"`
	Params common.Map `bson:",omitempty" json:",omitempty"`
}

type Index map[joiner.InterfaceKey][]Tagged

// TODO: don't remove "...Tags", it's necessary to resolve conflict in data_tagged.ActorKey

type Operator interface {
	AddTags(joiner.Link, []Tag, *crud.SaveOptions) error
	ReplaceTags(joiner.Link, []Tag, *crud.SaveOptions) error // or remove in particlar
	ListTags(joiner.Link, *crud.GetOptions) ([]Tag, error)   // i.e. parent sections if joiner.HandlerKey == "tagger"

	CountTags(*joiner.InterfaceKey, *crud.GetOptions) ([]TagCount, error)
	IndexTagged(*joiner.InterfaceKey, string, *crud.GetOptions) (Index, error)
}

//func CheckCycle(userIS auth.Key, operator ActorKey, id string, passedIDs []string) ([]string, error) {
//	for _, passedID := range passedIDs {
//		if id == passedID {
//			return nil, ErrSectionCycle
//		}
//	}
//
//	passedIDs = append(passedIDs, id)
//
//	parentIDs, err := operator.ParentIDs(userIS, id)
//	if err != nil {
//		return nil, errors.Wrapf(err, "can't get section's (%#s) parents IDs", id)
//	}
//
//	for _, parentID := range parentIDs {
//		passedIDs, err = CheckCycle(userIS, operator, parentID, passedIDs)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	return passedIDs, nil
//}
