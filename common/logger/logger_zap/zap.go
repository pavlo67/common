package logger_zap

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/strlib"
)

func New(cfg logger.Config) (logger.Operator, error) {
	c := zap.NewProductionConfig()
	c.DisableStacktrace = true
	c.Level.SetLevel(zapLevel(cfg.LogLevel))

	// TODO??? check if paths are correct

	if len(cfg.OutputPaths) < 1 {
		cfg.OutputPaths = []string{"stdout"}
	}
	c.OutputPaths = cfg.OutputPaths

	if len(cfg.ErrorOutputPaths) < 1 {
		cfg.ErrorOutputPaths = []string{"stderr"}
	}
	c.ErrorOutputPaths = cfg.ErrorOutputPaths

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

	commentPaths := cfg.OutputPaths
	for _, eop := range cfg.ErrorOutputPaths {
		if !strlib.In(commentPaths, eop) {
			commentPaths = append(commentPaths, eop)
		}
	}

	return &loggerZap{sugaredLogger: l.Sugar(), commentPaths: commentPaths}, nil
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

//var _ logger.Operator = &zap.SugaredLogger{}

var _ logger.Operator = &loggerZap{}

type loggerZap struct {
	sugaredLogger *zap.SugaredLogger
	commentPaths  []string
}

func (l loggerZap) Debug(args ...interface{}) {
	l.sugaredLogger.Debug(args...)
}

func (l loggerZap) Debugf(template string, args ...interface{}) {
	l.sugaredLogger.Debugf(template, args...)
}

func (l loggerZap) Info(args ...interface{}) {
	l.sugaredLogger.Info(args...)
}

func (l loggerZap) Infof(template string, args ...interface{}) {
	l.sugaredLogger.Infof(template, args...)
}

func (l loggerZap) Warn(args ...interface{}) {
	l.sugaredLogger.Warn(args...)
}

func (l loggerZap) Warnf(template string, args ...interface{}) {
	l.sugaredLogger.Warnf(template, args...)
}

func (l loggerZap) Error(args ...interface{}) {
	l.sugaredLogger.Error(args...)
}

func (l loggerZap) Errorf(template string, args ...interface{}) {
	l.sugaredLogger.Errorf(template, args...)
}

//Panic(args ...interface{}) {}
//Panicf(template string, args ...interface{}) {}

func (l loggerZap) Fatal(args ...interface{}) {
	l.sugaredLogger.Fatal(args...)
}

func (l loggerZap) Fatalf(template string, args ...interface{}) {
	l.sugaredLogger.Fatalf(template, args...)
}

func (l loggerZap) Comment(text string) {
	outstring := "\n\t\t" + text + "\n\n"
	for _, outPath := range l.commentPaths {
		switch outPath {
		case "stdout":
			fmt.Print(outstring)
		case "stderr":
			// to prevent duplicates in console
			// fmt.Fprint(os.Stderr, outPath+" "+outstring)
		default:
			f, err := os.OpenFile(outPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Fprint(os.Stderr, err)
			}
			defer f.Close()
			if _, err := f.WriteString(outstring); err != nil {
				fmt.Fprint(os.Stderr, err)
			}
		}
	}
}
