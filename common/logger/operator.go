package logger

import (
	"errors"
)

type Level int32

type Config struct {
	LogLevel    Level
	OutputPaths []string
	Encoding    string
}

const TraceLevel Level = -2
const DebugLevel Level = -1
const InfoLevel Level = 0
const WarnLevel Level = 1
const ErrorLevel Level = 2
const PanicLevel Level = 3
const FatalLevel Level = 4

type Operator interface {
	//Trace(args ...interface{})
	//Tracef(template string, args ...interface{})

	Debug(args ...interface{})
	Debugf(template string, args ...interface{})

	Info(args ...interface{})
	Infof(template string, args ...interface{})

	Warn(args ...interface{})
	Warnf(template string, args ...interface{})

	Error(args ...interface{})
	Errorf(template string, args ...interface{})

	Panic(args ...interface{})
	Panicf(template string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
}

func Init(loggerCfg Config) (Operator, error) {
	err := zapInit(loggerCfg)
	if err != nil {
		return nil, err
	}

	l := Operator(zapGet())
	if l == nil {
		return nil, errors.New("no logger ???")
	}

	return l, nil
}

func Get() Operator {
	return Operator(zapGet())
}
