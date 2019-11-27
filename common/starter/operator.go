package starter

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
)

type Operator interface {

	// Title returns started component name
	Name() string

	Init(conf *config.Config, l logger.Operator, options common.Map) (info []common.Map, err error)

	// Setup sets up the component
	Setup() error

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

//type Runner interface {
//	Prepare() error
//}