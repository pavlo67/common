package basis

import (
	"strconv"
)

type Options map[string]interface{}

func (p Options) StringDefault(key, defaultStr string) string {
	switch value := p[key].(type) {
	case string:
		return value
	case []string:
		if len(value) > 0 {
			return value[0]
		}
	case int:
		return strconv.Itoa(value)
	}

	return defaultStr
}

func (p Options) String(key string) (string, bool) {
	switch value := p[key].(type) {
	case string:
		return value, true
	case []string:
		if len(value) > 0 {
			return value[0], true
		}
	case int:
		return strconv.Itoa(value), true
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
	}

	return nil
}

//func (p Options) StringMapKeyDefault(key string, defaultMap map[string]string) map[string]string {
//	if valueMap, ok := p[key].(map[string]string); ok {
//		return valueMap
//	}
//
//	return defaultMap
//}
//
func (p Options) StringsMap() map[string]string {
	data := map[string]string{}

	for k, v := range p {
		if vStr, ok := v.(string); ok {
			data[k] = vStr
		} else if vStrs, ok := v.([]string); ok && len(vStrs) > 0 {
			data[k] = vStrs[0]
		}
	}

	return data
}
