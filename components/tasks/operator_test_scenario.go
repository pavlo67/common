package tasks

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/selectors"
)

type OperatorTestCase struct {
	Operator
	crud.Cleaner

	ToSave       Task
	ToSetResults Result
}

func TestCases(tasksOp Operator, cleanerOp crud.Cleaner) []OperatorTestCase {
	startedAt := time.Now().UTC()
	finishedAt := startedAt.Add(time.Second).UTC()

	return []OperatorTestCase{
		{
			Operator: tasksOp,
			Cleaner:  cleanerOp,
			ToSave: Task{
				WorkerType: "wt0",
				Params:     common.Map{"1": float64(2), "3": "4"},
			},

			ToSetResults: Result{
				Timing: Timing{
					StartedAt:  startedAt,
					FinishedAt: finishedAt,
				},
				Success:   true,
				Info:      common.Map{"1": float64(4)},
				Posterior: []common.ID{"6", "8"},
			},
		},
	}
}

const numRepeats = 3
const toReadI = 0       // must be < numRepeats
const toSetResultsI = 1 // must be < numRepeats
const toDeleteI = 2     // must be < numRepeats

func ChechReaded(t *testing.T, readed *Item, expectedID common.ID, expectedTask Task, l logger.Operator) {
	require.NotNil(t, readed)

	l.Infof("was saved: %#v", expectedTask)
	l.Infof("is readed: %#v", readed)

	require.Equal(t, expectedTask, readed.Task)
	require.Equal(t, expectedTask.Params, readed.Task.Params)
	require.Equal(t, expectedID, readed.ID)
	require.True(t, readed.History.CreatedAt.After(time.Time{}))
	require.True(t, readed.History.CreatedAt.Before(time.Now()))

	// TODO!!! check .Status
}

func OperatorTestScenario(t *testing.T, testCases []OperatorTestCase, l logger.Operator) {
	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	for i, tc := range testCases {
		l.Debug(i)

		var id [numRepeats]common.ID
		var toSave [numRepeats]Task

		// ClearDatabase ---------------------------------------------------------------------------------

		err := tc.Cleaner.Clean(nil, nil)
		require.NoError(t, err, "what is the error on .Cleaner()?")

		// test .Save ------------------------------------------------------------------------------------

		for i := 0; i < numRepeats; i++ {
			toSave[i] = tc.ToSave
			idI, err := tc.Save(toSave[i], nil)
			require.NoError(t, err)
			require.NotEqual(t, common.ID(""), idI)
			id[i] = idI
		}

		// test .Read ------------------------------------------------------------------------------------

		readedSaved, err := tc.Read(id[toReadI], nil)
		require.NoError(t, err)

		ChechReaded(t, readedSaved, id[toReadI], tc.ToSave, l)

		// test .SetResults & .Read ----------------------------------------------------------------------

		err = tc.SetResult(id[toSetResultsI], tc.ToSetResults, nil)
		require.NoError(t, err)

		readedUpdated, err := tc.Read(id[toSetResultsI], nil)
		require.NoError(t, err)
		require.NotNil(t, readedUpdated)
		for i, r := range readedUpdated.Results {
			readedUpdated.Results[i].StartedAt = r.StartedAt.UTC()
			readedUpdated.Results[i].FinishedAt = r.FinishedAt.UTC()
		}

		ChechReaded(t, readedUpdated, id[toSetResultsI], tc.ToSave, l)
		require.Equal(t, 1, len(readedUpdated.Results))
		require.True(t, reflect.DeepEqual(tc.ToSetResults, readedUpdated.Results[0]), fmt.Sprintf("\nexpected = %#v\n  readed = %#v", tc.ToSetResults, readedUpdated.Results[0]))

		// test .SetResults & .Read again  ---------------------------------------------------------------

		err = tc.SetResult(id[toSetResultsI], tc.ToSetResults, nil)
		require.NoError(t, err)

		readedUpdated2, err := tc.Read(id[toSetResultsI], nil)
		require.NoError(t, err)
		require.NotNil(t, readedUpdated2)
		for i, r := range readedUpdated2.Results {
			readedUpdated2.Results[i].StartedAt = r.StartedAt.UTC()
			readedUpdated2.Results[i].FinishedAt = r.FinishedAt.UTC()
		}

		ChechReaded(t, readedUpdated2, id[toSetResultsI], tc.ToSave, l)
		require.Equal(t, 2, len(readedUpdated2.Results))
		require.True(t, reflect.DeepEqual(tc.ToSetResults, readedUpdated2.Results[0]))
		require.True(t, reflect.DeepEqual(tc.ToSetResults, readedUpdated2.Results[1]))

		// check if another records are unchanged

		readedSaved, err = tc.Read(id[toReadI], nil)
		require.NoError(t, err)

		ChechReaded(t, readedSaved, id[toReadI], tc.ToSave, l)
		require.Equal(t, 0, len(readedSaved.Results))

		readedSaved, err = tc.Read(id[toDeleteI], nil)
		require.NoError(t, err)

		ChechReaded(t, readedSaved, id[toDeleteI], tc.ToSave, l)
		require.Equal(t, 0, len(readedSaved.Results))

		// test List -------------------------------------------------------------------------------------

		itemsAll, err := tc.List(nil, nil)
		require.NoError(t, err)
		require.True(t, len(itemsAll) == numRepeats)

		ChechReaded(t, &itemsAll[toReadI], id[toReadI], tc.ToSave, l)
		ChechReaded(t, &itemsAll[toSetResultsI], id[toSetResultsI], tc.ToSave, l)
		ChechReaded(t, &itemsAll[toDeleteI], id[toDeleteI], tc.ToSave, l)

		itemsOne, err := tc.List(selectors.In("id", id[toDeleteI]), nil)
		require.NoError(t, err)
		require.True(t, len(itemsOne) == 1)

		ChechReaded(t, &itemsOne[0], id[toDeleteI], tc.ToSave, l)

		// test .Remove ----------------------------------------------------------------------------------

		//err = tc.Remove(id[toDeleteI], nil)
		//require.NoError(t, err)
		//
		//readDeleted, err := tc.Read(id[toDeleteI], nil)
		//require.Error(t, err)
		//require.Nil(t, readDeleted)
		//
		//itemsAll, err = tc.List(nil, nil)
		//require.NoError(t, err)
		//require.True(t, len(itemsAll) == numRepeats-1)
	}
}
