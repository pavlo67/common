package joiner

import (
	"log"
	"reflect"
	"sync"

	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
)

type InterfaceKey string

type Component struct {
	Interface interface{}
	Key       InterfaceKey
}

type Operator interface {
	JoinInterface(interface{}, InterfaceKey) error
	Interface(InterfaceKey) interface{}
	ComponentsAll(InterfaceKey) []Component
	ComponentsAllWithSignature(ptrToInterface interface{}) []Component
	CloseAll()
}

var _ Operator = &joiner{}

func New() Operator {
	return &joiner{
		components: []Component{},
		mutex:      &sync.Mutex{},
	}
}

type joiner struct {
	components []Component
	mutex      *sync.Mutex
}

var ErrJoiningNil = errors.New("can't join nil interface")

func (j *joiner) JoinInterface(intrfc interface{}, key InterfaceKey) error {
	if j == nil {
		return errors.Wrap(basis.ErrNull, "on .JoinInterface()")
	}
	if intrfc == nil {
		return ErrJoiningNil
	}

	j.mutex.Lock()
	j.components = append(j.components, Component{intrfc, key})
	j.mutex.Unlock()

	return nil
}

func (j *joiner) Interface(key InterfaceKey) interface{} {
	if j == nil {
		log.Printf("on Operator.Component(%s): null Operator item", key)
	}
	if key == "" {
		return nil
	}

	j.mutex.Lock()
	defer j.mutex.Unlock()
	for _, comp := range j.components {
		if comp.Key == key {
			return comp.Interface
		}
	}

	return nil
}

func (j *joiner) ComponentsAll(key InterfaceKey) []Component {
	if j == nil {
		log.Printf("on Operator.ComponentsAll(%s): null Operator item", key)
	}
	j.mutex.Lock()
	defer j.mutex.Unlock()

	var components []Component

	for _, component := range j.components {
		//if ptrToInterface == nil || !CheckInterface(component, ptrToInterface) {
		//	continue
		//}

		if component.Key == key {
			components = append(components, component)
		}
	}

	return components
}

func (j *joiner) ComponentsAllWithSignature(ptrToInterface interface{}) []Component {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	var components []Component

	for _, component := range j.components {
		//if key != "" && component.Key != key {
		//	continue
		//}

		if ptrToInterface != nil && CheckInterface(component, ptrToInterface) {
			components = append(components, component)
		}
	}

	return components
}

func CheckInterface(component Component, ptrToInterface interface{}) bool {
	defer func() {
		recover()
	}()

	// ??? reflect.TypeOf(ptrToInterface).Elem()
	if component.Interface != nil && reflect.TypeOf(component.Interface).Implements(reflect.TypeOf(ptrToInterface).Elem()) {
		return true
	}

	return false
}

func (j *joiner) CloseAll() {
	if j == nil {
		log.Print("on Operator.Close(): null Operator item")
		return
	}

	closerComponents := j.ComponentsAllWithSignature((*Closer)(nil))

	for _, closerComponent := range closerComponents {
		if closer, _ := closerComponent.Interface.(Closer); closer != nil {
			err := closer.Close()
			if err != nil {
				log.Print("on Operator.Close(): ", err)
			}
		}
	}

}
