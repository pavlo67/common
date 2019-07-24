package starter

import (
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"
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

//type Runner interface {
//	Run() error
//}
