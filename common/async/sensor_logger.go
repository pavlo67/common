package async

import (
	"fmt"
	"time"

	"github.com/pavlo67/common/common/errors"
)

// LoggerSensor ------------------------------------------------------------------------------------------

var _ Sensor = NewLoggerSensor[int](1, 1, nil)

func NewLoggerSensor[S any](c0, c1 Code, fmt func(time.Time, *S)) *LoggerSensor[S] {
	return &LoggerSensor[S]{c0: c0, c1: c1, fmt: fmt}
}

type LoggerSensor[S any] struct {
	c0  Code
	c1  Code
	fmt func(time.Time, *S)
}

func (ls *LoggerSensor[S]) Code() (_, _ Code) {
	if ls == nil {
		return 0, 0
	}

	return ls.c0, ls.c1
}

const onLoggerSet = "on LoggerSensor.Set()"

func (ls *LoggerSensor[S]) Set(t time.Time, values interface{}) error {
	if ls == nil {
		return errors.New("ls == nil / " + onLoggerSet)
	}

	var sensorValues *S

	switch v := values.(type) {
	case S:
		sensorValues = &v
	case *S:
		sensorValues = v
	default:
		return fmt.Errorf("wrong values: %#v / "+onLoggerSet, values)
	}

	ls.fmt(t, sensorValues)

	return nil
}

func (ls *LoggerSensor[S]) Check() *time.Time {
	return nil
}
