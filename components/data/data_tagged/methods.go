package data_tagged

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/selectors"
	"github.com/pavlo67/workshop/common/selectors/logic"
	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/hypertext"
	"github.com/pavlo67/workshop/components/tags"
)

var _ Operator = &ws{}

// var _ crud.Cleaner = &ws{}

type ws struct {
	data.Operator
	Tagger
}

const onNewWorkspace = "on New(): "

func New(dataOp data.Operator, taggerOp tags.Operator) (Operator, crud.Cleaner, error) {
	if dataOp == nil {
		return nil, nil, errors.New(onNewWorkspace + ": no data.Operatoe")
	}

	wsOp := ws{
		Operator: dataOp,
		Tagger:   taggerOp,
	}
	return &wsOp, nil, nil
}

const onListWithTag = "on ws.ListWithTag(): "

func (wsOp *ws) ListWithTag(key *joiner.InterfaceKey, tagLabel string, selector *selectors.Term, options *crud.GetOptions) ([]data.Item, error) {
	if wsOp.Tagger == nil {
		return nil, errors.New(onListWithTag + ": no tagger.Operator")
	}

	index, err := wsOp.IndexTagged(key, tagLabel, options)
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
		selectorTagged = logic.AND(selectorTagged, *selector)
	}

	// l.Infof("%#v\n%#v", selectorTagged, options)
	return wsOp.List(selectorTagged, options)

	// TODO: check if all item.TypeKey are correct in the result of wsOp.ListTags
}

const onListWithText = "on ws.ListWithText(): "

func (wsOp *ws) ListWithText(*joiner.InterfaceKey, hypertext.ToSearch, *selectors.Term, *crud.GetOptions) ([]data.Item, error) {
	return nil, common.ErrNotImplemented
}

//// Search -----------------------------------------------------------------------------------------
//
//var rePhrase = regexp.MustCompile(`^\s*".*"\s*$`)
//var reDelimiter = regexp.MustCompile(`[\.,\s\t;:\-\+\!\?\(\)\{\}\[\]\/'"\*]+`)
//
//func (objOp *notesMySQL) ReadListByWords(userIS common.Key, options *content.ListOptions, searched string) (objects []notes.Item, allCnt uint64, err error) {
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
//func (objOp *notesMySQL) UpdateLinks(userIS common.Key, idStr string, linksListNew []links.Item, linkType string) error {
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
