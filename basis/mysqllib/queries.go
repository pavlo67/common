package mysqllib

import (
	"strconv"

	"github.com/pkg/errors"
)

const CantGetRowsAffected = "can't .RowsAffected (sql='%s', values='%#v')"
const NoRowOnQuery = "no row on query(sql='%s', values='%#v'"
const CantScanQueryRow = "can't scan query row (sql='%s', values='%#v')"

var ErrNoTable = errors.New("table doesn't exist")

const defaultPageLengthStr = "200"

func OrderAndLimit(sortBy []string, limits []uint64) string {
	var sortStr, limitsStr string
	if len(sortBy) > 0 {
		for _, s := range sortBy {
			if s == "" {
				continue
			}
			desc := ""
			if s[len(s)-1:] == "-" {
				s = s[:len(s)-1]
				desc = " desc"
			} else if s[len(s)-1:] == "+" {
				s = s[:len(s)-1]
			}
			if sortStr != "" {
				sortStr += ", "
			}
			sortStr += "`" + s + "`" + desc
		}
		if sortStr != "" {
			sortStr = " order by " + sortStr
		}
	}
	if len(limits) > 1 {
		// limit[0] can be equal to 0
		var pageLengthStr string
		if limits[1] > 0 {
			pageLengthStr = strconv.FormatUint(limits[1], 10)
		} else {
			pageLengthStr = defaultPageLengthStr
		}
		limitsStr = " limit " + strconv.FormatUint(limits[0], 10) + ", " + pageLengthStr
	} else if len(limits) > 0 {
		if limits[0] > 0 {
			limitsStr = " limit " + strconv.FormatUint(limits[0], 10)
		} else {
			limitsStr = " limit " + defaultPageLengthStr
		}
	}
	return sortStr + limitsStr
}
