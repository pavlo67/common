package news_starters

import (
	"github.com/pavlo67/punctum/starter"

	"github.com/pavlo67/punctum/auth/auth_ecdsa"
	"github.com/pavlo67/punctum/processor/founts/founts_leveldb"
	"github.com/pavlo67/punctum/processor/news/news_leveldb"
	"github.com/pavlo67/punctum/processor/news/router_news"
	"github.com/pavlo67/punctum/server/server_http/server_http_jschmhr"
)

func Starters(routerStarters ...starter.Starter) ([]starter.Starter, string) {
	var starters = []starter.Starter{
		// 1. operational interfaces
		{founts_leveldb.Starter(), nil},
		{news_leveldb.Starter(), nil},
	}

	return append(starters, routerStarters...), "NEWS BUILD"
}

func RouterHTTPStarters() []starter.Starter {
	return []starter.Starter{
		{auth_ecdsa.Starter(), nil},
		{server_http_jschmhr.Starter(), nil},
		{router_news.Starter(), nil},
	}
}