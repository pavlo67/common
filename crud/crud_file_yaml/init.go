package crud_file_yaml

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"
	"github.com/pavlo67/punctum/starter/logger"
)

var l logger.Operator
var _ starter.Operator = &crud_file_yamlStarter{}

func Starter() starter.Operator {
	return &crud_file_yamlStarter{}
}

type crud_file_yamlStarter struct {
	interfaceKey        joiner.InterfaceKey
	cleanerInterfaceKey joiner.InterfaceKey
	path                string
}

func (nms *crud_file_yamlStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (nms *crud_file_yamlStarter) Prepare(conf *config.PunctumConfig, options, runtimeOptions basis.Options) error {
	l = logger.Get()

	var ok bool
	nms.path, ok = options.String("path")
	if !ok {
		return errors.New("no path for crud_file_yamlStarter.Prepare()")
	}

	//nms.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(crud.InterfaceKey)))
	//nms.cleanerInterfaceKey = joiner.InterfaceKey(options.StringDefault("cleaner_interface_key", string(crud.CleanerInterfaceKey)))

	return nil
}

func (nms *crud_file_yamlStarter) Check() (info []starter.Info, err error) {
	return nil, nil
}

func (nms *crud_file_yamlStarter) Setup() error {
	return nil
}

func (nms *crud_file_yamlStarter) Init(joiner joiner.Operator) error {

	//var err error
	//notesOp, err := New(
	//	nms.mysqlConfig,
	//	nms.index.Params["table"],
	//	nms.jointLinks,
	//	nil, // grpsOp,
	//	nil, // linksOp,
	//	nil,
	//)
	//if err != nil {
	//	return errors.Wrap(err, "can't init notes_mysql")
	//}
	//
	//err = joiner.JoinInterface(notesOp, nms.interfaceKey)
	//if err != nil {
	//	return errors.Wrap(err, "can't join notes_mysql as notes.Operator interface")
	//}

	//err = joiner.JoinInterface(dataOp.Clean, ds.cleanerInterfaceKey)
	//if err != nil {
	//	return errors.Wrapf(err, "can't join datastoremysql.Operator.Clean as %s", ds.cleanerInterfaceKey)
	//}

	return nil
}
