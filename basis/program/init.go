package program

import (
	"github.com/pavlo67/punctum/basis/config"
	"github.com/pavlo67/punctum/basis/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init(cfgPath string, level zapcore.Level) (*config.PunctumConfig, *zap.SugaredLogger, error) {
	conf, err := config.Get(cfgPath)
	if err != nil {
		return conf, nil, err
	}

	err = logger.Init(logger.Config{LogLevel: level})

	return conf, logger.Get(), err
}
