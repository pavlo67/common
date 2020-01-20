package data

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/selectors"
)

const DataItemsTypeKey crud.TypeKey = ""

const KeyFieldName = "key"

const onExport = "on data.Export(): "

func Export(dataOp Operator, selector *selectors.Term, options *crud.GetOptions, l logger.Operator) (*crud.Data, error) {
	if dataOp == nil {
		return nil, errors.New(onExport + "no data.Operator")
	}

	// TODO!!! order by some item's autoincrement if original .Key isn't it
	if options == nil {
		options = &crud.GetOptions{OrderBy: []string{"id"}}
	} else {
		options.OrderBy = []string{"id"}
	}

	items, err := dataOp.List(selector, options)
	if err != nil {
		return nil, errors.Wrap(err, onExport)
	}

	content, err := json.Marshal(items)
	if err != nil {
		// TODO!!! cut the error message if it's too long
		return nil, errors.Wrapf(err, onExport+"can't .Marshal(%#v)", items)
	}

	return &crud.Data{TypeKey: DataItemsTypeKey, Content: content}, nil
}

const onImport = "on data.Import(): "

func Import(dataOp Operator, data crud.Data, options *crud.SaveOptions, l logger.Operator) error {
	if dataOp == nil {
		return errors.New(onImport + "no data.Operator")
	}

	if data.TypeKey != DataItemsTypeKey {
		return errors.Errorf(onImport+"wrong data.TypeKey(%s)", data.TypeKey)
	}

	var items []Item
	err := json.Unmarshal(data.Content, &items)
	if err != nil {
		return errors.Wrapf(err, onImport+"can't .Unmarshal(%s) into []data.Item", data.Content)
	}

	for i, item := range items {
		itemsOld, err := dataOp.List(selectors.Binary(selectors.Eq, KeyFieldName, selectors.Value{item.Key}), nil)
		if err != nil {
			return errors.Wrapf(err, onImport+"can't get old item for key '%s'", item.Key)
		}
		if len(itemsOld) != 1 {
			return errors.Errorf(onImport+"%d old items for key '%s' (instead one)", len(itemsOld))
		}

		// TODO: check .History

		item.ID = itemsOld[0].ID
		_, err = dataOp.Save(item, options)
		if err != nil {
			return errors.Wrapf(err, onImport+"can't save new item for key '%s'", item.Key)
		}

		l.Infof("imported item %d of %d: %s", i+1, len(items), item.Key)
	}

	return nil
}

//afterIDStr = strings.TrimSpace(afterIDStr)
//
//var selectorID, selectorFull *selectors.Term
//
//var afterID int
//if afterIDStr != "" {
//	var err error
//	afterID, err = strconv.Atoi(afterIDStr)
//	if err != nil {
//		return nil, errors.Errorf("can't strconv.Atoi(%s) for after_id parameter: %s", afterIDStr, err)
//	}
//
//	// TODO!!! selectorID with some item's autoincrement if original .Key isn't it (using .Key to find corresponding autoincrement value)
//	selectorID = selectors.Binary(selectors.Gt, "id", selectors.Value{afterID})
//}
//
//if selectorID == nil {
//	selectorFull = selector
//} else if selector == nil {
//	selectorFull = selectorID
//} else {
//	selector = logic.AND(selector, selectorID)
//}
