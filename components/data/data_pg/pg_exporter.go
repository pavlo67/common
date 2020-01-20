package data_pg

import (
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/selectors"

	"github.com/pavlo67/workshop/components/data"
)

func (dataOp *dataPg) Export(selector *selectors.Term, options *crud.GetOptions) (*crud.Data, error) {
	return data.Export(dataOp, selector, options, l)
}

func (dataOp *dataPg) Import(crudData crud.Data, options *crud.SaveOptions) error {
	return data.Import(dataOp, crudData, options, l)
}
