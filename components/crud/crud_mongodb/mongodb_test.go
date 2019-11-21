package crud_mongodb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/libraries/encodelib"
	"github.com/pavlo67/workshop/common/libraries/filelib"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/components/crud"
)

func TestCRUD(t *testing.T) {
	configPath := filelib.CurrentPath() + "../../../environments/test.yaml"

	cfg, err := config.Get(configPath, encodelib.MarshalerYAML)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	cfgMongoDB := config.Access{}

	err = cfg.Value("mongodb", &cfgMongoDB)
	require.NoError(t, err)

	l, err = logger.Init(logger.Config{})
	require.NoError(t, err)
	require.NotNil(t, l)

	crudOp, cleanerOp, mgoClient, err := NewCRUD(&cfgMongoDB, 5*time.Second, "crud", crud.Item{})
	require.NoError(t, err)

	testCases := []crud.OperatorTestCase{{
		Operator: crudOp,
		Cleaner:  cleanerOp,
		ToSave:   crud.Item{},
		ToUpdate: crud.Item{},
	}}

	crud.OperatorTest(t, testCases)

	mgoClient.Disconnect(nil)
}
