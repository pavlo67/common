package exporter

import (
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/selectors"
)

const InterfaceKey joiner.InterfaceKey = "exporter"

type Operator interface {
	Export(selector *selectors.Term, after string, options *crud.GetOptions) (*crud.Data, error)
	Import(crudData crud.Data, options *crud.SaveOptions) (till string, err error)
}

func ExporterTestScenario(exporterOp Operator) {

}
