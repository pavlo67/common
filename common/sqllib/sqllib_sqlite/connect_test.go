package sqllib_sqlite

import (
	"testing"

	"github.com/pavlo67/common/common/sqllib"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/apps"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/filelib"
)

func TestConnect(t *testing.T) {

	_, cfgService, _ := apps.PrepareTests(
		t,
		"test_service", "../../../"+apps.AppsSubpathDefault,
		"test",
		"", // "connect_test."+strconv.FormatInt(time.Now().Unix(), 10)+".log",
	)

	var cfgSqlite config.Access
	err := cfgService.Value("sqllib_sqlite", &cfgSqlite)
	require.NoError(t, err)

	cfgSqlite.Path, err = filelib.Dir(cfgSqlite.Path)
	require.NoError(t, err)
	require.NotEmpty(t, cfgSqlite.Path)

	cfgSqlite.Path += "test_connect.sqlite"

	db, err := Connect(cfgSqlite)
	require.NoError(t, err)
	require.NotNil(t, db)

	sqllib.TestDB(t, db)

}
