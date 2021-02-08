package starter

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
)

type Operator interface {

	// Title returns started component name
	Name() string

	Init(cfg *config.Config, options common.Map) error

	// Run inits the component to use in application
	Run(joiner.Operator) error
}

type Starter struct {
	Operator
	Options common.Map
}

func (starter Starter) CorrectedOptions(options common.Map) common.Map {
	newOptions := common.Map{}

	for k, v := range starter.Options {
		newOptions[k] = v
	}

	for k, v := range options {
		newOptions[k] = v
	}

	return newOptions
}
