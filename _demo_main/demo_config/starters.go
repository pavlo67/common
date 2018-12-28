package demo_config

import (
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/libs/filelib"
	"github.com/pavlo67/punctum/basis/starter"
	"github.com/pavlo67/punctum/server_http/server_http_jschmhr"

	"github.com/pavlo67/punctum/_demo_main/demo_server_http"
	"github.com/pavlo67/punctum/identity/identity_ecdsa"
)

func Starters() ([]starter.Starter, string) {
	paramsServerStatic := basis.Params{
		"static_path":   filelib.CurrentPath() + "../demo_static/",
		"template_path": filelib.CurrentPath() + "../demo_static/demo_server.html",
	}

	var starters []starter.Starter

	starters = append(starters, starter.Starter{identity_ecdsa.Starter(), nil})
	starters = append(starters, starter.Starter{server_http_jschmhr.Starter(), paramsServerStatic})
	starters = append(starters, starter.Starter{demo_server_http.Starter(), nil})

	return starters, "PUNCTUM DEMO BUILD"
}
