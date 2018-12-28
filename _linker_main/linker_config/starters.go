package linker_config

import (
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/starter"
	"github.com/pavlo67/punctum/identity/identity_ecdsa"
	"github.com/pavlo67/punctum/point/singlepoint_server_http"
	"github.com/pavlo67/punctum/server_http/server_http_jschmhr"
)

func Starters() ([]starter.Starter, string) {

	var starters []starter.Starter

	starters = append(starters, starter.Starter{identity_ecdsa.Starter(), nil})
	starters = append(starters, starter.Starter{server_http_jschmhr.Starter(), nil})
	starters = append(starters, starter.Starter{singlepoint_server_http.Starter(), basis.Params{"name": "linker"}})

	return starters, "LINKER SERVER BUILD"
}
