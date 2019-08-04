package collector_comp

import (
	"log"
	"strconv"

	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/filelib"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"

	"github.com/pavlo67/partes/serverhttp/serverhttp_jschmhr"
	"github.com/pavlo67/punctum/notebook/notes"
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/things_old/files"

	"github.com/pavlo67/punctum/basis/viewshtml/front"

	"github.com/pavlo67/punctum/collector/importer"
	"github.com/pavlo67/punctum/processor.old/news"
	"github.com/pavlo67/punctum/processor/sources"
)

const InterfaceKey joiner.InterfaceKey = "collector.comp.comp"

// Starter ...
func Starter() starter.Operator {
	return &fountComponent{}
}

var endpoints, itemsEndpoints map[string]config.Endpoint

var listeners map[string]config.Listener
var params map[string]string

var fountOp sources.Operator
var flowOp news.Operator
var objectsOp notes.Operator
var filesOp files.Operator

var pagination uint64

var importerOperators = map[string]importer.Operator{}

const defaultPagination = 200

var idFountScannerGroup = "5"

var ImportFlowsPath string
var pathRepository string

type fountComponent struct{}

func (sc *fountComponent) Name() string {
	return string(InterfaceKey)
}

func (sc *fountComponent) Check(conf config.Config, indexPath string) (info []joiner.Info, err error) {

	index, errs := config.ComponentIndex(indexPath, filelib.CurrentPath(), nil)
	endpoints = index.Endpoints
	listeners = index.Listeners
	params = index.Params

	pagination, err = strconv.ParseUint(params["pagination"], 10, 64)
	if err != nil {
		log.Println("bad pagination value for fount interface:", pagination)
		pagination = defaultPagination
	}
	if params["idFountScannerGroup"] != "" {
		idFountScannerGroup = params["idFountScannerGroup"]
	}
	ImportFlowsPath = endpoints["importFlows"].ServerPath

	pathRepository, errs = conf.Paths("file_repository_path", errs)

	return nil, errs.Err()
}

func (sc *fountComponent) Setup(config.Config, string, map[string]string) error {

	// TODO: execute setup if necessary...
	return nil
}

var htmlFront string

func (sc *fountComponent) Init() error {

	itemsOp, ok := joiner.Component(items_comp.InterfaceKey).(front.Operator)
	if !ok {
		return errors.New("can't get component(items.comp).Operator interface for _hotel/hotel.comp")
	}
	itemsEndpoints = itemsOp.Endpoints()

	serverOp, ok := joiner.Component(serverhttp_jschmhr.InterfaceKey).(serverhttp_jschmhr.Operator)
	if !ok {
		return errors.New("no serverhttp_jschmhr.Operator interface found for fount component")
	}

	var errs basis.Errors
	htmlFront, errs = serverOp.HandleJS(filelib.CurrentPath(), InterfaceKey, nil, endpoints, listeners)
	initHTML()

	errsBack := serverOp.InitEndpoints(
		endpoints,
		fountTemplator,
		nil,
		htmlHandlers,
		restHandlers,
		nil,
	)
	errs = append(errs, errsBack...)

	err := joiner.JoinInterface(front.New(endpoints, listeners), InterfaceKey)
	if err != nil {
		errs = append(errs, err)
	}

	fountOp, ok = joiner.Component(sources.InterfaceKey).(sources.Operator)
	if !ok {
		return errors.New("can't get interface for fount component")
	}
	flowOp, ok = joiner.Component(news.InterfaceKey).(news.Operator)
	if !ok {
		return errors.New("can't get flow interface for fount component")
	}
	objectsOp, ok = joiner.Component(notes.InterfaceKey).(notes.Operator)
	if !ok {
		return errors.New("can't get import interface for fount component")
	}
	filesOp, ok = joiner.Component(files.InterfaceKey).(files.Operator)
	if !ok {
		return errors.New("no fileslocal interface found for fount component")
	}

	apps := joiner.InterfacesAll((*importer.Operator)(nil), "")
	for _, app := range apps {
		if importerOp, ok := app.Interface.(importer.Operator); ok {
			importerOperators[string(app.Key)] = importerOp
			log.Println("find importer Operator:", app.Key)
		}
	}

	return nil
}
