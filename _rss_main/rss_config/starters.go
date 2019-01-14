package rss_config

import (
	"github.com/pavlo67/punctum/starter"

	"github.com/pavlo67/punctum/processor/founts/founts_leveldb"
	"github.com/pavlo67/punctum/processor/news/news_leveldb"
)

func Starters() ([]starter.Starter, string) {
	var starters []starter.Starter

	// 1. operational interfaces
	starters = append(starters, starter.Starter{founts_leveldb.Starter(), nil})
	starters = append(starters, starter.Starter{news_leveldb.Starter(), nil})

	return starters, "RSS BUILD"
}
