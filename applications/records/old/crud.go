package old

//var _ crud.Operator = &OperatorCRUD{}
//
//type OperatorCRUD struct {
//	Operator
//}
//
//func (opCRUD OperatorCRUD) Describe() (crud.Description, error) {
//	return crud.Description{
//		FieldsArr: []crud.Field{
//			{Key: "genus", Creatable: true, Editable: true},
//
//			{Key: "author", Creatable: true, Editable: true},
//			{Key: "name", Creatable: true, Editable: true},
//			{Key: "brief", Creatable: true, Editable: true},
//			{Key: "content", Creatable: true, Editable: true},
//			{Key: "links", Creatable: true, Editable: true},
//			{Key: "tags", Creatable: true, Editable: true},
//			{Key: "visibility", Creatable: true, Editable: true},
//			{Key: "count_linked"},
//
//			{Key: "r_view", Creatable: true, Editable: true, NotEmpty: true},
//			{Key: "r_owner", Creatable: true, Editable: true, NotEmpty: true},
//			{Key: "managers", Creatable: true, Editable: true},
//
//			{Key: "created_at", NotEmpty: true},
//			{Key: "updated_at"},
//
//			{Key: "global_is", Creatable: true},
//			{Key: "history"},
//			{Key: "status", Creatable: true, Editable: true},
//		},
//	}, nil
//}
//
//func DataToObject(data crud.Contentus) (*items.Item, error) {
//	if data == nil {
//		return nil, basis.ErrNullItem
//	}
//
//	var err error
//
//	var linksList []items.Link
//	if data["links"] != "" {
//		err = json.Unmarshal([]byte(data["links"]), &linksList)
//		if err != nil {
//			return nil, errors.Wrapf(err, `can't json.Unmarshal([]byte(data["links"]: %s)`, data["linksList"])
//		}
//	}
//
//	var managers rights.Managers
//	if data["managers"] != "" {
//		err = json.Unmarshal([]byte(data["managers"]), &managers)
//		if err != nil {
//			return nil, errors.Wrapf(err, `can't json.Unmarshal([]byte(data["managers"]: %s)`, data["managers"])
//		}
//	}
//
//	var createdAt, updatedAt time.Time
//	var updatedAtPtr *time.Time
//	if data["created_at"] != "" {
//		createdAt, err = time.Parse(time.RFC3339, data["created_at"])
//		if err != nil {
//			return nil, errors.Wrapf(err, `can't parse time from data["created_at"]: %s`, data["created_at"])
//		}
//	}
//	if data["updated_at"] != "" {
//		updatedAt, err = time.Parse(time.RFC3339, data["updated_at"])
//		if err != nil {
//			return nil, errors.Wrapf(err, `can't parse time from data["updated_at"]: %s`, data["updated_at"])
//		}
//		updatedAtPtr = &updatedAt
//	}
//
//	countLinked, _ := strconv.Atoi(data["count_linked"])
//
//	return &items.Item{
//		TargetID:          data["id"],
//		Genus:       data["genus"],
//		Author:      data["author"],
//		Title:        data["name"],
//		Brief:       data["brief"],
//		Contentus:     data["content"],
//		Links:       linksList,
//		Tags:        data["tags"],
//		Visibility:  data["visibility"],
//		CountLinked: uint16(countLinked),
//		RView:       basis.UserIS(data["r_view"]),
//		ROwner:      basis.UserIS(data["r_owner"]),
//		Managers:    managers,
//		SavedAt:   createdAt,
//		UpdatedAt:   updatedAtPtr,
//		GlobalIS:    data["global_is"],
//		History:     data["history"],
//		Status:      data["status"],
//	}, nil
//}
//
//func ObjectToData(obj *items.Item) (crud.Contentus, error) {
//	if obj == nil {
//		return nil, basis.ErrNullItem
//	}
//
//	var err error
//
//	var jsonLinks []byte
//	if obj.Links != nil {
//		jsonLinks, err = json.Marshal(obj.Links)
//		if err != nil {
//			return nil, errors.Wrapf(err, `can't json.Marshal(obj.Links): %#v)`, obj.Links)
//		}
//	}
//
//	var jsonManagers []byte
//	if len(obj.Managers) > 0 {
//		jsonManagers, err = json.Marshal(obj.Managers)
//		if err != nil {
//			return nil, errors.Wrapf(err, `can't json.Marshal(obj.Managers): %#v)`, obj.Managers)
//		}
//	}
//
//	createdAt := obj.SavedAt.Format(time.RFC3339)
//
//	var updatedAt string
//	if obj.UpdatedAt != nil {
//		updatedAt = obj.UpdatedAt.Format(time.RFC3339)
//	}
//
//	return crud.Contentus{
//		"id":           obj.TargetID,
//		"genus":        obj.Genus,
//		"author":       obj.Author,
//		"name":         obj.Title,
//		"brief":        obj.Brief,
//		"content":      obj.Contentus,
//		"links":        string(jsonLinks),
//		"tags":         obj.Tags,
//		"visibility":   obj.Visibility,
//		"count_linked": string(obj.CountLinked),
//		"r_view":       string(obj.RView),
//		"r_owner":      string(obj.ROwner),
//		"managers":     string(jsonManagers),
//		"created_at":   createdAt,
//		"updated_at":   updatedAt,
//		"global_is":    obj.GlobalIS,
//		"history":      obj.History,
//		"status":       obj.Status,
//	}, nil
//}
//
//func (opCRUD OperatorCRUD) Create(userIS basis.UserIS, data crud.Contentus) (id string, err error) {
//	obj, err := DataToObject(data)
//	if err != nil {
//		return "", err
//	}
//
//	return opCRUD.Operator.Create(userIS, *obj)
//}
//
//func (opCRUD OperatorCRUD) Read(userIS basis.UserIS, id string) (crud.Contentus, error) {
//	obj, err := opCRUD.Operator.Read(userIS, id)
//	if err != nil {
//		return nil, err
//	}
//
//	return ObjectToData(obj)
//}
//
//func (opCRUD OperatorCRUD) ReadList(userIS basis.UserIS, options *content.ListOptions) ([]crud.Contentus, uint64, error) {
//	objList, allCnt, err := opCRUD.Operator.ReadList(userIS, options)
//	if err != nil {
//		return nil, 0, err
//	}
//
//	var dataList []crud.Contentus
//	for _, obj := range objList {
//		data, err := ObjectToData(&obj)
//		if err != nil {
//			return dataList, allCnt, err
//		}
//		dataList = append(dataList, data)
//	}
//
//	return dataList, allCnt, nil
//}
//
//func (opCRUD OperatorCRUD) Update(userIS basis.UserIS, data crud.Contentus) (crud.Result, error) {
//	obj, err := DataToObject(data)
//	if err != nil {
//		return crud.Result{}, err
//	}
//
//	return opCRUD.Operator.Update(userIS, *obj)
//}
//
//func (opCRUD OperatorCRUD) TestCases(cleaner crud.Cleaner) ([]crud.OperatorTestCase, error) {
//	userIS := basis.UserIS("a/b/c")
//	userISAnother := basis.UserIS("d/e/f")
//	userISNil := basis.UserIS("")
//
//	name := "Nick One"
//
//	linksToCreatePrivate, _ := json.Marshal([]items.Link{
//		{TargetID: "1", Type: files.LinkType, Title: "Nick file test", ROwner: userIS, RView: userIS},
//		{TargetID: "Name2 file2 test2", Type: "tag", Title: "Name2 file2 test2", ROwner: userIS, RView: userIS},
//	})
//
//	toCreatePrivate := crud.Contentus{
//		"genus":    "test obj",
//		"language": "en",
//		"author":   "a",
//		"name":     string(name),
//		"summary":  "summary test",
//		"content":  "content test",
//		"links":    string(linksToCreatePrivate),
//		"tags":     "tag1;tag2;tag3;",
//		"r_view":   string(userIS),
//		"r_owner":  string(userIS),
//		// "managers":  "{}",
//		"global_is": "aaa/bb/3",
//		"history":   "123",
//		"status":    "",
//		"original":  "789",
//	}
//
//	toUpdatePrivate := crud.Contentus{
//		"genus":    "test obj1",
//		"language": "ua",
//		"author":   "a1",
//		"name":     string(name) + "1",
//		"summary":  "summary test1",
//		"content":  "content test1",
//		"links":    "",
//		"tags":     "tag11;tag2;tag3;",
//		"r_view":   string(userIS),
//		"r_owner":  string(userIS),
//		// "managers":  "{}",
//		"global_is": "aaa/bb/4",
//		"history":   "456",
//		"status":    "1",
//		"original":  "",
//	}
//
//	linksToCreatePublic, _ := json.Marshal([]items.Link{
//		{TargetID: "1", Type: files.LinkType, Title: "Nick file test", ROwner: userIS, RView: basis.Anyone},
//		{TargetID: "Name2 file2 test2", Type: "tag", Title: "Name2 file2 test2", ROwner: userIS, RView: basis.Anyone},
//	})
//
//	toCreatePublic := crud.Contentus{
//		"genus":    "test obj",
//		"language": "en",
//		"author":   "a",
//		"name":     string(name) + "2",
//		"summary":  "summary test",
//		"content":  "content test",
//		"links":    string(linksToCreatePublic),
//		"tags":     "",
//		"r_view":   string(basis.Anyone),
//		"r_owner":  string(userIS),
//		// "managers":  "{}",
//		"global_is": "aaa/bb/3",
//		"history":   "123",
//		"status":    "",
//		"original":  "789",
//	}
//
//	toUpdatePublic := crud.Contentus{
//		"genus":    "test obj1",
//		"language": "ua",
//		"author":   "a1",
//		"name":     string(name) + "3",
//		"summary":  "summary test1",
//		"content":  "content test1",
//		"links":    "",
//		"tags":     "tag4;tag5;tag3;",
//		"r_view":   string(basis.Anyone),
//		"r_owner":  string(userIS),
//		// "managers":  "{}",
//		"global_is": "aaa/bb/4",
//		"history":   "456",
//		"status":    "1",
//		"original":  "",
//	}
//
//	testCases := []crud.OperatorTestCase{
//
//		// 0. all ok for private record,
//		// can't create with identityNil,
//		// can't read, update or delete with identityAnother
//		{
//			Operator: opCRUD,
//			Cleaner:  cleaner,
//
//			ISToCreate:        userIS,
//			ISToCreateBad:     &userISNil,
//			ToSave:          toCreatePrivate,
//			ExpectedSaveErr: nil,
//
//			ISToRead:        userIS,
//			ISToReadBad:     &userISAnother,
//			ExpectedReadErr: nil,
//
//			ISToUpdate:        userIS,
//			ISToUpdateBad:     &userISAnother,
//			ToUpdate:          toUpdatePrivate,
//			ExpectedUpdateErr: nil,
//
//			ISToDelete:        userIS,
//			ISToDeleteBad:     &userISAnother,
//			ExpectedRemoveErr: nil,
//		},
//
//		// 1. all ok for private record,
//		// can't create with identityNil,
//		// can't read, update or delete with identityNil
//		{
//			Operator: opCRUD,
//			Cleaner:  cleaner,
//
//			ISToCreate:        userIS,
//			ISToCreateBad:     &userISNil,
//			ToSave:          toCreatePrivate,
//			ExpectedSaveErr: nil,
//
//			ISToRead:        userIS,
//			ISToReadBad:     &userISNil,
//			ExpectedReadErr: nil,
//
//			ISToUpdate:        userIS,
//			ISToUpdateBad:     &userISNil,
//			ToUpdate:          toUpdatePrivate,
//			ExpectedUpdateErr: nil,
//
//			ISToDelete:        userIS,
//			ISToDeleteBad:     &userISNil,
//			ExpectedRemoveErr: nil,
//		},
//
//		// 2. all ok for public record,
//		// can't create with identityNil,
//		// can read with identityAnother
//		// can't update or delete with identityAnother
//		{
//			Operator: opCRUD,
//			Cleaner:  cleaner,
//
//			ISToCreate:        userIS,
//			ISToCreateBad:     &userISNil,
//			ToSave:          toCreatePublic,
//			ExpectedSaveErr: nil,
//
//			ISToRead:        userIS,
//			ExpectedReadErr: nil,
//
//			ISToUpdate:        userIS,
//			ISToUpdateBad:     &userISAnother,
//			ToUpdate:          toUpdatePublic,
//			ExpectedUpdateErr: nil,
//
//			ISToDelete:        userIS,
//			ISToDeleteBad:     &userISAnother,
//			ExpectedRemoveErr: nil,
//		},
//
//		// 3. all ok for public record,
//		// can't create with identityNil,
//		// can read with identityNil,
//		// can't update or delete with identityNil
//		// close database
//		{
//			Operator: opCRUD,
//			Cleaner:  cleaner,
//
//			ISToCreate:        userIS,
//			ISToCreateBad:     &userISNil,
//			ToSave:          toCreatePublic,
//			ExpectedSaveErr: nil,
//
//			ISToRead:        userISNil,
//			ExpectedReadErr: nil,
//
//			ISToUpdate:        userIS,
//			ISToUpdateBad:     &userISNil,
//			ToUpdate:          toUpdatePublic,
//			ExpectedUpdateErr: nil,
//
//			ISToDelete:        userIS,
//			ISToDeleteBad:     &userISNil,
//			ExpectedRemoveErr: nil,
//		},
//	}
//
//	return testCases, nil
//}
//
////func (obj *MySQLObject) CountCRUD(selector selectors.Selector, joinTo crud.JoinTo, groupBy, sortBy []string) ([]crud.Count, error) {
////	if Tables[joinTo.ToTable] != "" {
////		joinTo.ToTable = Tables[joinTo.ToTable]
////	} else {
////		return nil, nil.New("can't find table code: " + joinTo.ToTable)
////	}
////	return mysqllib.Count(obj.dbh, selector, joinTo, groupBy, sortBy)
////}
////
////// Describe ... read crud.json5
////func (obj *MySQLObject) DescribeCRUD() (*crud.Description, error) {
////	return crud.Describe(filelib.CurrentPath() + "../")
////}
