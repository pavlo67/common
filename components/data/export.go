package data

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/selectors"
)

const onExport = "on data.Export()"

func Export(dataOp Operator, afterIDStr string, options *crud.GetOptions) ([]Item, error) {
	// TODO: remove limits
	// if options != nil {
	//	options.Limits = nil
	// }

	afterIDStr = strings.TrimSpace(afterIDStr)

	var term *selectors.Term

	var afterID int
	if afterIDStr != "" {
		var err error
		afterID, err = strconv.Atoi(afterIDStr)
		if err != nil {
			return nil, errors.Errorf("can't strconv.Atoi(%s) for after_id parameter: %s", afterIDStr, err)
		}

		// TODO!!! term with some item's autoincrement if original .ID isn't it (using .ID to find corresponding autoincrement value)
		term = selectors.Binary(selectors.Gt, "id", selectors.Value{afterID})
	}

	// TODO!!! order by some item's autoincrement if original .ID isn't it
	if options == nil {
		options = &crud.GetOptions{OrderBy: []string{"id"}}
	} else {
		options.OrderBy = []string{"id"}
	}

	return dataOp.List(term, options)
}
