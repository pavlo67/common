package auth_http

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/common/apps/demo/demo_server_http"
	"github.com/pavlo67/common/apps/demo/demo_settings"
)

func TestAuthHTTP(t *testing.T) {

	cfgService, l := config.PrepareTests(t, "../../../_envs/", "test", "")
	require.NotNil(t, cfgService)

	starters, err := demo_settings.Starters(cfgService, true)
	require.NoError(t, err)

	starters = append(
		starters,
		starter.Component{Starter(), common.Map{
			// "prefix":        demo_server_http.PrefixREST,
			"server_config": demo_server_http.ServerConfig,
		}},
	)

	joinerOp, err := starter.Run(starters, &cfgService, "CLI BUILD FOR TEST", l)
	require.NoError(t, err)
	require.NotNil(t, joinerOp)
	defer joinerOp.CloseAll()

	time.Sleep(time.Second)

	authOp, _ := joinerOp.Interface(InterfaceKey).(auth.Operator)
	require.NotNil(t, authOp)

	auth.OperatorTestScenarioPassword(t, authOp)

}
