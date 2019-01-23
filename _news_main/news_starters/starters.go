package news_starters

import (
	"github.com/pavlo67/punctum/starter"

	"github.com/pavlo67/punctum/auth/auth_ecdsa"
	"github.com/pavlo67/punctum/processor/founts/founts_leveldb"
	"github.com/pavlo67/punctum/processor/news/news_leveldb"
	"github.com/pavlo67/punctum/server/server_http/server_http_jschmhr"
)

func Starters() ([]starter.Starter, string) {
	var starters []starter.Starter

	// 1. operational interfaces
	starters = append(starters, starter.Starter{founts_leveldb.Starter(), nil})
	starters = append(starters, starter.Starter{news_leveldb.Starter(), nil})

	// 2. http server
	starters = append(starters, starter.Starter{auth_ecdsa.Starter(), nil})
	starters = append(starters, starter.Starter{server_http_jschmhr.Starter(), nil})

	return starters, "RSS BUILD"
}
