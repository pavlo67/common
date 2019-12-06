package joiner

import (
	"log"
	"reflect"
	"sync"

	"github.com/pkg/errors"
)

const InterfaceKeyFld = "interface_key"

type InterfaceKey string

type Component struct {
	InterfaceKey
	Interface interface{}
}

type Operator interface {
	Join(interface{}, InterfaceKey) error
	Interface(InterfaceKey) interface{}
	InterfacesAll(ptrToInterface interface{}) []Component
	CloseAll()
}

var _ Operator = &joiner{}

func New() Operator {
	return &joiner{
		components: map[InterfaceKey]interface{}{},
		mutex:      &sync.Mutex{},
	}
}

type joiner struct {
	components map[InterfaceKey]interface{}
	mutex      *sync.Mutex
}

var ErrJoiningNil = errors.New("can't join nil interface")
var ErrJoiningDuplicate = errors.New("can't join interface over joined before")

func (j *joiner) Join(intrfc interface{}, interfaceKey InterfaceKey) error {
	if j == nil {
		return errors.Errorf("got nil on .Join(%s)", interfaceKey)
	}
	if intrfc == nil {
		return errors.Wrapf(ErrJoiningNil, "on .Join(%s)", interfaceKey)
	}

	j.mutex.Lock()
	defer j.mutex.Unlock()

	if _, ok := j.components[interfaceKey]; ok {
		return errors.Wrapf(ErrJoiningDuplicate, "on .Join(%s)", interfaceKey)
	}

	j.components[interfaceKey] = intrfc

	return nil
}

func (j *joiner) Interface(interfaceKey InterfaceKey) interface{} {
	if j == nil {
		log.Printf("on Operator.Component(%s): null Operator item", interfaceKey)
	}

	j.mutex.Lock()
	defer j.mutex.Unlock()

	if intrfc, ok := j.components[interfaceKey]; ok {
		return intrfc
	}

	return nil
}

func (j *joiner) InterfacesAll(ptrToInterface interface{}) []Component {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	var components []Component

	for key, intrfc := range j.components {
		if CheckInterface(intrfc, ptrToInterface) {
			components = append(components, Component{InterfaceKey: key, Interface: intrfc})
		}
	}

	return components
}

func CheckInterface(intrfc interface{}, ptrToInterface interface{}) bool {
	defer func() {
		recover()
	}()

	// ??? reflect.TypeOf(ptrToInterface).Elem()
	// ??? if intrfc != nil
	if reflect.TypeOf(intrfc).Implements(reflect.TypeOf(ptrToInterface).Elem()) {
		return true
	}

	return false
}

func (j *joiner) CloseAll() {
	if j == nil {
		log.Print("on Operator.Close(): null Operator item")
		return
	}

	closerComponents := j.InterfacesAll((*Closer)(nil))

	for _, closerComponent := range closerComponents {
		if closer, _ := closerComponent.Interface.(Closer); closer != nil {
			err := closer.Close()
			if err != nil {
				log.Print("on Operator.Close(): ", err)
			}
		}
	}

}
