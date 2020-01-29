package transport

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"

	"github.com/pavlo67/workshop/components/packs"
	"github.com/pavlo67/workshop/components/runner"
)

//type OperatorTestCase struct {
//}

var createdAt = time.Now().UTC()

//func TestCases(PacksOp ActorKey, cleanerOp crud.Cleaner) []OperatorTestCase {
//	return []OperatorTestCase{
//		{
//			ActorKey: PacksOp,
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

const typeKeyTest crud.TypeKey = "test"
const typeKeyTestResponse crud.TypeKey = "test_response"

var paramsToTest = common.Map{
	"aa": "bb",
	"cc": float64(5),
}

var paramsToTestBytes, _ = json.Marshal(paramsToTest)

func OperatorTestScenario(t *testing.T, joinerOp joiner.Operator, transpOp Operator, l logger.Operator) {
	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	//for i, tc := range testCases {
	//	l.Debug(i)
	//}

	packOut := packs.Pack{
		Key:  "test1",
		From: "gatherer",
		To:   "gatherer",
		Data: crud.Data{
			TypeKey: typeKeyTest,
			Content: string(paramsToTestBytes),
		},
	}

	receiver := &actorReceiverEcho{l: l, t: t}
	err := joinerOp.Join(receiver, runner.DataInterfaceKey(typeKeyTest))
	require.NoError(t, err)

	receiverSender := &actorReceiverEchoTest{l: l, t: t}
	err = joinerOp.Join(receiverSender, runner.DataInterfaceKey(typeKeyTestResponse))
	require.NoError(t, err)

	_, packIn, err := transpOp.Send(&packOut)
	require.NoError(t, err)
	require.NotNil(t, packIn)

}

// receiver

var _ runner.Actor = &actorReceiverEcho{}

type actorReceiverEcho struct {
	l      logger.Operator
	t      *testing.T
	params common.Map
}

func (_ actorReceiverEcho) Name() string {
	return "actorReceiverEcho"
}

func (r *actorReceiverEcho) Init(params common.Map) (estimate *runner.Estimate, err error) {
	if r == nil {
		return nil, errors.New("no runner to init")
	}
	r.params = params
	return nil, nil
}

func (r actorReceiverEcho) Run() (info common.Map, posterior []joiner.Link, err error) {
	r.l.Infof("RECEIVER with params %#v", r.params)

	responseContent, _ := json.Marshal(r.params)

	info = common.Map{
		"response": crud.Data{
			TypeKey: typeKeyTestResponse,
			Content: string(responseContent),
		},
	}
	return info, nil, nil
}

// receiver/sender

var _ runner.Actor = &actorReceiverEchoTest{}

type actorReceiverEchoTest struct {
	l      logger.Operator
	t      *testing.T
	params common.Map
}

func (_ actorReceiverEchoTest) Name() string {
	return "actorReceiverEchoTest"
}

func (r *actorReceiverEchoTest) Init(params common.Map) (estimate *runner.Estimate, err error) {
	if r == nil {
		return nil, errors.New("no runner to init")
	}
	r.params = params
	return nil, nil
}

func (r actorReceiverEchoTest) Run() (info common.Map, posterior []joiner.Link, err error) {
	r.l.Infof("RECEIVER/SENDER with params %#v", r.params)

	// require.True(r.t, reflect.DeepEqual(paramsToTest, r.params))

	if reflect.DeepEqual(paramsToTest, r.params) {
		r.l.Info("!!!!!!!!!!!!!!!!!!!!")
	} else {
		r.l.Fatal("???????????????????")
	}

	// TODO: close HTTP server
	return nil, nil, nil
}
