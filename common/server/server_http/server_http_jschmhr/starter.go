package server_http_jschmhr

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/errors"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/server"
	"github.com/pavlo67/workshop/common/server/server_http"
	"github.com/pavlo67/workshop/common/starter"
)

func Starter() starter.Operator {
	return &server_http_jschmhrStarter{}
}

var l logger.Operator
var _ starter.Operator = &server_http_jschmhrStarter{}

type server_http_jschmhrStarter struct {
	config server.Config

	noEventsOp bool

	interfaceKey joiner.InterfaceKey
}

func (ss *server_http_jschmhrStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *server_http_jschmhrStarter) Init(cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	var errs errors.Errors
	l = lCommon

	ss.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(server_http.InterfaceKey)))
	ss.noEventsOp = options.IsTrue("no_events_op")

	var cfgServerHTTP server.Config
	err := cfg.Value("server_http", &cfgServerHTTP)
	if err != nil {
		return nil, err
	}

	ss.config = cfgServerHTTP

	return nil, errs.Err()
}

func (ss *server_http_jschmhrStarter) Setup() error {
	return nil
}

func (ss *server_http_jschmhrStarter) Run(joinerOp joiner.Operator) error {

	authOpNil := auth.Operator(nil)
	authComps := joinerOp.InterfacesAll(&authOpNil)

	var authOps []auth.Operator
	for _, authComp := range authComps {
		if authOp, ok := authComp.Interface.(auth.Operator); ok {
			authOps = append(authOps, authOp)
		}
	}

	//var eventsOpSystem events.OperatorSystem
	//var eventsOp events.Operator
	//if !ss.noEventsOp {
	//	eventsOpSystem, _ = joinerOp.Interface(events.InterfaceSystemKey).(events.OperatorSystem)
	//	if eventsOpSystem == nil {
	//		return fmt.Errorf("no events.OperatorSystem with key %s", events.InterfaceSystemKey)
	//	}
	//
	//	eventsOp, _ = joinerOp.Interface(events.InterfaceKey).(events.Operator)
	//	if eventsOp == nil {
	//		return fmt.Errorf("no events.Operator with key %s", events.InterfaceKey)
	//	}
	//}

	srvOp, err := New(ss.config.Port, ss.config.TLSCertFile, ss.config.TLSKeyFile, authOps, ss.noEventsOp)
	if err != nil {
		return errors.Wrap(err, "can't init serverHTTPJschmhr.UserKey")
	}

	err = joinerOp.Join(srvOp, ss.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join serverHTTPJschmhr srvOp as server.UserKey with key '%s'", ss.interfaceKey)
	}

	err = joinerOp.Join(ss.config.Port, server_http.PortInterfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join port with key '%s'", server_http.PortInterfaceKey)
	}

	err = joinerOp.Join(ss.config.NoHTTPS, server_http.NoHTTPSInterfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join no_https with key '%s'", server_http.NoHTTPSInterfaceKey)
	}

	return nil
	// return srvOp.Start()

}
