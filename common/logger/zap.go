package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var mutexZap = &sync.Mutex{}
var loggerZap *zap.SugaredLogger

var _ Operator = &zap.SugaredLogger{}

func zapInit(cfg Config) error {
	c := zap.NewProductionConfig()
	c.DisableStacktrace = true
	c.Level.SetLevel(zapLevel(cfg.LogLevel))
	c.OutputPaths = cfg.OutputPaths
	if len(c.OutputPaths) < 1 {
		c.OutputPaths = []string{"stdout"}
	}

	c.Encoding = cfg.Encoding
	if c.Encoding == "" {
		c.Encoding = "console"
	}
	c.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	c.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	l, err := c.Build()
	if err != nil {
		return err
	}

	mutexZap.Lock()
	loggerZap = l.Sugar()
	mutexZap.Unlock()

	return nil
}

func zapLevel(level Level) zapcore.Level {
	switch level {
	case TraceLevel:
		return zapcore.DebugLevel
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case PanicLevel:
		return zapcore.PanicLevel
	case FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func zapGet() *zap.SugaredLogger {
	mutexZap.Lock()
	l := loggerZap
	mutexZap.Unlock()

	if l == nil {
		panic("no loggerZap (zap.SugaredLogger) to use found!!!")
	}

	return l
}
