package basis

import (
	"strconv"
)

type Info map[string]interface{}

func (p Info) StringDefault(key, defaultStr string) string {
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

func (p Info) String(key string) (string, bool) {
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

func (p Info) Strings(key string) []string {
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

//func (p Info) StringMapKeyDefault(key string, defaultMap map[string]string) map[string]string {
//	if valueMap, ok := p[key].(map[string]string); ok {
//		return valueMap
//	}
//
//	return defaultMap
//}
//
func (p Info) StringsMap() map[string]string {
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
