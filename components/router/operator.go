package router

import (
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
)

const InterfaceKey joiner.InterfaceKey = "router"

// const CleanerInterfaceKey joiner.InterfaceKey = "route_cleaner"

type Domain string

type Routes map[Domain]config.Access

type Operator interface {
	Routes() (Routes, error)
}
