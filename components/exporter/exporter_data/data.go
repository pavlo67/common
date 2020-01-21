package exporter_data

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/selectors"
	"github.com/pavlo67/workshop/common/selectors/logic"

	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/exporter"
)

var _ exporter.Operator = &exporterData{}

type exporterData struct {
	dataOp       data.Operator
	interfaceKey joiner.InterfaceKey
}

const onNew = "on exporterData.New(): "

func New(dataOp data.Operator, interfaceKey joiner.InterfaceKey) (exporter.Operator, error) {
	if dataOp == nil {
		return nil, errors.New(onNew + "no data.Operator")
	}

	return &exporterData{dataOp, interfaceKey}, nil
}

const onExport = "on data.Export(): "

func (expOp *exporterData) Export(selector *selectors.Term, afterIDStr string, options *crud.GetOptions) (*crud.Data, error) {
	afterIDStr = strings.TrimSpace(afterIDStr)

	var selectorID, selectorFull *selectors.Term

	var afterID uint64
	if afterIDStr != "" {
		var err error
		afterID, err = strconv.ParseUint(afterIDStr, 10, 64)
		if err != nil {
			return nil, errors.Errorf("can't strconv.Atoi(%s) for after_id parameter: %s", afterIDStr, err)
		}

		// TODO!!! use other item's autoincrement if original .ID isn't autoincremental (using .Key to find corresponding autoincrement value)
		selectorID = selectors.Binary(selectors.Gt, "id", selectors.Value{afterID})
	}

	if selectorID == nil {
		selectorFull = selector
	} else if selector == nil {
		selectorFull = selectorID
	} else {
		selector = logic.AND(selector, selectorID)
	}

	// TODO!!! order by some other item's autoincrement if original .ID isn't autoincremental
	if options == nil {
		options = &crud.GetOptions{OrderBy: []string{"id"}}
	} else {
		options.OrderBy = []string{"id"}
	}

	items, err := expOp.dataOp.List(selectorFull, options)
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

func (expOp *exporterData) Import(crudData crud.Data, options *crud.SaveOptions) (string, error) {
	if crudData.TypeKey != DataItemsTypeKey {
		return "", errors.Errorf(onImport+"wrong crudData.TypeKey(%s)", crudData.TypeKey)
	}

	var items []data.Item
	err := json.Unmarshal(crudData.Content, &items)
	if err != nil {
		return "", errors.Wrapf(err, onImport+"can't .Unmarshal(%s) into []crudData.Item", crudData.Content)
	}

	var tillIDStr string
	if len(items) > 0 {
		for i, item := range items {
			itemsOld, err := expOp.dataOp.List(selectors.Binary(selectors.Eq, KeyFieldName, selectors.Value{item.Key}), nil)
			if err != nil {
				return "", errors.Wrapf(err, onImport+"can't get old item for key '%s'", item.Key)
			}
			if len(itemsOld) != 1 {
				return "", errors.Errorf(onImport+"%d old items for key '%s' (instead one)", len(itemsOld))
			}

			item.ID = itemsOld[0].ID
			_, err = expOp.dataOp.Save(item, options)
			if err != nil {
				return "", errors.Wrapf(err, onImport+"can't save new item for key '%s'", item.Key)
			}

			l.Infof("imported item %d of %d: %s", i+1, len(items), item.Key)
		}

		tillIDStr = string(items[len(items)-1].ID)
	}

	return tillIDStr, nil
}
