package crud_mongodb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"os"

	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/libraries/encodelib"
	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/components/crud"
)

type Test struct {
	AAA string
	BBB int
}

func TestCRUD(t *testing.T) {
	env := "test"
	err := os.Setenv("ENV", env)
	require.NoError(t, err)

	l, err = logger.Init(logger.Config{})
	require.NoError(t, err)
	require.NotNil(t, l)

	configPath := filelib.CurrentPath() + "../../../environments/" + env + ".yaml"
	cfg, err := config.Get(configPath, encodelib.MarshalerYAML)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	cfgMongoDB := config.Access{}
	err = cfg.Value("mongodb", &cfgMongoDB)
	require.NoError(t, err)

	crudOp, cleanerOp, mgoClient, err := NewCRUD(&cfgMongoDB, 5*time.Second, "test", "crud", crud.Item{Details: Test{}})
	require.NoError(t, err)

	testCases := []crud.OperatorTestCase{{
		Operator:      crudOp,
		Cleaner:       cleanerOp,
		DetailsToRead: &Test{},
		ToSave: crud.Item{
			Title:   "345456",
			Summary: "6578gj",
			URL:     "",
			Details: Test{
				AAA: "aaa",
				BBB: 222,
			},
		},
		ToUpdate: crud.Item{},
	}}

	crud.OperatorTestScenario(t, testCases, l)

	mgoClient.Disconnect(nil)
}
