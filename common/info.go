package common

type Info map[string]string

func (p Info) StringDefault(key, defaultStr string) string {
	// log.Printf("00000000: %T %#v", p[key], p[key])

	if v, ok := p[key]; ok {
		return v
	}

	return defaultStr
}
