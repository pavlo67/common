package router

import (
	"github.com/pavlo67/workshop/common/joiner"
)

const InterfaceKey joiner.InterfaceKey = "router"

// const CleanerInterfaceKey joiner.InterfaceKey = "route_cleaner"

type Domain string
type URL string

type Routes map[Domain]URL

type Operator interface {
	Routes() (Routes, error)
}
