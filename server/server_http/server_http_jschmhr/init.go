package server_http_jschmhr

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/auth"
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/filelib"
	"github.com/pavlo67/punctum/server/controller"
	"github.com/pavlo67/punctum/server/server_http"
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"
	"github.com/pavlo67/punctum/starter/logger"
)

func Starter() starter.Operator {
	return &server_http_jschmhrStarter{}
}

var l logger.Operator
var _ starter.Operator = &server_http_jschmhrStarter{}

type server_http_jschmhrStarter struct {
	interfaceKey       joiner.InterfaceKey
	interfaceKeyRouter joiner.InterfaceKey
	config             config.ServerTLS

	htmlTemplate string
	staticPath   string
}

func (ss *server_http_jschmhrStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *server_http_jschmhrStarter) Prepare(conf *config.Config, options, runtimeOptions basis.Info) error {
	l = logger.Get()

	var errs basis.Errors

	ss.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(server_http.InterfaceKey)))
	ss.interfaceKeyRouter = joiner.InterfaceKey(options.StringDefault("interface_key_router", string(controller.InterfaceKey)))

	ss.config, errs = conf.Server(options.StringDefault("config_server_key", "default"), errs)
	if ss.config.Port <= 0 {
		errs = append(errs, fmt.Errorf("wrong port for serverOp: %d", ss.config.Port))
	}

	templatePath := options.StringDefault("template_path", "")
	if templatePath == "" {
		l.Warn(`on server_http_jschmhr.Init(): empty options["template_path"]`)

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

	ss.staticPath = options.StringDefault("static_path", "")

	return errs.Err()
}

func (ss *server_http_jschmhrStarter) Check() (info []starter.Info, err error) {
	return nil, nil
}

func (ss *server_http_jschmhrStarter) Setup() error {
	return nil
}

func (ss *server_http_jschmhrStarter) Init(joiner joiner.Operator) error {
	identOpsMap := map[auth.CredsType][]auth.Operator{}

	identOpsPtr := joiner.ComponentsAll(auth.InterfaceKey)
	for _, identOpIntf := range identOpsPtr {
		if identOp, ok := identOpIntf.Interface.(auth.Operator); ok {
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
	)
	if err != nil {
		return errors.Wrap(err, "can't init serverHTTPJschmhr.Operator")
	}

	if ss.staticPath != "" {
		srvOp.HandleFiles("/static/*filepath", ss.staticPath, nil)
	}

	err = joiner.JoinInterface(srvOp, ss.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join serverHTTPJschmhr srvOp as server.Operator with key '%s'", ss.interfaceKey)
	}

	err = joiner.JoinInterface(srvOp, ss.interfaceKeyRouter)
	if err != nil {
		return errors.Wrapf(err, "can't join serverHTTPJschmhr srvOp as router.Operator with key '%s'", ss.interfaceKeyRouter)
	}

	return nil
}
