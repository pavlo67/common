package workspace

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/selectors"
	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/hypertext"
	"github.com/pavlo67/workshop/components/tagger"
)

var _ Operator = &ws{}

// var _ crud.Cleaner = &ws{}

type ws struct {
	data.Operator
	Tagger
}

const onNewWorkspace = "on NewWorkspace(): "

func NewWorkspace(dataOp data.Operator, taggerOp tagger.Operator) (Operator, crud.Cleaner, error) {
	if dataOp == nil {
		return nil, nil, errors.New(onNewWorkspace + ": no data.Operatoe")
	}
	if taggerOp == nil {
		return nil, nil, errors.New(onNewWorkspace + ": no tagger.Operatoe")
	}

	wsOp := ws{
		Operator: dataOp,
		Tagger:   taggerOp,
	}
	return &wsOp, nil, nil
}

const onListWithTag = "on ws.ListWithTag(): "

func (wsOp *ws) ListWithTag(selector *selectors.Term, tag tagger.Tag, options *crud.GetOptions) ([]data.Item, error) {
	_, err := wsOp.IndexWithTag(tag, options)
	if err != nil {
		return nil, errors.Wrap(err, onListWithTag)
	}

	// TODO: modify selector with index

	//	selectorTagged := selector.FieldStr("id", linkedIDs...)
	//	if options == nil {
	//		options = &content.ListOptions{Selector: selectorTagged}
	//	} else if options.Selector == nil {
	//		options.Selector = selectorTagged
	//	} else {
	//		options.Selector = selector.And(options.Selector, selectorTagged)
	//	}
	//
	//	for _, l := range linked {
	//		id, _ := strconv.ParseUint(l.ObjectID, 10, 64)
	//		if id > 0 {
	//			duplicatedID := false
	//			idStr := strings.TrimSpace(l.ObjectID)
	//			for _, parentID := range parentIDs {
	//				if idStr == parentID {
	//					duplicatedID = true
	//					continue
	//				}
	//			}
	//			if !duplicatedID {
	//				parentIDs = append(parentIDs, idStr)
	//			}
	//		}
	//	}
	//
	//	linkedObjs, allCnt, err = objOp.ReadList(userIS, options)
	//
	return wsOp.List(selector, options)

}

const onListWithText = "on ws.ListWithText(): "

func (wsOp *ws) ListWithText(*selectors.Term, hypertext.ToSearch, *crud.GetOptions) ([]data.Item, error) {
	return nil, common.ErrNotImplemented
}

//// Search -----------------------------------------------------------------------------------------
//
//var rePhrase = regexp.MustCompile(`^\s*".*"\s*$`)
//var reDelimiter = regexp.MustCompile(`[\.,\s\t;:\-\+\!\?\(\)\{\}\[\]\/'"\*]+`)
//
//func (objOp *notesMySQL) ReadListByWords(userIS common.ID, options *content.ListOptions, searched string) (objects []notes.Item, allCnt uint64, err error) {
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
//func (objOp *notesMySQL) UpdateLinks(userIS common.ID, idStr string, linksListNew []links.Item, linkType string) error {
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
//	values := []interface{}{jsonLinks, o.ID}
//	_, err = objOp.stmtUpdateLinks.Exec(values...)
//	if err != nil {
//		return errors.Wrapf(err, onUpdateLinks+": "+basis.CantExecQuery, objOp.sqlUpdateLinks, values)
//	}
//
//	return objOp.setLinks(userIS, o.ID, linksListCopy)
//}
