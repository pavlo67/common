package rest_flower_serverhttp_jsschmhr

import (
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/flow/datastore"
)

func New() *rest_datastore_serverhttp_jschmhr {
	return &rest_datastore_serverhttp_jschmhr{
		mutex:          &sync.Mutex{},
		datastoreTypes: nil,
		datastoreOps:   map[string]datastore.Operator{},
	}
}

type rest_datastore_serverhttp_jschmhr struct {
	mutex          *sync.Mutex
	datastoreTypes []string
	datastoreOps   map[string]datastore.Operator
}

const onnewOp = "on rest_datastore_serverhttp_jschmhr.initOps()"

func (rcOp *rest_datastore_serverhttp_jschmhr) newOp() error {
	return basis.ErrNotImplemented
	//
	//if rcOp == nil {
	//	return
	//}
	//
	//ops := rcOp.joiner.ComponentsAll(nil, "")
	//
	//for _, crudInt := range ops {
	//	// l.Info("???", crudInt)
	//
	//	crudOp, ok := crudInt.Interface.(crud.Operator)
	//	if !ok {
	//		// TODO: wtf?
	//		continue
	//	}
	//
	//	// l.Info("+++", crudOp)
	//
	//	crudType := strings.TrimSpace(string(crudInt.Key))
	//	if crudType == "" {
	//		l.Errorf(onnewOp+": no crudType defined for %#v", crudInt)
	//	}
	//
	//	rcOp.mutex.Lock()
	//
	//	if !str_json.In(rcOp.crudTypes, crudType) {
	//		rcOp.crudTypes = append(rcOp.crudTypes, crudType)
	//	}
	//
	//	label := crudType
	//	description, _ := crudOp.Describe()
	//	if description.Title != "" {
	//		label = description.Title
	//	}
	//
	//
	//
	//	rcOp.crudOps[crudType] = crudOp
	//	rcOp.mutex.Unlock()
	//}
}

const ongetOp = "on rest_datastore_serverhttp_jschmhr.getOp()"

func (rcOp *rest_datastore_serverhttp_jschmhr) getOp(datastoreType string) (datastore.Operator, error) {
	if rcOp == nil {
		return nil, errors.Wrap(basis.ErrNull, ongetOp+": no rest_datastore_serverhttp_jschmhr.Operator")
	}

	datastoreType = strings.TrimSpace(datastoreType)
	if datastoreType == "" {
		return nil, errors.Wrap(basis.ErrNull, ongetOp+": no datastoreType defined")
	}

	rcOp.mutex.Lock()
	defer rcOp.mutex.Unlock()

	datastoreOp := rcOp.datastoreOps[datastoreType]
	if datastoreOp == nil {
		return nil, errors.Wrapf(basis.ErrNull, ongetOp+": no datastore.Operator for datastoreType %s defined yet", datastoreType)
	}

	return datastoreOp, nil
}
