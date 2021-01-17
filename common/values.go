package common

import "strconv"

func Int64(val interface{}) *int64 {
	if val == nil {
		return nil
	}

	switch value := val.(type) {
	case string:
		if val64, err := strconv.ParseInt(value, 10, 64); err == nil {
			return &val64
		}
		return nil
	case int:
		val64 := int64(value)
		return &val64
	case int64:
		return &value
	case int32:
		val64 := int64(value)
		return &val64
	case int16:
		val64 := int64(value)
		return &val64
	case int8:
		val64 := int64(value)
		return &val64
	case uint:
		val64 := int64(value)
		return &val64
	case uint64:
		// TODO!!! check overflow
		val64 := int64(value)
		return &val64
	case uint32:
		val64 := int64(value)
		return &val64
	case uint16:
		val64 := int64(value)
		return &val64
	case uint8:
		val64 := int64(value)
		return &val64
	case float32:
		val64 := int64(value)
		return &val64
	case float64:
		val64 := int64(value)
		return &val64
	}

	return nil
}
