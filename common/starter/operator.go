package starter

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
)

type Operator interface {
	Name() string
	Prepare(cfg *config.Config, options common.Map) error
	Run(joiner.Operator) error
}

type Component struct {
	Operator
	Options common.Map
	*config.Config
}

func (starter Component) CorrectedOptions(options common.Map) common.Map {
	newOptions := common.Map{}

	for k, v := range starter.Options {
		newOptions[k] = v
	}

	for k, v := range options {
		newOptions[k] = v
	}

	return newOptions
}
