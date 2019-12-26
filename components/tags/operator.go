package tags

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"
)

const InterfaceKey joiner.InterfaceKey = "tags"
const CleanerInterfaceKey joiner.InterfaceKey = "tags_cleaner"

type Item struct {
	Label    string
	Relation string
}

type TagCount struct {
	Label     string
	Immediate uint64
	Full      uint64
}

type Tagged struct {
	ID       common.Key
	Relation string
}

type Index map[joiner.InterfaceKey][]Tagged

// TODO: don't remove "...Tags", it's necessary to resolve conflict in data_tagged.Operator

type Operator interface {
	AddTags(joiner.InterfaceKey, common.Key, []Item, *crud.SaveOptions) error
	ReplaceTags(joiner.InterfaceKey, common.Key, []Item, *crud.SaveOptions) error // or remove in particlar

	ListTags(joiner.InterfaceKey, common.Key, *crud.GetOptions) ([]Item, error) // i.e. parent sections if joiner.InterfaceKey == "tagger"
	CountTags(*joiner.InterfaceKey, *crud.GetOptions) ([]TagCount, error)

	IndexTagged(*joiner.InterfaceKey, string, *crud.GetOptions) (Index, error)
}

//func CheckCycle(userIS auth.Key, operator Operator, id string, passedIDs []string) ([]string, error) {
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
