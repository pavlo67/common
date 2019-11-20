package links

////LinkedType string `bson:"linked_type,omitempty" json:"linked_type"`
////LinkedID   string `bson:"linked_id,omitempty"   json:"linked_id"`
//
//
//const FieldTag = "tag"
//const FieldType = "type"
//

//func Filter(userIS common.ID, grpsOp groups.Operator, linksList []Item) (linksListFiltered []Item) {
//	for _, l := range linksList {
//		if groups.OneOf(userIS, grpsOp, l.RView, l.ROwner) {
//			linksListFiltered = append(linksListFiltered, l)
//		}
//	}
//	return linksListFiltered
//}

//func SelectString(userIS common.ID, linksOp Operator, linkType string) (selects viewshtml.SelectString) {
//	linked, err := linksOp.Query(userIS, selectors.FieldStr(FieldType, linkType))
//	if err != nil {
//		log.Println(err)
//	}
//
//	var index = make(map[string]string)
//	var names []string
//
//	for _, l := range linked {
//		names = append(names, l.Tag)
//		index[l.Tag] = l.LinkedID
//	}
//
//	sort.Strings(names)
//
//	for _, name := range names {
//		selects = append(selects, [2]string{name, index[name]})
//	}
//
//	return selects
//}

//func Correct(is common.ID, linksListTmp []Item) (linksList []Item) {
//	for _, l := range linksListTmp {
//		if l.ROwner == is { // l.Type == TypeTag &&
//			l.Title = strlib.ReSpaces.ReplaceAllString(strings.TrimSpace(l.Title), " ")
//			l.To = strings.TrimSpace(l.To)
//			linksList = append(linksList, l)
//		}
//	}
//
//	return linksList
//}
