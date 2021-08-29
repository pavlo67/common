package auth_http

import (
	"testing"

	"github.com/pavlo67/common/common/auth/auth_server_http"

	"github.com/pavlo67/common/common/auth"

	"github.com/pavlo67/common/apps/demo/demo_settings"
	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/config"
)

func TestAuthHTTP(t *testing.T) {

	var cfgService config.Config
	cfgService, l = config.PrepareTests(
		t,
		"../../../_environments/",
		"test",
		"", // "connect_test."+strconv.FormatInt(time.Now().Unix(), 10)+".log",
	)

	var cfgServerHTTP config.Access
	err := cfgService.Value("server_http", &cfgServerHTTP)
	require.NoError(t, err)

	serverConfig := demo_settings.ServerConfig

	err = serverConfig.CompleteDirectly(auth_server_http.Endpoints, cfgServerHTTP.Host, cfgServerHTTP.Port, demo_settings.PrefixREST)
	require.NoError(t, err)

	authOp, err := New(serverConfig)
	require.NoError(t, err)
	require.NotNil(t, authOp)

	auth.OperatorTestScenarioPassword(t, authOp)

}
