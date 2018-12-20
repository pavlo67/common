package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	LogLevel    zapcore.Level
	OutputPaths []string
	Encoding    string
}

var mutex = &sync.Mutex{}
var logger *zap.SugaredLogger

func Init(cfg Config) error {
	c := zap.NewProductionConfig()
	c.DisableStacktrace = true
	c.Level.SetLevel(cfg.LogLevel)
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

	mutex.Lock()
	logger = l.Sugar()
	mutex.Unlock()

	return nil
}

func Get() *zap.SugaredLogger {
	mutex.Lock()
	l := logger
	mutex.Unlock()

	if l == nil {
		panic("no logger (zap.SugaredLogger) to use found!!!")
	}

	return l
}
