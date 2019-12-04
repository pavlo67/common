package loader

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/selectors"
)

type Operator interface {
	Select(selector selectors.Term) ([]interface{}, error)
}

const onPack = "on loader.Pack()"

func Pack(op interface{}, selector selectors.Term) ([]byte, error) {
	loaderOp, ok := op.(Operator)
	if !ok {
		return nil, errors.Errorf(onPack+": no loader.Operator with %#v", op)
	}

	data, err := loaderOp.Select(selector)
	if err != nil {
		return nil, errors.Wrapf(err, onPack+": on loaderOp.Select(%#v) with %#v", selector, loaderOp)
	}

	return json.Marshal(data)
}
