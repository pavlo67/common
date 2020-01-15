package transport_http

import (
	"os"
	"testing"

	"github.com/pavlo67/workshop/common/scheduler/scheduler_timeout"
	"github.com/pavlo67/workshop/components/runner_factory/runner_factory_goroutine"
	"github.com/pavlo67/workshop/components/tasks/tasks_pg"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/auth/auth_ecdsa"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/control"
	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/serializer"
	"github.com/pavlo67/workshop/common/server/server_http/server_http_jschmhr"
	"github.com/pavlo67/workshop/common/starter"

	"github.com/pavlo67/workshop/components/packs/packs_pg"
	"github.com/pavlo67/workshop/components/transport"
	"github.com/pavlo67/workshop/components/transportrouter/transportrouter_stub"

	"github.com/pavlo67/workshop/apps/gatherer/gatherer_actions"
)

const serviceName = "gatherer"

func TestReceiverServerHTTP(t *testing.T) {
	env := "test"
	err := os.Setenv("ENV", env)
	require.NoError(t, err)

	l, err = logger.Init(logger.Config{})
	require.NoError(t, err)
	require.NotNil(t, l)

	currentPath := filelib.CurrentPath()

	configCommonPath := currentPath + "../../../environments/common." + env + ".yaml"
	cfgCommon, err := config.Get(configCommonPath, serializer.MarshalerYAML)
	require.NoError(t, err)

	var routesCfg map[string]config.Access
	err = cfgCommon.Value("routes", &routesCfg)
	require.NoError(t, err)

	serviceAccess, ok := routesCfg[serviceName]
	require.True(t, ok)

	port := serviceAccess.Port

	// gatherer config

	cfgServicePath := currentPath + "../../../environments/" + serviceName + "." + env + ".yaml"
	cfgService, err := config.Get(cfgServicePath, serializer.MarshalerYAML)
	require.NoError(t, err)

	// running starters

	label := "GATHERER/PG/TEST CLI BUILD"

	starters := []starter.Starter{

		// general purposes components
		{control.Starter(), nil},
		{auth_ecdsa.Starter(), nil},

		// action managers
		{tasks_pg.Starter(), nil},
		{runner_factory_goroutine.Starter(), nil},
		{scheduler_timeout.Starter(), nil},
		{server_http_jschmhr.Starter(), common.Map{"port": port}},

		// transport system
		{packs_pg.Starter(), nil},
		{transportrouter_stub.Starter(), nil},
		{Starter(), common.Map{"handler_key": transport.HandlerInterfaceKey, "domain": serviceName}},

		{gatherer_actions.Starter(), common.Map{
			"receiver_handler_key": transport.HandlerInterfaceKey,
		}},
	}

	joinerOp, err := starter.Run(starters, cfgCommon, cfgService, os.Args[1:], label)
	require.NoError(t, err)
	defer joinerOp.CloseAll()

	transpOp, ok := joinerOp.Interface(transport.InterfaceKey).(transport.Operator)
	require.True(t, ok)
	require.NotNil(t, transpOp)

	transport.OperatorTestScenario(t, joinerOp, transpOp, l)

	gatherer_actions.WG.Wait()
}
