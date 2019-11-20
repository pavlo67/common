package common

import (
	"strconv"
)

type Options map[string]interface{}

func (p Options) StringDefault(key, defaultStr string) string {
	switch value := p[key].(type) {
	case string:
		return value
	case int:
		return strconv.Itoa(value)
	case int64:
		return strconv.FormatInt(value, 10)
	case int32:
		return strconv.Itoa(int(value))
	case int16:
		return strconv.Itoa(int(value))
	case int8:
		return strconv.Itoa(int(value))
	}

	return defaultStr
}

func (p Options) String(key string) (string, bool) {
	switch value := p[key].(type) {
	case string:
		return value, true
	case int:
		return strconv.Itoa(value), true
	case int64:
		return strconv.FormatInt(value, 10), true
	case int32:
		return strconv.Itoa(int(value)), true
	case int16:
		return strconv.Itoa(int(value)), true
	case int8:
		return strconv.Itoa(int(value)), true
	}

	return "", false
}

func (p Options) Strings(key string) []string {
	switch value := p[key].(type) {
	case string:
		return []string{value}
	case []string:
		return value
	case int:
		return []string{strconv.Itoa(value)}
	case int64:
		return []string{strconv.FormatInt(value, 10)}
	case int32:
		return []string{strconv.Itoa(int(value))}
	case int16:
		return []string{strconv.Itoa(int(value))}
	case int8:
		return []string{strconv.Itoa(int(value))}
	}

	return nil
}
