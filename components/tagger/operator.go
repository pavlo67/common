package tagger

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"
)

const InterfaceKey joiner.InterfaceKey = "tagger"
const CleanerInterfaceKey joiner.InterfaceKey = "tagger_cleaner"

type Tag struct {
	Label    string
	Relation string
}

type TaggedCount struct {
	Immediate uint64
	Full      uint64
}

type Tagged struct {
	ID       common.ID
	Relation string
}

type Counter map[string]TaggedCount // parted by Tag.Label

type Index map[joiner.InterfaceKey][]Tagged

type Operator interface {
	AddTags(joiner.InterfaceKey, common.ID, []Tag, *crud.SaveOptions) error
	ReplaceTags(joiner.InterfaceKey, common.ID, []Tag, *crud.SaveOptions) error // or remove in particlar
	ListTags(joiner.InterfaceKey, common.ID, *crud.GetOptions) ([]Tag, error)   // i.e. parent sections if joiner.InterfaceKey == "tagger"

	//RemoveTags(joiner.InterfaceKey, common.ID, []string, *crud.SaveOptions) error
	//CleanTags(joiner.InterfaceKey, *selectors.Term, *crud.SaveOptions) error

	CountTagged(*joiner.InterfaceKey, *crud.GetOptions) (Counter, error)
	IndexWithTag(string, *crud.GetOptions) (Index, error)
}

//func CheckCycle(userIS auth.ID, operator Operator, id string, passedIDs []string) ([]string, error) {
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
