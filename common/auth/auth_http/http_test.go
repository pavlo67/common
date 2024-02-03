package auth_http

import (
	"testing"
	"time"

	"github.com/pavlo67/common/apps/demo"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/starter"
)

func TestAuthHTTP(t *testing.T) {

	cfgService, l := config.PrepareTests(t, "../../../_envs/", "test", "")
	require.NotNil(t, cfgService)

	starters, err := demo.Components(cfgService, true)
	require.NoError(t, err)

	starters = append(
		starters,
		starter.Component{Starter(), common.Map{
			// "prefix":        demo_server_http.PrefixREST,
			"server_config": demo.ServerConfig,
		}},
	)

	joinerOp, err := starter.Run(starters, &cfgService, "CLI BUILD FOR TEST", l)
	require.NoError(t, err)
	require.NotNil(t, joinerOp)
	defer joinerOp.CloseAll()

	authOp, _ := joinerOp.Interface(InterfaceKey).(auth.Operator)
	require.NotNil(t, authOp)

	//srvOp, _ := joinerOp.Interface(server_http.InterfaceKey).(server_http.Operator)
	//require.NotNil(t, authOp)

	err = demo.Run(joinerOp, false, l)
	require.NoError(t, err)

	time.Sleep(time.Second)

	auth.OperatorTestScenarioPassword(t, authOp)

	//err = srvOp.Shutdown(context.TODO())
	//require.NoError(t, err)

}
