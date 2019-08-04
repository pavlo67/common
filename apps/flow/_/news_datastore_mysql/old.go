package news_datastore_mysql

//if item.Original != nil {
//	original, err = json.Marshal(item.Original)
//	if err != nil {
//		return auth.IDentity{}, errors.Wrapf(err, "can't marshal item.Original: %v in flow.Create", item.Original)
//	}
//}
//if item.Media != nil {
//	media, err = json.Marshal(item.Media)
//	if err != nil {
//		return auth.IDentity{}, errors.Wrapf(err, "can't marshal item.Media: %v in flow.Create", item.Media)
//	}
//}
//
//if len([]rune(item.Title)) > MaxVarcharLen {
//	item.Title = string([]rune(item.Title)[:MaxVarcharLen])
//}
//if len([]rune(item.Summary)) > MaxVarcharLen {
//	item.Summary = string([]rune(item.Summary)[:MaxVarcharLen])
//}
//
////itemIdentity := itemIS.Identity()
//var item flow.Census
//var media []byte
//err := f.stmtRead.QueryRow(string(itemID)).Scan(&item.TargetID, &item.RView, &item.ROwner, &item.FountIS, &item.OriginalID, &item.SourceURL, &item.URL, &item.Title, &item.Summary, &item.Contentus, &item.Original, &media, &item.SavedAt, &item.ImportedTo)
//if err == sql.ErrNoRows {
//	return nil, errors.New("item not found")
//}
//if err != nil {
//	return nil, errors.Wrapf(err, "can't exec QueryRow: %s, is=%s", f.sqlRead, itemID)
//}
//if string(media) != "" {
//	err = json.Unmarshal(media, &item.Media)
//	if err != nil {
//		return nil, errors.Wrapf(err, "can't unmarshal flow.media: %s, is=%s", media, itemID)
//	}
//}
