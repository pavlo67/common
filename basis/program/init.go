package program

import (
	"github.com/pavlo67/punctum/basis/config"
	"github.com/pavlo67/punctum/basis/logger"
	"go.uber.org/zap"
)

func Init(cfgPath string, getIdentity bool) (*config.PunctumConfig, error) {
	conf, err := config.Get(cfgPath)
	if err != nil {
		return conf, err
	}

	err = logger.Init(logger.Config{
		LogLevel: zap.DebugLevel,
	})

	return conf, err
}
