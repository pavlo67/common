package importer

//// Object forms an items.Object from the imported entity
// func (entity Entity) Object() (obj *things.Object, err error) {
//	if entity.item == nil {
//		return nil, importer.ErrNilItem
//	}
//
//	item := entity.item
//
//	//language := ""
//	//if entity.rss != nil {
//	//	language = entity.rss.language
//	//}
//
//	createdLinks := []things.Link{}
//	if item.Author != nil {
//		//email, err := url.Parse(item.Author.Email)
//		//if err != nil {
//		//	email = nil
//		//}
//
//		createdLinks = append(createdLinks, things.Link{
//			Type: "author",
//			//Nick:    []items.Text{{Text: item.Author.Nick, Language: language}},
//			//Whereto: email,
//			Title: item.Author.Title,
//			To:   item.Author.Email,
//		})
//	}
//	if item.Link != "" {
//		//URL, err := url.Parse(item.Link)
//		//if err != nil {
//		//	URL = nil
//		//}
//
//		createdLinks = append(createdLinks, things.Link{
//			Type: "url",
//			//Nick:    []items.Text{{Text: item.Link, Language: language}},
//			//Whereto: URL,
//			Title: item.Link,
//			To:   item.Link,
//		})
//	}
//	if item.Image != nil {
//		//URL, err := url.Parse(item.Image.URL)
//		//if err != nil {
//		//	URL = nil
//		//}
//		createdLinks = append(createdLinks, things.Link{
//			Type: "image",
//			//Nick:    []items.Text{{Text: item.Label, Language: language}},
//			//Whereto: URL,
//			Title: item.Title,
//			To:   item.Image.URL,
//		})
//	}
//	for _, category := range item.Categories {
//		createdLinks = append(createdLinks, things.Link{
//			Type: links.TypeTag,
//			//Nick: []items.Text{{Text: category, Language: language}},
//			Title: category,
//		})
//	}
//
//	return &things.Object{
//		//Nick:    []items.Text{{Text: item.Label, Language: language}},
//		//Summary: []items.Text{{Text: item.Description, Language: language}},
//		Title:    item.Title,
//		Contentus: item.Description + " " + item.Contentus,
//		Tags:   createdLinks,
//	}, nil
//
// }
