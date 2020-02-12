package data

import (
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/selectors"
	"github.com/pavlo67/workshop/common/selectors/logic"
	"github.com/pavlo67/workshop/components/tagger"
)

const onListWithTag = "on dataTagged.ListTagged(): "

func ListTagged(dataOp Operator, taggerOp tagger.Operator, dataKey *joiner.InterfaceKey, tagLabel string, selector *selectors.Term, options *crud.GetOptions) ([]Item, error) {
	if dataOp == nil {
		return nil, errors.New(onListWithTag + ": no data.Operator")
	}

	if taggerOp == nil {
		return nil, errors.New(onListWithTag + ": no tagger.Operator")
	}

	index, err := taggerOp.IndexTagged(dataKey, tagLabel, options)
	if err != nil {
		return nil, errors.Wrap(err, onListWithTag)
	}
	var taggedIDs []interface{}
	for _, i := range index {
		for _, tagged := range i {
			taggedIDs = append(taggedIDs, string(tagged.ID))
		}
	}

	selectorTagged := selectors.In("id", taggedIDs...)
	if selector != nil {
		selectorTagged = logic.AND(selectorTagged, selector)
	}

	// l.Infof("%#v\n%#v", selectorTagged, options)
	return dataOp.List(selectorTagged, options)

	// TODO: check if all item.TypeKey are correct in the result of wsOp.ListTags
}

//// Search -----------------------------------------------------------------------------------------
//
//var rePhrase = regexp.MustCompile(`^\s*".*"\s*$`)
//var reDelimiter = regexp.MustCompile(`[\.,\s\t;:\-\+\!\?\(\)\{\}\[\]\/'"\*]+`)
//
//func (objOp *notesMySQL) ReadListByWords(userIS common.Key, options *content.ListOptions, searched string) (objects []notes.Tag, allCnt uint64, err error) {
//	if !rePhrase.MatchString(searched) {
//		words := reDelimiter.Split(searched, -1)
//		searched = ""
//		for _, w := range words {
//			if len(w) > 2 {
//				searched += " +" + w
//			}
//		}
//	}
//
//	selectorSearched := selector.Match("name,content,tags", searched, "IN BOOLEAN MODE")
//	if options == nil {
//		options = &content.ListOptions{Selector: selectorSearched}
//	} else if options.Selector == nil {
//		options.Selector = selectorSearched
//	} else {
//		options.Selector = selector.And(options.Selector, selectorSearched)
//	}
//	return objOp.ReadList(userIS, options)
//}
//

//
//// UpdateLinks corrects object links within and without object itself
//
//const onUpdateLinks = "on notesMySQL.UpdateLinks"
//
//func (objOp *notesMySQL) UpdateLinks(userIS common.Key, idStr string, linksListNew []links.Tag, linkType string) error {
//	// TODO: lock object record for update (use history!!!)
//
//	o, err := objOp.Read(userIS, idStr)
//	if err != nil {
//		return errors.Wrap(err, onUpdateLinks)
//	}
//
//	linksList := notes.PrepareLinks(userIS, objOp.grpsOp, o.ROwner, o.Tags, linksListNew, objOp.jointLinks, linkType)
//	linksListCopy := linksList
//	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//
//	var jsonLinks []byte
//	if len(linksList) > 0 {
//		jsonLinks, err = json.Marshal(linksList)
//		if err != nil {
//			return errors.Wrapf(err, onUpdateLinks+": can't marshal .Tags(%#v)", linksList)
//		}
//	}
//
//	values := []interface{}{jsonLinks, o.Key}
//	_, err = objOp.stmtUpdateLinks.Exec(values...)
//	if err != nil {
//		return errors.Wrapf(err, onUpdateLinks+": "+basis.CantExecQuery, objOp.sqlUpdateLinks, values)
//	}
//
//	return objOp.setLinks(userIS, o.Key, linksListCopy)
//}
