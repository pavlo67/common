package logger_zap

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/pavlo67/common/common/logger"
)

var _ logger.Operator = &zap.SugaredLogger{}

func New(cfg logger.Config) (logger.Operator, error) {
	c := zap.NewProductionConfig()
	c.DisableStacktrace = true
	c.Level.SetLevel(zapLevel(cfg.LogLevel))

	if c.OutputPaths = cfg.OutputPaths; len(c.OutputPaths) < 1 {
		c.OutputPaths = []string{"stdout"}
	}
	if c.ErrorOutputPaths = cfg.ErrorOutputPaths; len(c.ErrorOutputPaths) < 1 {
		c.ErrorOutputPaths = []string{"stderr"}
	}

	c.Encoding = cfg.Encoding
	if c.Encoding == "" {
		c.Encoding = "console"
	}
	c.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	c.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	l, err := c.Build()
	if err != nil {
		return nil, fmt.Errorf("can't create logger (%#v --> %#v), got %s", cfg, c, err)
	}

	return l.Sugar(), nil
}

func zapLevel(level logger.Level) zapcore.Level {
	switch level {
	case logger.TraceLevel:
		return zapcore.DebugLevel
	case logger.DebugLevel:
		return zapcore.DebugLevel
	case logger.InfoLevel:
		return zapcore.InfoLevel
	case logger.WarnLevel:
		return zapcore.WarnLevel
	case logger.ErrorLevel:
		return zapcore.ErrorLevel
	//case PanicLevel:
	//	return zapcore.PanicLevel
	case logger.FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
