package starter

import (
	"github.com/pavlo67/workshop/basis/common"
	"github.com/pavlo67/workshop/basis/config"
	"github.com/pavlo67/workshop/basis/joiner"
)

//type Info struct {
//	Service string
//	Path    string
//	Status  string
//	Details interface{}
//}

type Operator interface {

	// Title returns started component name
	Name() string

	Init(conf *config.Config, options common.Info) (info []common.Info, err error)

	// Setup sets up the component
	Setup() error

	// Run inits the component to use in application
	Run(joiner.Operator) error
}

type Starter struct {
	Operator
	Options common.Info
}

func (starter Starter) CorrectedOptions(options common.Info) common.Info {
	newOptions := common.Info{}

	for k, v := range starter.Options {
		newOptions[k] = v
	}

	for k, v := range options {
		newOptions[k] = v
	}

	return newOptions
}

//type Runner interface {
//	Run() error
//}
