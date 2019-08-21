package old

import (
	"log"
	"strings"

	"github.com/pavlo67/constructor/components/auth"
	"github.com/pavlo67/constructor/components/common"
	"github.com/pavlo67/constructor/components/common/strlib"
	"github.com/pavlo67/partes/crud/selectors"

	"github.com/pavlo67/constructor/confidenter/groups"

	"github.com/pavlo67/constructor/notebook/links"
)

func PrepareLinks(userIS auth.ID, grpsOp groups.Operator, objectROwner auth.ID, linksListOld, linksListNew []links.Item, jointLinks bool, linkTypesOnly ...string) (linkslistFinal []links.Item) {
	// only userIS's links can be changed
	// to remove group's links user should get corresponding group's confidenter for the action

	var linksListFinal []links.Item

LINKS_OLD:
	for _, l := range linksListOld {
		if l.ROwner == "" {
			l.ROwner = objectROwner
		}

		if l.ROwner != userIS && !jointLinks {
			// keep alien tags if all they aren't joint
			linksListFinal = append(linksListFinal, l)

		} else if len(linkTypesOnly) <= 0 {
			// remove all old own tags if no types selected
			continue

		} else {
			// remove old own tags with selected types only
			for _, linkType := range linkTypesOnly {
				if l.Type == linkType {
					continue LINKS_OLD
				}
			}
			linksListFinal = append(linksListFinal, l)

		}
	}

	for _, l := range linksListNew {
		if l.ROwner == "" || !groups.OneOf(userIS, grpsOp, l.ROwner) {
			l.ROwner = userIS
		}

		if l.RView == "" || !groups.OneOf(userIS, grpsOp, l.RView) {
			// to prevent unexpected link invisibility
			l.RView = common.Anyone
		}

		linksListFinal = append(linksListFinal, l)
	}

	return linksListFinal
}

func AddTags(userIS auth.ID, o *Item, tags string) {
	if o == nil {
		return
	}

TAG:
	for _, t := range strlib.ReSemicolon.Split(tags, -1) {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}

		for _, l := range o.Links {
			if l.Type == links.TypeTag && strings.ToLower(l.Name) == strings.ToLower(t) {
				continue TAG
			}
		}
		o.Links = append(o.Links, links.Item{Type: links.TypeTag, Name: t})
	}
	o.Tags = tags
}

func Linked(userIS auth.ID, objectsOp Operator, linksOp links.Operator, o *Item) []Item {
	if o == nil { // || o.CountLinked <= 0
		return nil
	}

	linked, err := linksOp.QueryByObjectID(userIS, o.ID)
	if err != nil {
		log.Print(err)
		return nil
	}

	if len(linked) < 1 {
		return nil
	}

	var ids []interface{}
	for _, l := range linked {
		ids = append(ids, l.LinkedID)
	}
	selector := selectors.FieldEqual("id", ids...)
	options := &content.ListOptions{Selector: selector, SortBy: []string{"name"}}
	linkedObjects, _, err := objectsOp.ReadList(userIS, options)
	if err != nil {
		log.Print(err)
	}

	return linkedObjects
}
