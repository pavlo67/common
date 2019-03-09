package crud_file

import (
	"os"

	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/crud"
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"
	"github.com/pavlo67/punctum/starter/logger"
)

var l logger.Operator
var _ starter.Operator = &crud_fileStarter{}

func Starter() starter.Operator {
	return &crud_fileStarter{}
}

type crud_fileStarter struct {
	interfaceKey joiner.InterfaceKey
	path         string
	marshaler    basis.Marshaler
	mapper       crud.Mapper
}

func (nms *crud_fileStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (nms *crud_fileStarter) Prepare(conf *config.PunctumConfig, options, runtimeOptions basis.Options) error {
	l = logger.Get()

	var ok bool
	nms.path, ok = options.String("path")
	if !ok {
		return errors.New("no path for crud_fileStarter.Prepare()")
	}

	nms.marshaler, ok = options["marshaler"].(basis.Marshaler)
	if !ok || nms.marshaler == nil {
		return errors.New("no marshaler for crud_fileStarter.Prepare()")
	}

	nms.mapper, ok = options["mapper"].(crud.Mapper)
	if !ok || nms.mapper == nil {
		return errors.New("no mapper for crud_fileStarter.Prepare()")
	}

	nms.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(crud.InterfaceKey)))

	return nil
}

func (nms *crud_fileStarter) Check() (info []starter.Info, err error) {
	fileinfo, err := os.Stat(nms.path)

	if err != nil {
		return nil, errors.Wrapf(err, "on check directory '%s'", nms.path)
	}

	if !fileinfo.IsDir() {
		return nil, errors.Errorf("'%s' isn't a directory", nms.path)
	}

	return nil, nil
}

func (nms *crud_fileStarter) Setup() error {
	err := os.MkdirAll(nms.path, 0755)
	if err != nil {
		return errors.Wrapf(err, "on create directory '%s'", nms.path)
	}

	return nil
}

func (nms *crud_fileStarter) Init(joiner joiner.Operator) error {

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
