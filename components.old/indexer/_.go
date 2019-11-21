//func (item *Item) PartesTexti() ([]textus.Pars, error) {
//	if item == nil || item.Content == nil {
//		return nil, basis.ErrNull
//	}
//
//	return []textus.Pars{
//		{
//			Fons:            item.OriginData,
//			Origo:           item.Original,
//			ClavisContentus: item.ContentKey,
//			Contentus: &textus.Contentus{
//				Titulus:    item.Content.Title,
//				Index:      item.Content.Summary,
//				Textus:     item.Content.Text,
//				Appendices: map[string][]string{"tags": item.Content.Tags},
//			},
//		},
//	}, nil
//
//}

//func (src *OriginData) Key(keyAdd string) string {
//	if src == nil {
//		return ""
//	}
//
//	url := strings.TrimSpace(src.URL)
//
//	pos := strings.Index(url, "#")
//	if pos >= 0 {
//		url = url[:pos]
//	}
//
//	if url == "" {
//		return ""
//	}
//
//	if len(keyAdd) > 0 {
//		url += "#" + keyAdd
//	}
//
//	sourceID := strings.TrimSpace(src.ID)
//	if sourceID == "" {
//		return url
//	}
//
//	return url + "#" + sourceID
//}
