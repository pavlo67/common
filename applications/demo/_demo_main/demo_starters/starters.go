package demo_starters

import (
	"github.com/pavlo67/constructor/applications/demo/_demo_main/demo_server_http"
	"github.com/pavlo67/constructor/components/auth/auth_ecdsa"
	"github.com/pavlo67/constructor/components/common"
	"github.com/pavlo67/constructor/components/common/filelib"
	"github.com/pavlo67/constructor/components/common/starter"
	"github.com/pavlo67/constructor/components/server/server_http/server_http_jschmhr"
)

func Starters() ([]starter.Starter, string) {
	paramsServerStatic := common.Info{
		"static_path": filelib.CurrentPath() + "../demo_server_http/static/",
	}

	var starters []starter.Starter

	starters = append(starters, starter.Starter{auth_ecdsa.Starter(), nil})
	starters = append(starters, starter.Starter{server_http_jschmhr.Starter(), paramsServerStatic})
	starters = append(starters, starter.Starter{demo_server_http.Starter(), nil})

	return starters, "PUNCTUM DEMO BUILD"
}
