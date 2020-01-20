package crud

import "github.com/pavlo67/workshop/common/selectors"

type Exporter interface {
	Export(selector *selectors.Term, options *GetOptions) (*Data, error)
	Import(data Data, options *SaveOptions) error
}

func ExporterTestScenario(exporterOp Exporter) {

}
