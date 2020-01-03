package router

import (
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/identity"
	"github.com/pavlo67/workshop/common/joiner"
)

const InterfaceKey joiner.InterfaceKey = "router"

// const CleanerInterfaceKey joiner.HandlerKey = "route_cleaner"

type Routes map[identity.Domain]config.Access

type Operator interface {
	Routes() (Routes, error)
}
