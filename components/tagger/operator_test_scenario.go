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
	ToTag           joiner.Link
	Tags            []Tag
	IsErrorExpected bool
}

type TagToCheck struct {
	Tag             Tag
	Tagged          Index
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

	params1 := common.Map{"a": "b"}
	params2 := common.Map{"c": "d"}

	tags1 := []Tag{{"1", params1}, {"2", nil}, {"3", nil}}
	tags2 := []Tag{{"3", params2}, {"5", nil}, {"6", nil}}

	key := InterfaceKey

	return []TestCase{
		// 0 all ok
		{
			Operator: taggerOp,
			Steps: []TestStep{
				{
					TagsToChange: TagsToChange{
						Action: "add",
						ToTag:  joiner.Link{InterfaceKey: InterfaceKey, ID: id1},
						Tags:   tags1,
					},
					TagsToCheck: []TagToCheck{
						{Tag: Tag{"1", nil}, Tagged: Index{key: []Tagged{{ID: id1, Params: params1}}}},
						{Tag: Tag{"2", nil}, Tagged: Index{key: []Tagged{{ID: id1}}}},
						{Tag: Tag{"3", nil}, Tagged: Index{key: []Tagged{{ID: id1}}}},
						{Tag: Tag{"4", nil}, Tagged: Index{}},
						{Tag: Tag{"5", nil}, Tagged: Index{}},
						{Tag: Tag{"6", nil}, Tagged: Index{}},
					},
				},
				{
					TagsToChange: TagsToChange{
						Action: "add",
						ToTag:  joiner.Link{InterfaceKey: InterfaceKey, ID: id2},
						Tags:   tags2,
					},
					TagsToCheck: []TagToCheck{
						{Tag: Tag{"1", nil}, Tagged: Index{key: []Tagged{{ID: id1, Params: params1}}}},
						{Tag: Tag{"2", nil}, Tagged: Index{key: []Tagged{{ID: id1}}}},
						{Tag: Tag{"3", nil}, Tagged: Index{key: []Tagged{{ID: id1}, {ID: id2, Params: params2}}}},
						{Tag: Tag{"4", nil}, Tagged: Index{}},
						{Tag: Tag{"5", nil}, Tagged: Index{key: []Tagged{{ID: id2}}}},
						{Tag: Tag{"6", nil}, Tagged: Index{key: []Tagged{{ID: id2}}}},
					},
				},
				{
					TagsToChange: TagsToChange{
						Action: "replace",
						ToTag:  joiner.Link{InterfaceKey: InterfaceKey, ID: id1},
						Tags:   nil,
					},
					TagsToCheck: []TagToCheck{
						{Tag: Tag{"1", nil}, Tagged: Index{}},
						{Tag: Tag{"2", nil}, Tagged: Index{}},
						{Tag: Tag{"3", nil}, Tagged: Index{key: []Tagged{{ID: id2, Params: params2}}}},
						{Tag: Tag{"4", nil}, Tagged: Index{}},
						{Tag: Tag{"5", nil}, Tagged: Index{key: []Tagged{{ID: id2}}}},
						{Tag: Tag{"6", nil}, Tagged: Index{key: []Tagged{{ID: id2}}}},
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

		err := cleanerOp.Clean(nil, nil)
		require.NoError(t, err)

		for j, step := range tc.Steps {
			l.Infof("\tstep #%d", j)

			var err error
			switch step.Action {
			case "add":
				err = tc.Operator.AddTags(step.ToTag, step.Tags, nil)
			//case "remove":
			//	err = tc.Actor.RemoveTags(step.ID, step.ID, step.Tags, nil)
			case "replace":
				err = tc.Operator.ReplaceTags(step.ToTag, step.Tags, nil)
			case "tags":
				var tags []Tag
				tags, err = tc.Operator.ListTags(step.ToTag, nil)
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
				tagged, err := tc.Operator.IndexTagged(nil, tagToCheck.Tag.Label, nil)

				if tagToCheck.IsErrorExpected {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					require.Equal(t, tagToCheck.Tagged, tagged, "was checked: "+tagToCheck.Tag.Label)
				}
			}
		}
	}
}

//func prepareTest(t *testing.T, operator Actor, settaggerSteps []SettaggerStep) {
//	// ClearDatabase ------------------------------------------------------------------------------------
//
//	//err := goroutine.go.Clean()
//	//require.NoError(t, err, "what is the error on .Clean()?")
//
//	// test Settagger --------------------------------------------------------------------------------------
//
//	for j, sl := range settaggerSteps {
//		fmt.Println("." + strconv.Itoa(j))
//
//		if sl.ExpectedErr != nil {
//			_, err := operator.Settagger(sl.IS, sl.LinkedKey, sl.LinkedID, sl.Tags)
//			require.ErrStr(t, err, "where is an error on .Seltagger(%#v, %s, %s, %#v)?", sl.IS, sl.LinkedKey, sl.LinkedID, sl.Tags)
//			continue
//		}
//
//		if sl.ISBad != nil {
//			_, err := operator.Settagger(*sl.ISBad, sl.LinkedKey, sl.LinkedID, sl.Tags)
//			require.ErrStr(t, err, "where is an error on .Seltagger(%#v, %s, %s, %#v)?", *sl.ISBad, sl.LinkedKey, sl.LinkedID, sl.Tags)
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
//		prepareTest(t, tc.Actor, tc.SettaggerSteps)
//
//		// test QueryByObjectID --------------------------------------------------------------------------------------
//
//		if tc.ExpectedErr != nil {
//			_, err := tc.Actor.QueryByObjectID(tc.IS, tc.ObjectID)
//			require.ErrStr(t, err, "where is an error on .QueryByObjectID(%#v, %s)?", tc.IS, tc.ObjectID)
//			continue
//		}
//
//		if tc.ISBad != nil {
//			_, err := tc.Actor.QueryByObjectID(*tc.ISBad, tc.ObjectID)
//			require.ErrStr(t, err, "where is an error on .QueryByObjectID(%#v, %s)?", *tc.ISBad, tc.ObjectID)
//		}
//
//		linked, err := tc.Actor.QueryByObjectID(tc.IS, tc.ObjectID)
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
//		res[v.Tag] = v.CountTags
//	}
//
//	return res
//}
//
//func Hash(li Linked) string {
//	return li.ObjectID + " " + li.TypeKey + " " + li.LinkedID + " " + li.LinkedType + " " + li.Tag
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
