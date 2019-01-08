package crud_serverhttp_jschmhr0

import (
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/pavlo67/partes/crud"
	"github.com/pavlo67/partes/fronthttp/componenthtml"
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/strlib"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"
)

func New(joiner joiner.Operator, endpoints map[string]config.Endpoint) *crud_serverhttp_jschmhr {
	if joiner == nil {
		l.Error("null joiner for crud_serverhttp_jschmhr.New()")
		return nil
	}

	return &crud_serverhttp_jschmhr{
		joiner:    joiner,
		mutex:     &sync.Mutex{},
		crudTypes: nil,
		crudOps:   map[string]crud.Operator{},
		endpoints: endpoints,
	}
}

var _ componenthtml.Operator = &crud_serverhttp_jschmhr{}

type crud_serverhttp_jschmhr struct {
	joiner     joiner.Operator
	mutex      *sync.Mutex
	crudTypes  []string
	crudOps    map[string]crud.Operator
	endpoints  map[string]config.Endpoint
	menuItems  []componenthtml.MenuItem
	staticHTML map[string]string
}

func (rcOp *crud_serverhttp_jschmhr) Name() string {
	return "CRUD"
}

func (rcOp *crud_serverhttp_jschmhr) Endpoints() map[string]config.Endpoint {
	if rcOp == nil {
		return nil
	}

	return rcOp.endpoints
}

func (rcOp *crud_serverhttp_jschmhr) Listeners() map[string]config.Listener {
	return nil
}

func (rcOp *crud_serverhttp_jschmhr) Menu(key string) []componenthtml.MenuItem {
	if rcOp == nil {
		return nil
	}

	return rcOp.menuItems
}

func (rcOp *crud_serverhttp_jschmhr) StaticHTML() map[string]string {
	if rcOp == nil {
		return nil
	}

	return rcOp.staticHTML
}

const oninitOps = "on crud_serverhttp_jschmhr.initOps()"

func (rcOp *crud_serverhttp_jschmhr) initOps() { // useOnly []joiner.InterfaceKey
	if rcOp == nil {
		return
	}

	ops := rcOp.joiner.ComponentsAll(nil, "")

	for _, crudInt := range ops {
		// l.Info("???", crudInt)

		crudOp, ok := crudInt.Interface.(crud.Operator)
		if !ok {
			// TODO: wtf?
			continue
		}

		// if len(useOnly) > 0 &&

		crudType := strings.TrimSpace(string(crudInt.Key))
		if crudType == "" {
			l.Errorf(oninitOps+": no crudType defined for %#v", crudInt)
		}

		rcOp.mutex.Lock()

		if !strlib.In(rcOp.crudTypes, crudType) {
			rcOp.crudTypes = append(rcOp.crudTypes, crudType)
		}

		label := crudType
		description, _ := crudOp.Describe()
		if description.Title != "" {
			label = description.Title
		}

		var url string
		if ep, ok := rcOp.endpoints["read_list"]; ok {
			url = ep.Path(crudType)
		}

		menuItem := componenthtml.MenuItem{
			Key:   crudType,
			Label: label,
			URL:   url,
		}

		iam := false
		for i, menuItem := range rcOp.menuItems {
			if menuItem.Key == crudType {
				rcOp.menuItems[i] = menuItem
				iam = true
				break
			}
		}

		if !iam {
			rcOp.menuItems = append(
				rcOp.menuItems,
				menuItem,
			)
		}

		rcOp.crudOps[crudType] = crudOp
		rcOp.mutex.Unlock()
	}
}

const ongetOp = "on crud_serverhttp_jschmhr.getOp()"

func (rcOp *crud_serverhttp_jschmhr) getOp(crudType string) (crud.Operator, error) {
	if rcOp == nil {
		return nil, errors.Wrap(basis.ErrNull, ongetOp+": no restcrud.Operator")
	}

	crudType = strings.TrimSpace(crudType)
	if crudType == "" {
		return nil, errors.Wrap(basis.ErrNull, ongetOp+": no crudType defined")
	}

	rcOp.mutex.Lock()
	defer rcOp.mutex.Unlock()

	crudOp := rcOp.crudOps[crudType]
	if crudOp == nil {
		return nil, errors.Wrapf(basis.ErrNull, ongetOp+": no crud.Operator for crudType %s defined yet", crudType)
	}

	return crudOp, nil
}
