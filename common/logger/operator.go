package logger

import (
	"github.com/pavlo67/common/common/joiner"
)

const InterfaceKey joiner.InterfaceKey = "logger"

type Level int

type Config struct {
	LogLevel         Level
	OutputPaths      []string
	ErrorOutputPaths []string
	Encoding         string
}

const TraceLevel Level = -2
const DebugLevel Level = -1
const InfoLevel Level = 0
const WarnLevel Level = 1
const ErrorLevel Level = 2

// const PanicLevel Level = 3
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

	//Panic(args ...interface{})
	//Panicf(template string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
}

type OperatorComments interface {
	Comment(text string)
}
