package starter

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
)

type Operator interface {
	Name() string
	Run(*config.Config, common.Map, joiner.Operator, logger.Operator) error
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
