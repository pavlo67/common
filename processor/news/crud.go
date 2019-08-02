package news

import (
	"github.com/pavlo67/constructor/starter/joiner"
)

const InterfaceKeyCRUD joiner.ComponentKey = "flow.crud"

//var _ crud.Operator = &OperatorCRUD{}
//
//type OperatorCRUD struct {
//	Operator
//}
//
//func (opCRUD OperatorCRUD) Describe() (crud.Description, error) {
//	return crud.Description{
//		Title: "flow",
//		Fields: []crud.Field{
//			{Key: "id", Unique: true, AutoUnique: true},
//			{Key: "source_url", Creatable: true},
//			{Key: "source_time", Creatable: true},
//			{Key: "source_key", Creatable: true},
//
//			{Key: "original", Creatable: true},
//
//			{Key: "title", Creatable: true},
//			{Key: "summary", Creatable: true},
//			{Key: "text", Creatable: true},
//			{Key: "tags", Creatable: true},
//			{Key: "embedded", Creatable: true},
//			{Key: "href", Creatable: true},
//			{Key: "content_key", Creatable: true},
//
//			{Key: "status", Creatable: true},
//			{Key: "history"},
//			{Key: "stored_at"},
//		},
//		SortByDefault: []string{"id"},
//		Exemplar:      &Item{},
//		ViewList: []viewshtml.Field{
//			{Key: "id", Type: "view"},
//
//			{Key: "title", Type: "view"},
//			{Key: "summary", Type: "view"},
//			{Key: "text", Type: "view"},
//			{Key: "tags", Type: "view"},
//			{Key: "href", Type: "view"},
//
//			{Key: "source_url", Type: "view"},
//			{Key: "source_time", Type: "view"},
//			{Key: "source_key", Type: "view"},
//			{Key: "stored_at", Type: "view"},
//			{Key: "status"},
//		},
//	}, nil
//}
//
//func (opCRUD OperatorCRUD) StringMapToNative(dataMap crud.StringMap) (interface{}, error) {
//	if dataMap == nil {
//		return nil, basis.ErrNull
//	}
//
//	var embedded []Embedded
//	if dataMap["embedded"] != "" {
//		err := json.Unmarshal([]byte(dataMap["embedded"]), &embedded)
//		if err != nil {
//			return nil, errors.Wrapf(err, `can't json.Unmarshal([]byte(data["embedded"]: %s)`, dataMap["embedded"])
//		}
//	}
//
//	var tags []string
//	if dataMap["tags"] != "" {
//		err := json.Unmarshal([]byte(dataMap["tags"]), &tags)
//		if err != nil {
//			return nil, errors.Wrapf(err, `can't json.Unmarshal([]byte(data["tags"]: %s)`, dataMap["tags"])
//		}
//	}
//
//	var storedAtPtr *time.Time
//	if dataMap["stored_at"] != "" {
//		storedAt, err := time.Parse(time.RFC3339, dataMap["stored_at"])
//		if err != nil {
//			return nil, errors.Wrapf(err, `can't parse time from data["stored_at"]: %s`, dataMap["stored_at"])
//		}
//		storedAtPtr = &storedAt
//	}
//
//	var sourceTimePtr *time.Time
//	if dataMap["source_time"] != "" {
//		sourceTime, err := time.Parse(time.RFC3339, dataMap["source_time"])
//		if err != nil {
//			return nil, errors.Wrapf(err, `can't parse time from data["source_time"]: %s`, dataMap["source_time"])
//		}
//		sourceTimePtr = &sourceTime
//	}
//
//	return &Item{
//		ID: dataMap["id"],
//		Source: flow.Source{
//			URL:  dataMap["source_url"],
//			Time: sourceTimePtr,
//			Key:  dataMap["source_key"],
//		},
//		Original:   dataMap["original"],
//		ContentKey: dataMap["content_key"],
//		Content: &Content{
//			Title:    dataMap["title"],
//			Summary:  dataMap["summary"],
//			Text:     dataMap["text"],
//			Tags:     tags,
//			Embedded: embedded,
//			Href:     dataMap["href"],
//		},
//		Status:   dataMap["status"],
//		History:  dataMap["history"],
//		SavedAt: storedAtPtr,
//	}, nil
//}
//
//func (opCRUD OperatorCRUD) NativeToStringMap(native interface{}) (crud.StringMap, error) {
//	if native == nil {
//		return nil, basis.ErrNull
//	}
//
//	itemPtr, ok := native.(*Item)
//	if !ok {
//		item, ok := native.(Item)
//		if !ok {
//			return nil, basis.ErrWrongDataType
//		}
//		itemPtr = &item
//	}
//
//	var embeddedJSON []byte
//	if len(itemPtr.Embedded) > 0 {
//		var err error
//		embeddedJSON, err = json.Marshal(itemPtr.Embedded)
//		if err != nil {
//			return nil, errors.Wrapf(err, `can't json.Marshal(obj.Embedded): %#v)`, itemPtr.Embedded)
//		}
//	}
//
//	var tagsJSON []byte
//	if len(itemPtr.Tags) > 0 {
//		var err error
//		tagsJSON, err = json.Marshal(itemPtr.Tags)
//		if err != nil {
//			return nil, errors.Wrapf(err, `can't json.Marshal(obj.Tags): %#v)`, itemPtr.Tags)
//		}
//	}
//
//	var storedAtStr string
//	if itemPtr.SavedAt != nil {
//		storedAtStr = itemPtr.SavedAt.Format(time.RFC3339)
//	}
//
//	var sourceTimeStr string
//	if itemPtr.Source.Time != nil {
//		sourceTimeStr = itemPtr.Source.Time.Format(time.RFC3339)
//	}
//
//	return crud.StringMap{
//		"id":          itemPtr.ID,
//		"source_url":  itemPtr.Source.URL,
//		"source_time": sourceTimeStr,
//		"source_key":  itemPtr.Source.Key,
//
//		"original": itemPtr.Original,
//
//		"title":       itemPtr.Content.Title,
//		"summary":     itemPtr.Content.Summary,
//		"text":        itemPtr.Content.Text,
//		"embedded":    string(embeddedJSON),
//		"tags":        string(tagsJSON),
//		"href":        itemPtr.Content.Href,
//		"content_key": itemPtr.ContentKey,
//
//		"status":    itemPtr.Status,
//		"history":   itemPtr.History,
//		"stored_at": storedAtStr,
//	}, nil
//}
//
//const onIDFromNative = "on sources.OperatorCRUD.IDFromNative()"
//
//func (opCRUD OperatorCRUD) IDFromNative(native interface{}) (string, error) {
//	item, ok := native.(*Item)
//	if !ok {
//		return "", errors.Wrapf(basis.ErrWrongDataType, onIDFromNative+": expected *flow.Item, actual = %T", native)
//	}
//
//	return item.ID, nil
//}
//
//func (opCRUD OperatorCRUD) Create(userIS auth.ID, native interface{}) (id string, err error) {
//	flow, ok := native.(*Item)
//	if !ok {
//		return "", basis.ErrWrongDataType
//	}
//	if flow == nil {
//		return "", basis.ErrNull
//	}
//
//	return "", opCRUD.Operator.Save(flow)
//}
//
//func (opCRUD OperatorCRUD) Read(userIS auth.ID, id string) (interface{}, error) {
//	items, _, err := opCRUD.Operator.ReadList(
//		&content.ListOptions{
//			Selector: selectors.FieldStr("id", id),
//		},
//	)
//	if err != nil {
//		return nil, err
//	}
//	if len(items) != 1 {
//		return nil, errors.Errorf("wrong number of items: %d", len(items))
//
//	}
//
//	return items[0], nil
//}
//
//func (opCRUD OperatorCRUD) ReadList(userIS auth.ID, options *content.ListOptions) ([]interface{}, uint64, error) {
//	srcList, allCnt, err := opCRUD.Operator.ReadList(options)
//	if err != nil {
//		return nil, 0, err
//	}
//
//	var intfsList []interface{}
//	for _, src := range srcList {
//		intfsList = append(intfsList, src)
//	}
//
//	return intfsList, allCnt, nil
//}
//
//func (opCRUD OperatorCRUD) Update(userIS auth.ID, native interface{}) (crud.Result, error) {
//	return crud.Result{}, basis.ErrNotImplemented
//}
//
//func (opCRUD OperatorCRUD) Delete(userIS auth.ID, id string) (crud.Result, error) {
//	return opCRUD.Operator.Delete(
//		&content.ListOptions{
//			Selector: selectors.FieldStr("id", id),
//		},
//	)
//}
//
////func (opCRUD OperatorCRUD) TestCases(cleaner crud.Cleaner) ([]crud.OperatorTestCase, error) {
////
////	userIS := basis.UserIS("a/b/c")
////	userISAnother := basis.UserIS("d/e/f")
////	userISNil := basis.UserIS("")
////
////	toCreatePrivate := crud.StringMap{
////		"url":        "url1",
////		"title":      "title1",
////		"type":       "type1",
////		"params_raw": `{"a":"1"}`,
////		"r_view":     string(userIS),
////		"r_owner":    string(userIS),
////	}
////
////	toUpdatePrivate := crud.StringMap{
////		"url":        "url1u",
////		"title":      "title1u",
////		"type":       "type1u",
////		"params_raw": `{"a":"1u"}`,
////		"r_view":     string(userIS),
////		"r_owner":    string(userIS),
////	}
////
////	toCreatePublic := crud.StringMap{
////		"url":        "url2",
////		"title":      "title2",
////		"type":       "type2",
////		"params_raw": `{"a":"2"}`,
////		"r_view":     string(basis.Anyone),
////		"r_owner":    string(userIS),
////	}
////
////	toUpdatePublic := crud.StringMap{
////		"url":        "url2u",
////		"title":      "title2u",
////		"type":       "type2u",
////		"params_raw": `{"a":"2u"}`,
////		"r_view":     string(basis.Anyone),
////		"r_owner":    string(userIS),
////	}
////
////	testCases := []crud.OperatorTestCase{
////
////		// 0. all ok for private record,
////		// can't create with identityNil,
////		// can't read, update or delete with identityAnother
////		{
////			Operator: opCRUD,
////			Cleaner:  cleaner,
////
////			ISToCreate:        userIS,
////			ISToCreateBad:     &userISNil,
////			ToSave:          toCreatePrivate,
////			ExpectedSaveErr: nil,
////
////			ISToRead:        userIS,
////			ISToReadBad:     &userISAnother,
////			ExpectedReadErr: nil,
////
////			ISToUpdate:        userIS,
////			ISToUpdateBad:     &userISAnother,
////			ToUpdate:          toUpdatePrivate,
////			ExpectedUpdateErr: nil,
////
////			ISToDelete:        userIS,
////			ISToDeleteBad:     &userISAnother,
////			ExpectedRemoveErr: nil,
////		},
////
////		// 1. all ok for private record,
////		// can't create with identityNil,
////		// can't read, update or delete with identityNil
////		{
////			Operator: opCRUD,
////			Cleaner:  cleaner,
////
////			ISToCreate:        userIS,
////			ISToCreateBad:     &userISNil,
////			ToSave:          toCreatePrivate,
////			ExpectedSaveErr: nil,
////
////			ISToRead:        userIS,
////			ISToReadBad:     &userISNil,
////			ExpectedReadErr: nil,
////
////			ISToUpdate:        userIS,
////			ISToUpdateBad:     &userISNil,
////			ToUpdate:          toUpdatePrivate,
////			ExpectedUpdateErr: nil,
////
////			ISToDelete:        userIS,
////			ISToDeleteBad:     &userISNil,
////			ExpectedRemoveErr: nil,
////		},
////
////		// 2. all ok for public record,
////		// can't create with identityNil,
////		// can read with identityAnother
////		// can't update or delete with identityAnother
////		{
////			Operator: opCRUD,
////			Cleaner:  cleaner,
////
////			ISToCreate:        userIS,
////			ISToCreateBad:     &userISNil,
////			ToSave:          toCreatePublic,
////			ExpectedSaveErr: nil,
////
////			ISToRead:        userIS,
////			ExpectedReadErr: nil,
////
////			ISToUpdate:        userIS,
////			ISToUpdateBad:     &userISAnother,
////			ToUpdate:          toUpdatePublic,
////			ExpectedUpdateErr: nil,
////
////			ISToDelete:        userIS,
////			ISToDeleteBad:     &userISAnother,
////			ExpectedRemoveErr: nil,
////		},
////
////		// 3. all ok for public record,
////		// can't create with identityNil,
////		// can read with identityNil,
////		// can't update or delete with identityNil
////		// close database
////		{
////			Operator: opCRUD,
////			Cleaner:  cleaner,
////
////			ISToCreate:        userIS,
////			ISToCreateBad:     &userISNil,
////			ToSave:          toCreatePublic,
////			ExpectedSaveErr: nil,
////
////			ISToRead:        userISNil,
////			ExpectedReadErr: nil,
////
////			ISToUpdate:        userIS,
////			ISToUpdateBad:     &userISNil,
////			ToUpdate:          toUpdatePublic,
////			ExpectedUpdateErr: nil,
////
////			ISToDelete:        userIS,
////			ISToDeleteBad:     &userISNil,
////			ExpectedRemoveErr: nil,
////		},
////	}
////
////	return testCases, nil
////}
