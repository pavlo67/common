package joiner

import (
	"fmt"
	"log"
	"reflect"
	"sync"

	"github.com/pavlo67/common/common"

	"github.com/pavlo67/common/common/errors"
)

type InterfaceKey string

type Component struct {
	InterfaceKey
	Interface interface{}
}

type ID common.IDStr

type Operator interface {
	Join(interface{}, InterfaceKey) error
	Interface(InterfaceKey) interface{}
	InterfacesAll(ptrToInterface interface{}) []Component
	CloseAll()
}

type Link struct {
	InterfaceKey InterfaceKey `bson:",omitempty" json:",omitempty"`
	ID           ID           `bson:",omitempty" json:",omitempty"`
}

var _ Operator = &joiner{}

func New() Operator {
	return &joiner{
		components: map[InterfaceKey]interface{}{},
		mutex:      &sync.RWMutex{},
	}
}

type joiner struct {
	components map[InterfaceKey]interface{}
	mutex      *sync.RWMutex
}

var ErrJoiningNil = errors.New("can't join nil interface")
var ErrJoiningDuplicate = errors.New("can't join interface over joined before")

func (j *joiner) Join(intrfc interface{}, interfaceKey InterfaceKey) error {
	if j == nil {
		return fmt.Errorf("got nil on .Join(%s)", interfaceKey)
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
		log.Printf("on ActorKey.Component(%s): null ActorKey item", interfaceKey)
	}

	j.mutex.RLock()
	defer j.mutex.RUnlock()

	if intrfc, ok := j.components[interfaceKey]; ok {
		return intrfc
	}

	return nil
}

func (j *joiner) InterfacesAll(ptrToInterface interface{}) []Component {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

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
		log.Print("on ActorKey.Close(): null ActorKey item")
		return
	}

	closerComponents := j.InterfacesAll((*Closer)(nil))

	for _, closerComponent := range closerComponents {
		if closer, _ := closerComponent.Interface.(Closer); closer != nil {
			err := closer.Close()
			if err != nil {
				log.Print("on ActorKey.Close(): ", err)
			}
		}
	}

}
