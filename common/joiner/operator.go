package joiner

import (
	"fmt"
	"log"
	"reflect"
	"sync"

	"github.com/pavlo67/common/common/logger"

	"github.com/pavlo67/common/common"
	"github.com/pkg/errors"
)

type Component struct {
	common.InterfaceKey
	Interface interface{}
}

type ID common.IDStr

type Operator interface {
	Join(interface{}, common.InterfaceKey) error
	Interface(common.InterfaceKey) interface{}
	InterfacesAll(ptrToInterface interface{}) []Component
	CloseAll()
}

var _ Operator = &joiner{}

func New(options common.Map, l logger.Operator) Operator {
	return &joiner{
		l:          l,
		options:    options,
		components: map[common.InterfaceKey]interface{}{},
		mutex:      &sync.RWMutex{},
	}
}

type joiner struct {
	l          logger.Operator
	options    common.Map
	components map[common.InterfaceKey]interface{}
	mutex      *sync.RWMutex
}

var ErrJoiningNil = errors.New("can't join nil interface")
var ErrJoiningDuplicate = errors.New("can't join interface over joined before")

func (j *joiner) Join(intrfc interface{}, interfaceKey common.InterfaceKey) error {
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

	if j.l != nil && !j.options.IsTrue("silent") {
		j.l.Infof("joined %s", interfaceKey)
	}

	return nil
}

func (j *joiner) Interface(interfaceKey common.InterfaceKey) interface{} {
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
