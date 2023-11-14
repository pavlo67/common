package f

import (
	"encoding/json"
	"fmt"
)

func J(v interface{}) string {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%+v", v)
	}

	return string(jsonBytes)
}
