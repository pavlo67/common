package sources

//import (
//	"encoding/json"
//	"time"
//
//	"github.com/pkg/errors"
//
//	"github.com/pavlo67/partes/crud"
//	"github.com/pavlo67/punctum/basis"
//	"github.com/pavlo67/punctum/joiner"
//	"github.com/pavlo67/punctum/confidenter/rights"
//)
//
//const InterfaceKeyCRUD joiner.InterfaceKey = "sources.crud"
//
//var _ crud.Operator = &OperatorCRUD{}
//
//type OperatorCRUD struct {
//	Operator
//}
//
//func (opCRUD OperatorCRUD) Describe() (crud.Description, error) {
//	return crud.Description{
//		Title: "sources",
//		FieldsArr: []crud.Field{
//			{Key: "id", Unique: true, AutoUnique: true},
//			{Key: "url", Creatable: true, Updatable: true},
//			{Key: "title", Creatable: true, Updatable: true},
//			{Key: "type", Creatable: true, Updatable: true},
//			{Key: "params_raw", Creatable: true, Updatable: true},
//
//			{Key: "r_view", Creatable: true, Updatable: true, NotEmpty: true},
//			{Key: "r_owner", Creatable: true, Updatable: true, NotEmpty: true},
//			{Key: "managers", Creatable: true, Updatable: true},
//
//			{Key: "created_at", NotEmpty: true},
//			{Key: "updated_at"},
//		},
//		SortByDefault: []string{"id"},
//		Exemplar:      &Item{},
//	}, nil
//}
//
//func (opCRUD OperatorCRUD) StringMapToNative(data crud.StringMap) (interface{}, error) {
//	if data == nil {
//		return nil, basis.ErrNull
//	}
//
//	var managers rights.Managers
//	if data["managers"] != "" {
//		err := json.Unmarshal([]byte(data["managers"]), &managers)
//		if err != nil {
//			return nil, errors.Wrapf(err, `can't json.Unmarshal([]byte(data["managers"]: %s)`, data["managers"])
//		}
//	}
//
//	var params basis.Options
//	if data["params_raw"] != "" {
//		err := json.Unmarshal([]byte(data["params_raw"]), &params)
//		if err != nil {
//			return nil, errors.Wrapf(err, `can't json.Unmarshal([]byte(data["params"]: %s)`, data["params"])
//		}
//	}
//
//	var createdAt time.Time
//	if data["created_at"] != "" {
//		var err error
//		createdAt, err = time.Parse(time.RFC3339, data["created_at"])
//		if err != nil {
//			return nil, errors.Wrapf(err, `can't parse time from data["created_at"]: %s`, data["created_at"])
//		}
//	}
//
//	var updatedAtPtr *time.Time
//	if data["updated_at"] != "" {
//		updatedAt, err := time.Parse(time.RFC3339, data["updated_at"])
//		if err != nil {
//			return nil, errors.Wrapf(err, `can't parse time from data["updated_at"]: %s`, data["updated_at"])
//		}
//		updatedAtPtr = &updatedAt
//	}
//
//	return &Item{
//		ID:        data["id"],
//		URL:       data["url"],
//		Title:     data["title"],
//		Type:      joiner.InterfaceKey(data["type"]),
//		Options:    params,
//		RView:     auth.ID(data["r_view"]),
//		ROwner:    auth.ID(data["r_owner"]),
//		Managers:  managers,
//		SavedAt: createdAt,
//		UpdatedAt: updatedAtPtr,
//	}, nil
//}
//
//func (opCRUD OperatorCRUD) NativeToStringMap(native interface{}) (crud.StringMap, error) {
//	if native == nil {
//		return nil, basis.ErrNull
//	}
//
//	source, ok := native.(*Item)
//	if !ok {
//		sourceItem, ok := native.(Item)
//		if !ok {
//			return nil, basis.ErrWrongDataType
//		}
//		source = &sourceItem
//	}
//
//	var managersJSON []byte
//	if len(source.Managers) > 0 {
//		var err error
//		managersJSON, err = json.Marshal(source.Managers)
//		if err != nil {
//			return nil, errors.Wrapf(err, `can't json.Marshal(obj.Managers): %#v)`, source.Managers)
//		}
//	}
//
//	createdAtStr := source.SavedAt.Format(time.RFC3339)
//
//	var updatedAtStr string
//	if source.UpdatedAt != nil {
//		updatedAtStr = source.UpdatedAt.Format(time.RFC3339)
//	}
//
//	return crud.StringMap{
//		"id":         source.ID,
//		"url":        source.URL,
//		"title":      source.Title,
//		"type":       string(source.Type),
//		"params_raw": source.ParamsRaw,
//		"r_view":     string(source.RView),
//		"r_owner":    string(source.ROwner),
//		"managers":   string(managersJSON),
//		"created_at": createdAtStr,
//		"updated_at": updatedAtStr,
//	}, nil
//}
//
//const onIDFromNative = "on sources.OperatorCRUD.IDFromNative()"
//
//func (opCRUD OperatorCRUD) IDFromNative(native interface{}) (string, error) {
//	item, ok := native.(*Item)
//	if !ok {
//		return "", errors.Wrapf(basis.ErrWrongDataType, onIDFromNative+": expected crud.NativeMap, actual = %T", native)
//	}
//
//	return item.ID, nil
//}
//
//func (opCRUD OperatorCRUD) Create(userIS auth.ID, native interface{}) (id string, err error) {
//	source, ok := native.(*Item)
//	if !ok {
//		return "", basis.ErrWrongDataType
//	}
//	if source == nil {
//		return "", basis.ErrNull
//	}
//
//	return opCRUD.Operator.Create(userIS, *source)
//}
//
//func (opCRUD OperatorCRUD) Read(userIS auth.ID, id string) (interface{}, error) {
//	return opCRUD.Operator.Read(userIS, id)
//}
//
//func (opCRUD OperatorCRUD) ReadList(userIS auth.ID, options *content.ListOptions) ([]interface{}, uint64, error) {
//	srcList, allCnt, err := opCRUD.Operator.ReadList(userIS, options)
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
//	source, ok := native.(*Item)
//	if !ok {
//		return crud.Result{}, basis.ErrWrongDataType
//	}
//	if source == nil {
//		return crud.Result{}, basis.ErrNull
//	}
//
//	return opCRUD.Operator.Update(userIS, *source)
//
//}
//
//func (opCRUD OperatorCRUD) TestCases(cleaner crud.Cleaner) ([]crud.OperatorTestCase, error) {
//
//	userIS := auth.ID("a/b/c")
//	userISAnother := auth.ID("d/e/f")
//	userISNil := auth.ID("")
//
//	toCreatePrivate := crud.StringMap{
//		"url":        "url1",
//		"title":      "title1",
//		"type":       "type1",
//		"params_raw": `{"a":"1"}`,
//		"r_view":     string(userIS),
//		"r_owner":    string(userIS),
//	}
//
//	toUpdatePrivate := crud.StringMap{
//		"url":        "url1u",
//		"title":      "title1u",
//		"type":       "type1u",
//		"params_raw": `{"a":"1u"}`,
//		"r_view":     string(userIS),
//		"r_owner":    string(userIS),
//	}
//
//	toCreatePublic := crud.StringMap{
//		"url":        "url2",
//		"title":      "title2",
//		"type":       "type2",
//		"params_raw": `{"a":"2"}`,
//		"r_view":     string(basis.Anyone),
//		"r_owner":    string(userIS),
//	}
//
//	toUpdatePublic := crud.StringMap{
//		"url":        "url2u",
//		"title":      "title2u",
//		"type":       "type2u",
//		"params_raw": `{"a":"2u"}`,
//		"r_view":     string(basis.Anyone),
//		"r_owner":    string(userIS),
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
