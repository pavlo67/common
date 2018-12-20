package program

import (
	"reflect"
	"sync"

	"log"

	"github.com/pavlo67/punctum/basis"
	"github.com/pkg/errors"
)

type InterfaceKey string

type Interface struct {
	Interface interface{}
	Key       InterfaceKey
}

type Joiner interface {
	JoinInterface(intrfc interface{}, key InterfaceKey) error
	Interface(key InterfaceKey) interface{}
	InterfacesAll(ptrToInterface interface{}, key InterfaceKey) []Interface
	CloseAll()
}

var _ Joiner = &joiner{}

func NewJoiner() Joiner {
	return &joiner{
		appInterfaces: []Interface{},
		mutex:         &sync.Mutex{},
	}
}

type joiner struct {
	appInterfaces []Interface
	mutex         *sync.Mutex
}

var ErrJoiningNil = errors.New("can't join nil interface")

func (j *joiner) JoinInterface(intrfc interface{}, key InterfaceKey) error {
	if j == nil {
		return errors.Wrap(basis.ErrNullItem, "on .JoinInterface()")
	}
	if intrfc == nil {
		return ErrJoiningNil
	}

	j.mutex.Lock()
	j.appInterfaces = append(j.appInterfaces, Interface{intrfc, key})
	j.mutex.Unlock()

	return nil
}

func (j *joiner) Interface(key InterfaceKey) interface{} {
	if j == nil {
		log.Printf("on Joiner.Interface(%s): null Joiner item", key)
	}
	if key == "" {
		return nil
	}

	j.mutex.Lock()
	defer j.mutex.Unlock()
	for _, app := range j.appInterfaces {
		if app.Key == key {
			return app.Interface
		}
	}

	return nil
}

func (j *joiner) InterfacesAll(ptrToInterface interface{}, key InterfaceKey) []Interface {
	if j == nil {
		log.Printf("on Joiner.InterfacesAll(%T, %s): null Joiner item", ptrToInterface, key)
	}

	apps := []Interface{}

	j.mutex.Lock()
	defer j.mutex.Unlock()

	for _, app := range j.appInterfaces {

		if key != "" && app.Key != key {
			continue
		}

		if ptrToInterface == nil || checkInterface(app, ptrToInterface) {
			apps = append(apps, app)
		}

	}

	return apps
}

func checkInterface(app Interface, ptrToInterface interface{}) bool {
	defer func() {
		recover()
	}()

	// ??? reflect.TypeOf(ptrToInterface).Elem()
	if app.Interface != nil && reflect.TypeOf(app.Interface).Implements(reflect.TypeOf(ptrToInterface).Elem()) {
		return true
	}

	return false
}

func (j *joiner) CloseAll() {
	if j == nil {
		log.Print("on Joiner.Close(): null Joiner item")
		return
	}
	closersInt := j.InterfacesAll((*basis.Closer)(nil), "")
	for _, cl := range closersInt {
		if cl.Interface != nil {
			if closer, ok := cl.Interface.(basis.Closer); ok {
				err := closer.Close()
				if err != nil {
					log.Print("on Joiner.Close(): ", err)
				}
			}
		}
	}
}

// DEPRECATED
// func GetInterfaceBySignature(ptrToInterface interface{}, key InterfaceKey) interface{} {
//	defer func() {
//		recover()
//	}()
//
//	mutex.Lock()
//	defer mutex.Unlock()
//
//	for _, app := range appInterfaces {
//		if key != "" && app.Key != key {
//			continue
//		}
//
//		if app.Interface != nil && reflect.TypeOf(app.Interface).Implements(reflect.TypeOf(ptrToInterface).Elem()) {
//			return app.Interface
//		}
//	}
//
//	return nil
// }
