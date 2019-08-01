package starter

import (
	"github.com/pavlo67/constructor/basis"
	"github.com/pavlo67/constructor/starter/config"
	"github.com/pavlo67/constructor/starter/joiner"
)

//type Info struct {
//	Service string
//	Path    string
//	Status  string
//	Details interface{}
//}

type Operator interface {

	// Name returns started component name
	Name() string

	Init(conf *config.Config, options basis.Info) (info []basis.Info, err error)

	// Setup sets up the component
	Setup() error

	// Run inits the component to use in application
	Run(joiner.Operator) error
}

type Starter struct {
	Operator
	Options basis.Info
}

func (starter Starter) CorrectedOptions(options basis.Info) basis.Info {
	newOptions := basis.Info{}

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
