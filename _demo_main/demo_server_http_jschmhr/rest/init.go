package rest_flower_serverhttp_jsschmhr

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/config"
	"github.com/pavlo67/punctum/basis/filelib"
	"github.com/pavlo67/punctum/basis/logger"
	"github.com/pavlo67/punctum/basis/program"
	"github.com/pavlo67/punctum/basis/starter"

	"github.com/pavlo67/partes/serverhttp/serverhttp_jschmhr"
)

var l *zap.SugaredLogger

func Starter() starter.Operator {
	return &rest_serverhttp_jschmhrStarter{}
}

type rest_serverhttp_jschmhrStarter struct {
	index config.ServerComponentsIndex
	// interfaceKey program.InterfaceKey
}

var _ starter.Operator = &rest_serverhttp_jschmhrStarter{}

func (css *rest_serverhttp_jschmhrStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (css *rest_serverhttp_jschmhrStarter) Prepare(conf *config.PunctumConfig, params basis.Params) error {
	l = logger.Get()

	// css.interfaceKey = program.InterfaceKey(params.StringKeyDefault("interface_key", string(InterfaceKey)))

	var errs basis.Errors
	indexPath := params.StringKeyDefault("index_path", filelib.CurrentPath())
	css.index, errs = config.ComponentIndex(indexPath, errs)

	return errs.Err()
}

func (css *rest_serverhttp_jschmhrStarter) Check() (info []program.Info, err error) {
	return nil, nil
}

func (css *rest_serverhttp_jschmhrStarter) Setup() error {
	return nil
}

func (css *rest_serverhttp_jschmhrStarter) Init(joiner program.Joiner) error {
	serverOp, ok := joiner.Interface(serverhttp_jschmhr.InterfaceKey).(serverhttp_jschmhr.Operator)
	if !ok {
		return errors.New("no serverhttp_jschmhr.Operator interface found for frontrest_serverhttp_jschmhrStarterStarter.Init()")
	}

	//cssOp := New()

	restHandlers := map[string]serverhttp_jschmhr.RESTHandler{
		//"new":        cssOp.NewDatastore,
		//"save":       cssOp.Save,
		//"read_list":  cssOp.ReadList,
		//"delete":     cssOp.DeleteList,
		//"key_exists": cssOp.KeyExists,
		//"last_key":   cssOp.LastKey,
	}

	errs := serverhttp_jschmhr.InitEndpoints(
		serverOp,
		css.index.Endpoints,
		nil,
		restHandlers,
		nil,
		nil,
	)

	return errs.Err()

}
