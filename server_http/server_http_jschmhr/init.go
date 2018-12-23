package server_http_jschmhr

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/config"
	"github.com/pavlo67/punctum/basis/filelib"
	"github.com/pavlo67/punctum/basis/logger"
	"github.com/pavlo67/punctum/basis/program"
	"github.com/pavlo67/punctum/basis/starter"
	"github.com/pavlo67/punctum/identity"
	"github.com/pavlo67/punctum/server_http"
)

func Starter() starter.Operator {
	return &server_http_jschmhrStarter{}
}

var l *zap.SugaredLogger
var _ starter.Operator = &server_http_jschmhrStarter{}

type server_http_jschmhrStarter struct {
	interfaceKey program.InterfaceKey
	config       config.ServerTLS

	htmlTemplate string
}

func (ss *server_http_jschmhrStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *server_http_jschmhrStarter) Prepare(conf *config.PunctumConfig, params basis.Params) error {
	l = logger.Get()

	var errs basis.Errors

	ss.interfaceKey = program.InterfaceKey(params.StringKeyDefault("interface_key", string(server_http.InterfaceKey)))

	ss.config, errs = conf.Server(params.StringKeyDefault("config_key", "data"), errs)
	if ss.config.Port <= 0 {
		errs = append(errs, fmt.Errorf("wrong port for serverOp: %d", ss.config.Port))
	}

	templatePath := params.StringKeyDefault("template_path", "")
	if templatePath == "" {
		l.Warn(`on serverHTTPJschmhr.Prepare(): empty params["template_path"]`)

	} else {
		if templatePath[0] != '/' {
			templatePath = filelib.CurrentPath() + templatePath
		}

		htmlTemplate, err := ioutil.ReadFile(templatePath)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "error reading template data from '%s'", templatePath))
		} else if len(htmlTemplate) < 1 {
			errs = append(errs, errors.Errorf("empty template data file: '%s'", templatePath))
		}
		ss.htmlTemplate = string(htmlTemplate)
	}

	return errs.Err()
}

func (ss *server_http_jschmhrStarter) Check() (info []starter.Info, err error) {
	return nil, nil
}

func (ss *server_http_jschmhrStarter) Setup() error {
	return nil
}

func (ss *server_http_jschmhrStarter) Init(joiner program.Joiner) error {
	identOpsMap := map[identity.CredsType][]identity.Operator{}

	identOpsPtr := joiner.InterfacesAll(identity.InterfaceKey)
	for _, identOpIntf := range identOpsPtr {
		if identOp, ok := identOpIntf.Interface.(identity.Operator); ok {
			credsTypes, err := identOp.Accepts()
			if err != nil {
				l.Error(err)
			}
			for _, credsType := range credsTypes {
				identOpsMap[credsType] = append(identOpsMap[credsType], identOp)
			}
		}
	}

	srvOp, err := New(
		ss.config.Port,
		ss.config.TLSCertFile,
		ss.config.TLSKeyFile,
		identOpsMap,
		ss.htmlTemplate,
	)
	if err != nil {
		return errors.Wrap(err, "can't init serverHTTPJschmhr.Operator")
	}

	err = joiner.JoinInterface(srvOp, ss.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join serverHTTPJschmhr srvOp as server.Operator with key '%s'", ss.interfaceKey)
	}

	return nil
}
