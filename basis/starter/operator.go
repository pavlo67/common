package starter

import (
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/config"
	"github.com/pavlo67/punctum/basis/program"
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

	Prepare(conf *config.PunctumConfig, params basis.Params) error

	// Check checks status of the component that implements this interface
	Check() (info []Info, err error)

	// Setup sets up the component
	Setup() error

	// Init inits the component to use in application
	Init(program.Joiner) error
}

type Starter struct {
	Operator
	Params basis.Params
}

//type Runner interface {
//	Run() error
//}
