package starter

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
)

type Operator interface {
	Name() string
	Run(*config.Config, common.Map, joiner.Operator, logger.Operator) error
}

type Component struct {
	Operator
	Options common.Map
}
