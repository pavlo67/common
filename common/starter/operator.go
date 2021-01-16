package starter

import (
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/data"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
)

type Operator interface {

	// Title returns started component name
	Name() string

	Init(cfg *config.Config, l logger.Operator, options data.Map) (info []data.Map, err error)

	// Setup sets up the component
	Setup() error

	// Run inits the component to use in application
	Run(joiner.Operator) error
}

type Starter struct {
	Operator
	Options data.Map
}

func (starter Starter) CorrectedOptions(options data.Map) data.Map {
	newOptions := data.Map{}

	for k, v := range starter.Options {
		newOptions[k] = v
	}

	for k, v := range options {
		newOptions[k] = v
	}

	return newOptions
}

//type Runner interface {
//	Prepare() error
//}
