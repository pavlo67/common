package demo_api

import (
	"fmt"
	"strconv"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/libraries/filelib"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/server/server_http"
	"github.com/pavlo67/common/common/starter"
)

func Starter() starter.Operator {
	return &demoStarter{}
}

var l logger.Operator

var _ starter.Operator = &demoStarter{}

type demoStarter struct {
	prefix string
	// baseDir string

	// skipAbsentEndpoints bool
}

func (ss *demoStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *demoStarter) Init(cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon
	if l == nil {
		return nil, fmt.Errorf("no logger for %s:-(", ss.Name())
	}

	ss.prefix = options.StringDefault("prefix", "")
	// ss.skipAbsentEndpoints = options.IsTrue("skip_absent_endpoints")

	return nil, nil
}

func (ss *demoStarter) Setup() error {
	return nil
}

// Swagger-UI sorts sections due to the first their path occurrences, so:
// 1. unauthorized       /auth
// 2. admin              /front/add_plan
// 3. any_authenticated

var Endpoints = server_http.Endpoints{
	auth.EPAuth: {Path: "/auth", Tags: []string{"unauthorized"}, HandlerKey: auth.AuthHandlerKey},
}

func (ss *demoStarter) Run(joinerOp joiner.Operator) error {
	srvOp, ok := joinerOp.Interface(server_http.InterfaceKey).(server_http.Operator)
	if !ok {
		return fmt.Errorf("no server_http.UserKey with key %s", server_http.InterfaceKey)
	}

	srvPort, _ := joinerOp.Interface(server_http.PortInterfaceKey).(int)
	if srvPort <= 0 {
		return fmt.Errorf("no server_http.Port with key %s", server_http.PortInterfaceKey)
	}

	isHTTPS, ok := joinerOp.Interface(server_http.HTTPSInterfaceKey).(bool)
	if !ok {
		return fmt.Errorf("no server_http.HTTPS info with key %s", server_http.HTTPSInterfaceKey)
	}

	//err = joinerOp.Join(ss.config.NoHTTPS, server_http.SwaggerNoHTTPSKey)
	//if err != nil {
	//	return errors.Wrapf(err, "can't join no_https with key '%s'", server_http.SwaggerNoHTTPSKey)
	//}
	//noHTTPS := joinerOp.Interface(server_http.SwaggerNoHTTPSKey).(bool)
	//if !ok {
	//	return fmt.Errorf("no server_http.NoHTTPS with key %s", server_http.SwaggerNoHTTPSKey)
	//}

	for key, ep := range Endpoints {
		ep.Handler, ok = joinerOp.Interface(ep.HandlerKey).(*server_http.Endpoint)
		if ok {
			Endpoints[key] = ep
			//} else if ss.skipAbsentEndpoints {
			//	ep.Skip = true
			//	Endpoints[key] = ep
		} else {
			return fmt.Errorf("no server_http.Endpoint with key %s", ep.HandlerKey)
		}
	}

	cfg := server_http.Config{
		Title:     "Demo REST API",
		Version:   "0.0.1",
		Prefix:    ss.prefix,
		Endpoints: Endpoints,
	}

	if err := server_http.InitEndpointsWithSwaggerV2(
		cfg,
		":"+strconv.Itoa(srvPort),
		!isHTTPS,
		srvOp,
		filelib.CurrentPath()+"api-docs/",
		"api-docs",
		l,
	); err != nil {
		l.Error("on server_http.InitEndpointsWithSwaggerV2(): ", err)
	}

	WG.Add(1)

	// TODO!!! customize it
	// if isHTTPS {
	//	go http.ListenAndServe(":80", http.HandlerFunc(server_http.Redirect))
	// }

	go func() {
		defer WG.Done()
		if err := srvOp.Start(); err != nil {
			l.Error("on srvOp.Start(): ", err)
		}
	}()

	return nil
}
