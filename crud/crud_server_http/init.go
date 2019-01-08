package crud_serverhttp_jschmhr0

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/filelib"
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"
	"github.com/pavlo67/punctum/starter/logger"

	"github.com/pavlo67/partes/fronthttp/componenthtml"
)

const InterfaceKey joiner.InterfaceKey = "crud_serverhttp_jschmhr"

var l *zap.SugaredLogger

func Starter() starter.Operator {
	return &crud_serverhttp_jschmhrStarter{}
}

type crud_serverhttp_jschmhrStarter struct {
	index        config.ServerComponentsIndex
	interfaceKey joiner.InterfaceKey
}

var _ starter.Operator = &crud_serverhttp_jschmhrStarter{}

func (fcs *crud_serverhttp_jschmhrStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (fcs *crud_serverhttp_jschmhrStarter) Prepare(conf *config.PunctumConfig, params basis.Params) error {
	l = logger.zapGet()

	fcs.interfaceKey = joiner.InterfaceKey(params.StringKeyDefault("interface_key", string(InterfaceKey)))

	var errs basis.Errors
	indexPath := params.StringKeyDefault("index_path", filelib.CurrentPath())
	fcs.index, errs = config.ComponentIndex(indexPath, errs)

	return errs.Err()
}

func (fcs *crud_serverhttp_jschmhrStarter) Check() (info []joiner.Info, err error) {
	return nil, nil
}

func (fcs *crud_serverhttp_jschmhrStarter) Setup() error {
	return nil
}

func (fcs *crud_serverhttp_jschmhrStarter) Init(joiner joiner.Operator) error {
	serverOp, ok := joiner.Interface(serverhttp_jschmhr.InterfaceKey).(serverhttp_jschmhr.Operator)
	if !ok {
		return errors.New("no serverhttp_jschmhr.Operator interface found for frontcrud_serverhttp_jschmhrStarterStarter.zapInit()")
	}

	fcOp := New(joiner, fcs.index.Endpoints)

	// TODO: repeat after some new crudOps are added
	fcOp.initOps()

	htmlHandlers := map[string]serverhttp_jschmhr.HTMLHandler{
		"read":      fcOp.Read,
		"read_list": fcOp.ReadList,
	}

	restHandlers := map[string]serverhttp_jschmhr.RESTHandler{
		"save_rest":        fcOp.SaveREST,
		"update_list_rest": fcOp.UpdateListREST,
		"read_rest":        fcOp.ReadREST,
		"read_list_rest":   fcOp.ReadListREST,
		"delete_rest":      fcOp.DeleteREST,
	}

	errs := serverhttp_jschmhr.InitEndpoints(
		serverOp,
		fcs.index.Endpoints,
		htmlHandlers,
		restHandlers,
		nil,
		nil,
	)

	err := joiner.JoinInterface(fcOp, fcs.interfaceKey)
	errs = errs.Append(err)

	htmlFront, err := serverhttp.ManageJSFile(
		serverOp,
		filelib.CurrentPath()+"front.js",
		"/js_crud.js",
		fcs.interfaceKey,
		nil,
		fcs.index.Endpoints,
		fcs.index.Listeners,
	)
	if err != nil {
		return errors.Wrapf(err, "can't serverhttp.ManageJSFile(%s)", filelib.CurrentPath()+"front.js")
	}

	op := componenthtml.New(
		"auth",
		fcs.index.Endpoints,
		fcs.index.Listeners,
		map[string]string{"front": htmlFront},
		nil,
	)

	err = joiner.JoinInterface(op, fcs.interfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't join serverhttp_jschmhr as componenthtml.Operator interface")
	}

	return errs.Append(err).Err()

}
