package joiner

import (
	"log"
	"reflect"
	"sync"

	"github.com/pkg/errors"

	"github.com/pavlo67/constructor/basis"
)

type ComponentKey string

type Component struct {
	Worker interface{}
	Key    ComponentKey
}

type Operator interface {
	JoinComponent(interface{}, ComponentKey) error
	Worker(ComponentKey) interface{}
	ComponentsAllWithInterface(ptrToInterface interface{}) []Component
	CloseAll()

	// DEPRECATED
	Interface(ComponentKey) interface{}

}

var _ Operator = &joiner{}

func New() Operator {
	return &joiner{
		components: map[ComponentKey]Component{},
		mutex:      &sync.Mutex{},
	}
}

type joiner struct {
	components map[ComponentKey]Component
	mutex      *sync.Mutex
}

var ErrJoiningNil = errors.New("can't join nil interface")

func (j *joiner) JoinComponent(intrfc interface{}, key ComponentKey) error {
	if j == nil {
		return errors.Wrap(basis.ErrNull, "on .JoinComponent()")
	}
	if intrfc == nil {
		return ErrJoiningNil
	}

	j.mutex.Lock()
	j.components[key] = Component{intrfc, key}
	j.mutex.Unlock()

	return nil
}

// DEPRECATED
func (j *joiner) Interface(key ComponentKey) interface{} {
	return j.Worker(key)
}

func (j *joiner) Worker(key ComponentKey) interface{} {
	if j == nil {
		log.Printf("on Operator.Component(%s): null Operator item", key)
	}
	if key == "" {
		return nil
	}

	j.mutex.Lock()
	defer j.mutex.Unlock()

	if comp, ok := j.components[key]; ok {
		return comp.Worker
	}
	//for _, comp := range j.components {
	//	if comp.Key == key {
	//		return comp.Worker
	//	}
	//}

	return nil
}

//func (j *joiner) ComponentsAll(key ComponentKey) []Component {
//	if j == nil {
//		log.Printf("on Operator.ComponentsAll(%s): null Operator item", key)
//	}
//	j.mutex.Lock()
//	defer j.mutex.Unlock()
//
//	var components []Component
//
//	for _, component := range j.components {
//		//if ptrToInterface == nil || !CheckInterface(component, ptrToInterface) {
//		//	continue
//		//}
//
//		if component.Key == key {
//			components = append(components, component)
//		}
//	}
//
//	return components
//}

func (j *joiner) ComponentsAllWithInterface(ptrToInterface interface{}) []Component {
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
	if component.Worker != nil && reflect.TypeOf(component.Worker).Implements(reflect.TypeOf(ptrToInterface).Elem()) {
		return true
	}

	return false
}

func (j *joiner) CloseAll() {
	if j == nil {
		log.Print("on Operator.Close(): null Operator item")
		return
	}

	closerComponents := j.ComponentsAllWithInterface((*Closer)(nil))

	for _, closerComponent := range closerComponents {
		if closer, _ := closerComponent.Worker.(Closer); closer != nil {
			err := closer.Close()
			if err != nil {
				log.Print("on Operator.Close(): ", err)
			}
		}
	}

}
