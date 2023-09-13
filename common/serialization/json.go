package serialization

import (
	"encoding/json"
	"github.com/pavlo67/common/common/errors"
)

const onJSONList = "on format.JSONList()"

func JSONList[T any](list []T, prefix, indent string) ([]byte, error) {
	prefixBytes, indentBytes := []byte(prefix), []byte(indent)

	jsonBytes := []byte{'[', '\n'}

	for i, item := range list {
		itemBytes, err := json.Marshal(item)
		if err != nil {
			return nil, errors.Wrap(err, onJSONList)
		}

		if i > 0 {
			jsonBytes = append(append(jsonBytes, ','), prefixBytes...)
		}
		jsonBytes = append(append(jsonBytes, indentBytes...), itemBytes...)
	}

	return append(jsonBytes, '\n', ']'), nil

}
