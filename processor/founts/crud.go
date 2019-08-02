package founts

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"

	"github.com/pavlo67/constructor/auth"
	"github.com/pavlo67/constructor/basis"
	"github.com/pavlo67/constructor/crud"
	"github.com/pavlo67/constructor/processor"
	"github.com/pavlo67/constructor/starter/joiner"
)

const InterfaceKeyCRUD joiner.ComponentKey = "founts.crud"

var _ crud.Operator = &OperatorCRUD{}

type OperatorCRUD struct {
	Operator
}

func (opCRUD OperatorCRUD) Description() crud.Description {
	return crud.Description{
		ExemplarNative: &Item{},
		Fields: []crud.Field{
			{Key: "url", Creatable: true, Unique: true, Primary: true},
			{Key: "log", Creatable: true, Updatable: true},
			{Key: "saved_at", NotEmpty: true},
		},
	}
}

func (opCRUD OperatorCRUD) StringMapToNative(data crud.StringMap) (interface{}, error) {
	if data == nil {
		return nil, basis.ErrNull
	}

	var logItems []processor.LogItem
	if data["log"] != "" {
		err := json.Unmarshal([]byte(data["log"]), &logItems)
		if err != nil {
			return nil, errors.Wrapf(err, `can't json.Unmarshal([]byte(data["log"]: %s)`, data["log"])
		}
	}

	var savedAt time.Time
	if data["saved_at"] != "" {
		var err error
		savedAt, err = time.Parse(time.RFC3339, data["saved_at"])
		if err != nil {
			return nil, errors.Wrapf(err, `can't parse time from data["saved_at"]: %s`, data["saved_at"])
		}
	}

	return &Item{
		URL:     data["url"],
		Log:     logItems,
		SavedAt: savedAt,
	}, nil
}

func (opCRUD OperatorCRUD) NativeToStringMap(native interface{}) (crud.StringMap, error) {
	if native == nil {
		return nil, basis.ErrNull
	}

	fount, ok := native.(*Item)
	if !ok {
		fountItem, ok := native.(Item)
		if !ok {
			return nil, basis.ErrWrongDataType
		}
		fount = &fountItem
	}

	var logJSON []byte
	if len(fount.Log) > 0 {
		var err error
		logJSON, err = json.Marshal(fount.Log)
		if err != nil {
			return nil, errors.Wrapf(err, `can't json.Marshal(obj.Log): %#v)`, fount.Log)
		}
	}

	savedAtStr := fount.SavedAt.Format(time.RFC3339)

	return crud.StringMap{
		"url":      fount.URL,
		"log":      string(logJSON),
		"saved_at": savedAtStr,
	}, nil
}

const onIDFromNative = "on founts.OperatorCRUD.IDFromNative()"

func (opCRUD OperatorCRUD) IDFromNative(native interface{}) (string, error) {
	item, ok := native.(*Item)
	if !ok {
		return "", errors.Wrapf(basis.ErrWrongDataType, onIDFromNative+": expected crud.NativeMap, actual = %T", native)
	}

	return item.URL, nil
}

func (opCRUD OperatorCRUD) Create(_ auth.ID, native interface{}) (string, error) {
	fount, ok := native.(*Item)
	if !ok {
		return "", basis.ErrWrongDataType
	}
	if fount == nil {
		return "", basis.ErrNull
	}

	_, err := opCRUD.Operator.Read(fount.URL)
	if err == nil {
		return "", errors.New("already exists")
	} else if err != leveldb.ErrNotFound {
		return "", err
	}

	return fount.URL, opCRUD.Operator.Save(fount.URL, fount.Log...)
}

func (opCRUD OperatorCRUD) Read(_ auth.ID, url string) (interface{}, error) {
	return opCRUD.Operator.Read(url)
}

func (opCRUD OperatorCRUD) ReadList(userIS auth.ID, options content.ListOptions) ([]interface{}, *uint64, error) {
	srcList, allCnt, err := opCRUD.Operator.ReadList(options)
	if err != nil {
		return nil, nil, err
	}

	var intfsList []interface{}
	for _, src := range srcList {
		intfsList = append(intfsList, src)
	}

	return intfsList, allCnt, nil
}

func (opCRUD OperatorCRUD) Update(_ auth.ID, url string, native interface{}) error {
	fount, ok := native.(*Item)
	if !ok {
		return basis.ErrWrongDataType
	}
	if fount == nil {
		return basis.ErrNull
	}

	_, err := opCRUD.Operator.Read(url)
	if err != nil {
		return err
	}

	return opCRUD.Operator.Save(url, fount.Log...)

}

func (opCRUD OperatorCRUD) Delete(_ auth.ID, url string) error {
	return opCRUD.Operator.Delete(url)
}

func (opCRUD OperatorCRUD) TestCases(cleaner crud.Cleaner) ([]crud.OperatorTestCase, error) {

	url := "url1"

	log1, _ := json.Marshal([]processor.LogItem{{
		Started:  time.Now(),
		Finished: time.Now(),
		Success:  false,
		Info:     "???",
	}})

	log2, _ := json.Marshal([]processor.LogItem{{
		Started:  time.Now(),
		Finished: time.Now(),
		Success:  true,
		Info:     "!!!",
	}})

	toCreate := crud.StringMap{
		"url": url,
		"log": string(log1),
	}

	toUpdate := crud.StringMap{
		"url": url,
		"log": string(log2),
	}

	testCases := []crud.OperatorTestCase{

		// 0. all ok
		{
			Operator: opCRUD,
			Cleaner:  cleaner,

			KeyField: "url",

			ToCreate: toCreate,
			ToUpdate: toUpdate,
		},
	}

	return testCases, nil
}
