package basis

type Values []struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (values Values) ByName(key string) string {
	for _, v := range values {

		if v.Key == key {
			return v.Value
		}
	}

	return ""
}
