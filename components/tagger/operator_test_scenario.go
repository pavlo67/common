package tagger

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
)

type TagsToChange struct {
	Action          string
	Key             joiner.InterfaceKey
	ID              common.ID
	Tags            []Tag
	IsErrorExpected bool
}

type TagToCheck struct {
	Tag             Tag
	Tagged          crud.Index
	IsErrorExpected bool
}

type TestStep struct {
	TagsToChange
	TagsToCheck []TagToCheck
}

type TestCase struct {
	Operator Operator
	Steps    []TestStep
}

func QueryTagsTestCases(taggerOp Operator) []TestCase {
	id1 := common.ID("11")
	id2 := common.ID("22")

	tags1 := []Tag{"1", "2", "3"}
	tags2 := []Tag{"5", "6", "3"}

	key := string(InterfaceKey)

	return []TestCase{
		// 0 all ok
		{
			Operator: taggerOp,
			Steps: []TestStep{
				{
					TagsToChange: TagsToChange{
						Action: "save",
						Key:    InterfaceKey,
						ID:     id1,
						Tags:   tags1,
					},
					TagsToCheck: []TagToCheck{
						{Tag: "1", Tagged: crud.Index{key: []common.ID{id1}}},
						{Tag: "2", Tagged: crud.Index{key: []common.ID{id1}}},
						{Tag: "3", Tagged: crud.Index{key: []common.ID{id1}}},
						{Tag: "4", Tagged: crud.Index{}},
						{Tag: "5", Tagged: crud.Index{}},
						{Tag: "6", Tagged: crud.Index{}},
					},
				},
				{
					TagsToChange: TagsToChange{
						Action: "save",
						Key:    InterfaceKey,
						ID:     id2,
						Tags:   tags2,
					},
					TagsToCheck: []TagToCheck{
						{Tag: "1", Tagged: crud.Index{key: []common.ID{id1}}},
						{Tag: "2", Tagged: crud.Index{key: []common.ID{id1}}},
						{Tag: "3", Tagged: crud.Index{key: []common.ID{id1, id2}}},
						{Tag: "4", Tagged: crud.Index{}},
						{Tag: "5", Tagged: crud.Index{key: []common.ID{id2}}},
						{Tag: "6", Tagged: crud.Index{key: []common.ID{id2}}},
					},
				},
				{
					TagsToChange: TagsToChange{
						Action: "replace",
						Key:    InterfaceKey,
						ID:     id1,
						Tags:   nil,
					},
					TagsToCheck: []TagToCheck{
						{Tag: "1", Tagged: crud.Index{}},
						{Tag: "2", Tagged: crud.Index{}},
						{Tag: "3", Tagged: crud.Index{key: []common.ID{id2}}},
						{Tag: "4", Tagged: crud.Index{}},
						{Tag: "5", Tagged: crud.Index{key: []common.ID{id2}}},
						{Tag: "6", Tagged: crud.Index{key: []common.ID{id2}}},
					},
				},
			},
		},
	}
}

func OperatorTestScenario(t *testing.T, testCases []TestCase, cleanerOp crud.Cleaner, l logger.Operator) {
	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	for i, tc := range testCases {
		l.Infof("test #%d", i)

		err := cleanerOp.Clean(nil)
		require.NoError(t, err)

		for j, step := range tc.Steps {
			l.Infof("\tstep #%d", j)

			var err error
			switch step.Action {
			case "save":
				err = tc.Operator.SaveTags(step.Key, step.ID, step.Tags, nil)
			case "remove":
				err = tc.Operator.RemoveTags(step.Key, step.ID, step.Tags, nil)
			case "replace":
				err = tc.Operator.ReplaceTags(step.Key, step.ID, step.Tags, nil)
			case "tags":
				var tags []Tag
				tags, err = tc.Operator.ListTags(step.Key, step.ID, nil)
				if !step.TagsToChange.IsErrorExpected {
					require.Equal(t, step.Tags, tags)
				}
			case "":
				l.Debug("no action!")
			default:
				l.Errorf("wrong action: '%s'", step.Action)
			}

			if step.TagsToChange.IsErrorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			for _, tagToCheck := range step.TagsToCheck {
				tagged, err := tc.Operator.IndexWithTag(tagToCheck.Tag, nil)

				if tagToCheck.IsErrorExpected {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					require.Equal(t, tagToCheck.Tagged, tagged)
				}
			}
		}
	}
}

//func prepareTest(t *testing.T, operator Operator, settaggerSteps []SettaggerStep) {
//	// ClearDatabase ------------------------------------------------------------------------------------
//
//	//err := operator.go.Clean()
//	//require.NoError(t, err, "what is the error on .Clean()?")
//
//	// test Settagger --------------------------------------------------------------------------------------
//
//	for j, sl := range settaggerSteps {
//		fmt.Println("." + strconv.Itoa(j))
//
//		if sl.ExpectedErr != nil {
//			_, err := operator.Settagger(sl.IS, sl.LinkedKey, sl.LinkedID, sl.Tags)
//			require.Error(t, err, "where is an error on .Seltagger(%#v, %s, %s, %#v)?", sl.IS, sl.LinkedKey, sl.LinkedID, sl.Tags)
//			continue
//		}
//
//		if sl.ISBad != nil {
//			_, err := operator.Settagger(*sl.ISBad, sl.LinkedKey, sl.LinkedID, sl.Tags)
//			require.Error(t, err, "where is an error on .Seltagger(%#v, %s, %s, %#v)?", *sl.ISBad, sl.LinkedKey, sl.LinkedID, sl.Tags)
//		}
//
//		res, err := operator.Settagger(sl.IS, sl.LinkedKey, sl.LinkedID, sl.Tags)
//		require.NoError(t, err, "what is an error on .Seltagger(%#v, %s, %s, %#v)?", sl.IS, sl.LinkedKey, sl.LinkedID, sl.Tags)
//		require.Equal(t, len(sl.ExpectedTagInfo), len(res), "len(sl.ExpectedTagInfo = %#v) != len(res = %#v)", sl.ExpectedTagInfo, res)
//
//		exp := mapLinkedInfo(sl.ExpectedTagInfo)
//		act := mapLinkedInfo(res)
//		for objectID, countLinked := range exp {
//			require.Equal(t, countLinked, act[objectID], "linkedInfo[%s] isn't correct (%#v)", objectID, res)
//		}
//
//	}
//
//}

//func QueryByObjectIDTest(t *testing.T, testCases []QueryByObjectIDTestCase) {
//	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
//		t.Fatal("No test environment!!!")
//	}
//
//	for i, tc := range testCases {
//		fmt.Println("QueryByObjectIDTest: ", i)
//
//		prepareTest(t, tc.Operator, tc.SettaggerSteps)
//
//		// test QueryByObjectID --------------------------------------------------------------------------------------
//
//		if tc.ExpectedErr != nil {
//			_, err := tc.Operator.QueryByObjectID(tc.IS, tc.ObjectID)
//			require.Error(t, err, "where is an error on .QueryByObjectID(%#v, %s)?", tc.IS, tc.ObjectID)
//			continue
//		}
//
//		if tc.ISBad != nil {
//			_, err := tc.Operator.QueryByObjectID(*tc.ISBad, tc.ObjectID)
//			require.Error(t, err, "where is an error on .QueryByObjectID(%#v, %s)?", *tc.ISBad, tc.ObjectID)
//		}
//
//		linked, err := tc.Operator.QueryByObjectID(tc.IS, tc.ObjectID)
//		require.NoError(t, err, "what is an error on .QueryByObjectID(%#v, %s)?", tc.IS, tc.ObjectID)
//		require.Equal(t, len(tc.ExpectedLinked), len(linked), "len(tc.ExpectedLinked = %#v) != len(linked = %#v)", tc.ExpectedLinked, linked)
//
//		sort.Sort(byLinked(tc.ExpectedLinked))
//		sort.Sort(byLinked(linked))
//
//		for i, li := range tc.ExpectedLinked {
//			require.Equal(t, Hash(li), Hash(linked[i]), "linked[%d] isn't correct", i)
//		}
//
//	}
//}

//func mapLinkedInfo(linkedInfo []LinkedInfo) map[string]uint {
//	res := map[string]uint{}
//	for _, v := range linkedInfo {
//		res[v.ObjectID] = v.CountLinked
//	}
//
//	return res
//}
//
//func mapTagInfo(tagInfo []TagInfo) map[string]uint64 {
//	res := map[string]uint64{}
//	for _, v := range tagInfo {
//		res[v.Tag] = v.Count
//	}
//
//	return res
//}
//
//func Hash(li Linked) string {
//	return li.ObjectID + " " + li.Type + " " + li.LinkedID + " " + li.LinkedType + " " + li.Tag
//}
//
//type byLinked []Linked
//
//func (s byLinked) Len() int {
//	return len(s)
//}
//func (s byLinked) Swap(i, j int) {
//	s[i], s[j] = s[j], s[i]
//}
//
//func (s byLinked) Less(i, j int) bool {
//	return Hash(s[i]) < Hash(s[j])
//}
