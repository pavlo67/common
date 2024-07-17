package sensors

import (
	"fmt"
	"sync"
	"time"

	"github.com/pavlo67/common/common/errors"
)

type Code = uint16

// implementations of Item must be thread-safe
type Item interface {
	Code() (_, _ Code)
	Set(_ time.Time, data interface{}) error
	Check() *time.Time
}

// GenericSensor ------------------------------------------------------------------------------------------

var _ Item = NewGenericSensor[int](1, 1)

func NewGenericSensor[S any](c0, c1 Code) *GenericSensor[S] {
	return &GenericSensor[S]{c0: c0, c1: c1, m: &sync.Mutex{}}
}

type GenericSensor[S any] struct {
	c0 Code
	c1 Code
	t  time.Time
	v  *S
	m  *sync.Mutex
}

func (gs *GenericSensor[S]) Code() (_, _ Code) {
	if gs == nil {
		return 0, 0
	}

	return gs.c0, gs.c1
}

const onSet = "on GenericSensor.Set()"

func (gs *GenericSensor[S]) Set(t time.Time, values interface{}) error {
	if gs == nil {
		return errors.New("gs == nil / " + onSet)
	}

	var sensorValues *S

	switch v := values.(type) {
	case S:
		sensorValues = &v
	case *S:
		sensorValues = v
	default:
		return fmt.Errorf("wrong values: %#v / "+onSet, values)
	}

	gs.m.Lock()
	defer gs.m.Unlock()
	gs.t, gs.v = t, sensorValues

	return nil
}

func (gs *GenericSensor[S]) SetTyped(t time.Time, values S) {
	if gs == nil {
		// TODO???
		return
	}

	gs.m.Lock()
	defer gs.m.Unlock()
	gs.t, gs.v = t, &values
}

func (gs *GenericSensor[S]) GetTyped() (time.Time, *S) {
	if gs == nil {
		return time.Time{}, nil
	}

	gs.m.Lock()
	defer gs.m.Unlock()

	return gs.t, gs.v
}

func (gs *GenericSensor[S]) Check() *time.Time {
	if gs == nil {
		return nil
	}

	gs.m.Lock()
	defer gs.m.Unlock()

	if gs.v == nil {
		return nil
	}
	return &gs.t
}
