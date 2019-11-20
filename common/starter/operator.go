package starter

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
)

type Operator interface {

	// Title returns started component name
	Name() string

	Init(conf *config.Config, options common.Options) (info []common.Options, err error)

	// Setup sets up the component
	Setup() error

	// Run inits the component to use in application
	Run(joiner.Operator) error
}

type Starter struct {
	Operator
	Options common.Options
}

func (starter Starter) CorrectedOptions(options common.Options) common.Options {
	newOptions := common.Options{}

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
