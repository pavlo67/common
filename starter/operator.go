package starter

import (
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"
)

type Info struct {
	Service string
	Path    string
	Status  string
	Details interface{}
}

type Operator interface {

	// Name returns started component name
	Name() string

	Prepare(conf *config.Config, params, runtimeOptions basis.Options) error

	// Check checks status of the component that implements this interface
	Check() (info []Info, err error)

	// Setup sets up the component
	Setup() error

	// Init inits the component to use in application
	Init(joiner.Operator) error
}

type Starter struct {
	Operator
	Options basis.Options
}

//type Runner interface {
//	Run() error
//}
