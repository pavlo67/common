package tagger

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"
)

const InterfaceKey joiner.InterfaceKey = "tagger"
const CleanerInterfaceKey joiner.InterfaceKey = "tag_cleaner"

type Tag struct {
	Label  string
	Params common.Map
}

type TagCount struct {
	Label     string
	Immediate uint64
	Full      uint64
}

type Tagged struct {
	ID     common.ID
	Params common.Map
}

type Index map[joiner.InterfaceKey][]Tagged

// TODO: don't remove "...Tags", it's necessary to resolve conflict in data_tagged.Actor

type Operator interface {
	AddTags(joiner.Link, []Tag, *crud.SaveOptions) error
	ReplaceTags(joiner.Link, []Tag, *crud.SaveOptions) error // or remove in particlar
	ListTags(joiner.Link, *crud.GetOptions) ([]Tag, error)   // i.e. parent sections if joiner.HandlerKey == "tagger"

	CountTags(*joiner.InterfaceKey, *crud.GetOptions) ([]TagCount, error)
	IndexTagged(*joiner.InterfaceKey, string, *crud.GetOptions) (Index, error)
}

//func CheckCycle(userIS auth.ID, operator Actor, id string, passedIDs []string) ([]string, error) {
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
