package receiver_server_http

import (
	"os"
	"testing"

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
	"github.com/pavlo67/workshop/components/receiver"
	"github.com/pavlo67/workshop/components/router/router_stub"
	"github.com/pavlo67/workshop/components/sender"
	"github.com/pavlo67/workshop/components/sender/sender_http"

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
		{server_http_jschmhr.Starter(), common.Map{"port": port}},

		// transport system
		{packs_pg.Starter(), nil},
		{router_stub.Starter(), nil},
		{sender_http.Starter(), nil},
		{Starter(), common.Map{"handler_key": receiver.HandlerInterfaceKey}},

		{gatherer_actions.Starter(), common.Map{
			"receiver_handler_key": receiver.HandlerInterfaceKey,
		}},
	}

	joinerOp, err := starter.Run(starters, cfgCommon, cfgService, os.Args[1:], label)
	require.NoError(t, err)
	defer joinerOp.CloseAll()

	receiverOp, ok := joinerOp.Interface(receiver.InterfaceKey).(receiver.Operator)
	require.True(t, ok)
	require.NotNil(t, receiverOp)

	senderOp, ok := joinerOp.Interface(sender.InterfaceKey).(sender.Operator)
	require.True(t, ok)
	require.NotNil(t, senderOp)

	receiver.OperatorTestScenario(t, receiverOp, senderOp, l)
}
