package actor

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/joiner"
)

type Operator interface {
	Name() string
	Run(params common.Map) (posterior []joiner.Link, info common.Map, err error)
}
