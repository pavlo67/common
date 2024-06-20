package logger_zap

import (
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/filelib"
	"github.com/pavlo67/common/common/imagelib"
	"github.com/pavlo67/common/common/logger"
)

type loggerZap struct {
	zap.SugaredLogger
	logger.Config
}

func (op loggerZap) Comment(text string) {
	for _, commentPath := range append(op.Config.OutputPaths, op.Config.ErrorPaths...) {
		if commentPath == "stdout" {
			fmt.Fprintf(os.Stdout, "\n%s\n", text)
		} else if commentPath == "stderr" {
			fmt.Fprintf(os.Stderr, "\n%s\n", text)
		} else {
			f, err := os.OpenFile(commentPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				op.Errorf("CAN'T OPEN PATH (%s) TO COMMENT: %s", commentPath, err)
				continue
			}
			defer f.Close()
			if _, err = f.WriteString(text); err != nil {
				op.Errorf("CAN'T WRITE COMMENT TO %s: %s", commentPath, err)
			}
		}
	}
}

func (op *loggerZap) SetPath(basePath string) {
	if op == nil {
		return
	}
	if basePath = strings.TrimSpace(basePath); basePath == "" {
		op.Config.BasePath = ""
		return
	}

	var err error
	if basePath, err = filelib.Dir(basePath); err != nil {
		op.Errorf("can't create basePath (%s): %s / on logger.SetPath()", basePath, err)
		return
	}

	op.Config.BasePath = basePath
}

func (op loggerZap) File(path string, data []byte) {
	if op.Config.SaveFiles {
		basedPaths, err := logger.ModifyPaths([]string{path}, op.Config.BasePath)
		if err != nil {
			op.Error(err)
		} else if err := os.WriteFile(basedPaths[0], data, 0644); err != nil {
			op.Errorf("CAN'T WRITE TO FILE %s: %s", path, err)
		} else {
			op.Infof("FILE IS WRITTEN  %s", basedPaths[0])
		}
	}
}

func (op loggerZap) Image(path string, getImage logger.GetImage, opts common.Map) {
	if op.Config.SaveFiles {
		img, info, err := getImage.Image(opts)
		if info != "" {
			op.File(path+".info", []byte(info))

			//_, filename, line, _ := runtime.Caller(1)
			//op.Infof("from %s:%d: "+info, filename, line)
		}
		if img != nil {
			basedPaths, err := logger.ModifyPaths([]string{path}, op.Config.BasePath)
			if err != nil {
				op.Error(err)
			} else if err = imagelib.Save(img, basedPaths[0]); err != nil {
				op.Error(err)
			}
		}
		if err != nil {
			op.Error(err)
		}
	}
}

func (op *loggerZap) SetKey(key string) {
	if op == nil {
		return
	}
	op.Config.Key = key
}

func (op loggerZap) Key() string {
	return op.Config.Key
}

var _ logger.Operator = &loggerZap{}

func New(cfg logger.Config) (logger.Operator, error) {
	c := zap.NewProductionConfig()
	c.DisableStacktrace = true
	c.Level.SetLevel(zapLevel(cfg.LogLevel))

	var err error

	if len(cfg.OutputPaths) < 1 {
		c.OutputPaths = []string{"stdout"}
	} else {
		c.OutputPaths, err = logger.ModifyPaths(cfg.OutputPaths, cfg.BasePath)
		if err != nil {
			return nil, err
		}
	}

	if len(cfg.ErrorPaths) < 1 {
		c.ErrorOutputPaths = []string{"stderr"}
	} else {
		c.ErrorOutputPaths, err = logger.ModifyPaths(cfg.ErrorPaths, cfg.BasePath)
		if err != nil {
			return nil, err
		}

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

	if cfg.Key == "" {
		cfg.Key = time.Now().Format(time.RFC3339)[:19]
	}

	return &loggerZap{SugaredLogger: *l.Sugar(), Config: cfg}, nil
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
