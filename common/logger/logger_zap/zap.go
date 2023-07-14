package logger_zap

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/pavlo67/common/common/imagelib"
	"github.com/pavlo67/common/common/logger"
)

type loggerZap struct {
	zap.SugaredLogger
	cfg logger.Config
}

func (loggerOp *loggerZap) Comment(text string) {
	for _, commentPath := range append(loggerOp.cfg.OutputPaths, loggerOp.cfg.ErrorPaths...) {
		if commentPath == "stdout" {
			fmt.Fprintf(os.Stdout, "\n%s\n", text)
		} else if commentPath == "stderr" {
			fmt.Fprintf(os.Stderr, "\n%s\n", text)
		} else {
			f, err := os.OpenFile(commentPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				loggerOp.Errorf("CAN'T OPEN PATH (%s) TO COMMENT: %s", commentPath, err)
				continue
			}
			defer f.Close()
			if _, err = f.WriteString(text); err != nil {
				loggerOp.Errorf("CAN'T WRITE COMMENT TO %s: %s", commentPath, err)
			}
		}
	}
}

func (loggerOp *loggerZap) File(path string, data []byte) {
	if loggerOp.cfg.SaveFiles {
		basedPaths := logger.ModifyPaths([]string{path}, loggerOp.cfg.BasePath)
		if err := os.WriteFile(basedPaths[0], data, 0644); err != nil {
			loggerOp.Errorf("CAN'T WRITE TO FILE %s: %s", path, err)
		}
	}
}

func (loggerOp *loggerZap) Image(path string, getImage imagelib.GetImage) {
	if loggerOp.cfg.SaveFiles {
		img, info, err := getImage.Image()
		if info != "" {
			loggerOp.Info(info)
		}
		if img != nil {
			basedPaths := logger.ModifyPaths([]string{path}, loggerOp.cfg.BasePath)
			if err = imagelib.SavePNG(img, basedPaths[0]); err != nil {
				loggerOp.Error(err)
			}
		}
		if err != nil {
			loggerOp.Error(err)
		}
	}
}

func (loggerOp *loggerZap) NoOps() {
}

var _ logger.Operator = &loggerZap{}

func New(cfg logger.Config) (logger.Operator, error) {
	c := zap.NewProductionConfig()
	c.DisableStacktrace = true
	c.Level.SetLevel(zapLevel(cfg.LogLevel))

	if len(cfg.OutputPaths) < 1 {
		c.OutputPaths = []string{"stdout"}
	} else {
		c.OutputPaths = logger.ModifyPaths(cfg.OutputPaths, cfg.BasePath)
	}

	if len(cfg.ErrorPaths) < 1 {
		c.ErrorOutputPaths = []string{"stderr"}
	} else {
		c.ErrorOutputPaths = logger.ModifyPaths(cfg.ErrorPaths, cfg.BasePath)
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

	return &loggerZap{SugaredLogger: *l.Sugar(), cfg: cfg}, nil
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
