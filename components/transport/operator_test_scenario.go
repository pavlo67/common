package transport

import (
	"os"
	"testing"

	"github.com/pavlo67/workshop/common/identity"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/components/packs"
	"github.com/stretchr/testify/require"
)

//type OperatorTestCase struct {
//	Operator
//}
//
//var createdAt = time.Now().UTC()
//
//func TestCases(PacksOp Operator, cleanerOp crud.Cleaner) []OperatorTestCase {
//	return []OperatorTestCase{
//		{
//			Operator: PacksOp,
//			Cleaner:  cleanerOp,
//			ToSave: Pack{
//				Key: "test_key1",
//				From:        "test_key2",
//				To:          []identity.Key{"qwerwqer", "!!!"},
//				Options:     common.Map{"1": float64(2)},
//				TypeKey:     "no_type",
//				Content:     map[string]string{"6": ";klj"},
//				History: []crud.Action{{
//					Key:    crud.CreatedAction,
//					DoneAt: createdAt,
//				}},
//			},
//
//			ToAddHistory: []crud.Action{{
//				Key:    "action1",
//				DoneAt: time.Now().UTC(),
//			}},
//		},
//	}
//}

func OperatorTestScenario(t *testing.T, transpOp Operator, l logger.Operator) {
	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	//for i, tc := range testCases {
	//	l.Debug(i)
	//}

	transpOp.AddHandler("", TestTypeKey, &handlerEcho{})

	packOut := packs.Pack{
		Key:     "test1",
		From:    "test2",
		To:      identity.Key("gatherer_transport/aaa"),
		TypeKey: TestTypeKey,
		Content: map[string]interface{}{"aaa": "bbb"},
	}

	_, packIn, err := transpOp.Send(&packOut)
	require.NoError(t, err)
	require.NotNil(t, packIn)

	l.Infof("--> %#v", packOut)
	l.Infof("<-- %#v", *packIn)

	packOut.History = nil
	packIn.History = nil

	require.Equal(t, packOut, *packIn)

}

const TestTypeKey identity.Key = "test"

var _ packs.Handler = &handlerEcho{}

type handlerEcho struct{}

func (_ handlerEcho) Handle(pack *packs.Pack) (*packs.Pack, error) {
	return pack, nil
}
