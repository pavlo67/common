package links

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"sort"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/components/auth"
)

// TODO: test .Selector

// TODO: wtf???

const TypeTag = "tag"
const TypeOriginal = "original"
const TypeAuthor = "author"

// SetLinks ------------------------------------------------------------------------------------------------

type SetLinksStep struct {
	IS                 common.ID
	ISBad              *common.ID
	LinkedType         string
	LinkedID           string
	Links              []Item
	ExpectedErr        error
	ExpectedLinkedInfo []LinkedInfo
}

// QueryTags -----------------------------------------------------------------------------------------------

type QueryTagsStep struct {
	IS              common.ID
	ISBad           *common.ID
	Selector        *libs.Term
	ExpectedErr     error
	ExpectedTagInfo []TagInfo
}

type QueryTagsTestCase struct {
	Operator Operator

	SetLinksSteps  []SetLinksStep
	QueryTagsSteps []QueryTagsStep
}

// QueryTagsByOwner ---------------------------------------------------------------------------------------

type QueryTagsByOwnerStep struct {
	IS              common.ID
	ISBad           *common.ID
	ROwner          common.ID
	ExpectedErr     error
	ExpectedTagInfo []TagInfo
}

type QueryTagsByOwnerTestCase struct {
	Operator Operator

	SetLinksSteps         []SetLinksStep
	QueryTagsByOwnerSteps []QueryTagsByOwnerStep
}

// QueryByObjectID ----------------------------------------------------------------------------------------

type QueryByObjectIDTestCase struct {
	Operator      Operator
	SetLinksSteps []SetLinksStep

	IS             common.ID
	ISBad          *common.ID
	ObjectID       string
	ExpectedErr    error
	ExpectedLinked []Linked
}

// QueryByTag ---------------------------------------------------------------------------------------------

type QueryByTagTestCase struct {
	Operator      Operator
	SetLinksSteps []SetLinksStep

	IS             common.ID
	ISBad          *common.ID
	Tag            string
	ExpectedErr    error
	ExpectedLinked []Linked
}

// Query --------------------------------------------------------------------------------------------------

type QueryStep struct {
	IS             common.ID
	ISBad          *common.ID
	Selector       *libs.Term
	ExpectedErr    error
	ExpectedLinked []Linked
}

type QueryTestCase struct {
	Operator Operator

	SetLinksSteps []SetLinksStep
	QuerySteps    []QueryStep
}

// --------------------------------------------------------------------------------------------------------

const linkedType1 = "linkedType1"
const linkedType2 = "linkedType2"

const linkedID1 = "1"
const linkedID2 = "2"

const linkName1 = "linkName1"
const linkName2 = "linkName2"
const linkName3 = "linkName3"

const objectID1 = "1"
const objectID2 = "2"
const objectID3 = "3"
const objectID4 = "4"

var userISNil common.ID

func setLinkSteps(userISToSet, userISToSetAnother common.ID) []SetLinksStep {
	is := userISToSet
	isAnother := userISToSetAnother

	return []SetLinksStep{
		// .0 added public tags.comp
		{
			IS:         userISToSet,
			ISBad:      &userISNil,
			LinkedType: linkedType1,
			LinkedID:   linkedID1,
			Links: []Item{
				{Type: TypeTag, Name: linkName1, To: objectID1, RView: auth.Anyone, ROwner: is},
				{Type: TypeTag, Name: linkName2, To: objectID2, RView: auth.Anyone, ROwner: is},

				{Type: "", Name: linkName3, To: objectID2, RView: auth.Anyone, ROwner: is},
				// set and counted
			},
			ExpectedErr: nil,
			ExpectedLinkedInfo: []LinkedInfo{
				{ObjectID: objectID1, CountLinked: 1},
				{ObjectID: objectID2, CountLinked: 2},
			},
		},

		// .1 added (but not counted) private tags.comp
		{
			IS:         userISToSet,
			ISBad:      &userISNil,
			LinkedType: linkedType1,
			LinkedID:   linkedID2,
			Links: []Item{
				{Type: TypeTag, Name: linkName1, To: objectID1, RView: auth.Anyone, ROwner: is},

				{Type: TypeTag, Name: linkName1, To: objectID3, RView: is, ROwner: is},
				// ignored in linked count, queried by QueryBy...
				{Type: TypeTag, Name: linkName2, To: objectID4, RView: is, ROwner: is},
				// ignored in linked count, queried by QueryBy...
			},
			ExpectedLinkedInfo: []LinkedInfo{
				{ObjectID: objectID1, CountLinked: 2},
			},
		},

		// .2 added tags.comp from another user (except "alien" one)
		{
			IS:         userISToSetAnother,
			ISBad:      &userISNil,
			LinkedType: linkedType2,
			LinkedID:   linkedID1,
			Links: []Item{
				{Type: TypeTag, Name: linkName1, To: objectID1, RView: auth.Anyone, ROwner: isAnother},
				{Type: TypeTag, Name: linkName1, To: objectID2, RView: auth.Anyone, ROwner: isAnother},

				{Type: TypeTag, Name: linkName2, To: objectID4, RView: auth.Anyone, ROwner: isAnother},
				{Type: TypeTag, Name: linkName2, To: objectID2, RView: auth.Anyone, ROwner: is},
				// not set
			},
			ExpectedLinkedInfo: []LinkedInfo{
				{ObjectID: objectID1, CountLinked: 3},
				{ObjectID: objectID2, CountLinked: 3},
				{ObjectID: objectID4, CountLinked: 1},
			},
		},
	}
}

func QueryTagsByOwnerTestCases(linksOp Operator, userISToSet, userISToSetAnother common.ID) []QueryTagsByOwnerTestCase {
	return []QueryTagsByOwnerTestCase{
		// 0 all ok for userISToSet
		{
			Operator:      linksOp,
			SetLinksSteps: setLinkSteps(userISToSet, userISToSetAnother),
			QueryTagsByOwnerSteps: []QueryTagsByOwnerStep{
				{
					IS:     userISToSet,
					ROwner: userISToSet,
					ExpectedTagInfo: []TagInfo{
						{Tag: "linkName1", Count: 3},
						{Tag: "linkName2", Count: 2},
						{Tag: "linkName3", Count: 1},
					},
				},
			},
		},

		// 1 all ok for userISToSetAnother
		{
			Operator:      linksOp,
			SetLinksSteps: setLinkSteps(userISToSet, userISToSetAnother),
			QueryTagsByOwnerSteps: []QueryTagsByOwnerStep{
				{
					IS:     userISToSetAnother,
					ROwner: userISToSetAnother,
					ExpectedTagInfo: []TagInfo{
						{Tag: "linkName1", Count: 2},
						{Tag: "linkName2", Count: 1},
					},
				},
			},
		},

		// 2 all ok for userISToSetAnother (viewing userISToSet)
		{
			Operator:      linksOp,
			SetLinksSteps: setLinkSteps(userISToSet, userISToSetAnother),
			QueryTagsByOwnerSteps: []QueryTagsByOwnerStep{
				{
					IS:     userISToSetAnother,
					ROwner: userISToSet,
					ExpectedTagInfo: []TagInfo{
						{Tag: "linkName1", Count: 2},
						{Tag: "linkName2", Count: 1},
						{Tag: "linkName3", Count: 1},
					},
				},
			},
		},

		// 3 all ok for userISNil
		{
			Operator:      linksOp,
			SetLinksSteps: setLinkSteps(userISToSet, userISToSetAnother),
			QueryTagsByOwnerSteps: []QueryTagsByOwnerStep{
				{
					IS:              userISNil,
					ROwner:          userISNil,
					ExpectedTagInfo: []TagInfo{},
				},
			},
		},

		// 4 all ok for userISNil (viewing userISToSet)
		{
			Operator:      linksOp,
			SetLinksSteps: setLinkSteps(userISToSet, userISToSetAnother),
			QueryTagsByOwnerSteps: []QueryTagsByOwnerStep{
				{
					IS:     userISNil,
					ROwner: userISToSet,
					ExpectedTagInfo: []TagInfo{
						{Tag: "linkName1", Count: 2},
						{Tag: "linkName2", Count: 1},
						{Tag: "linkName3", Count: 1},
					},
				},
			},
		},
	}
}

func QueryTagsByOwnerTest(t *testing.T, testCases []QueryTagsByOwnerTestCase) {
	if err := os.Setenv("ENV", "test"); err != nil {
		t.Fatal("No test environment!!!")
	}

	for i, tc := range testCases {
		fmt.Println("QueryTagsByOwnerTest: ", i)

		prepareTest(t, tc.Operator, tc.SetLinksSteps)

		// test QueryTagsByOwner -------------------------------------------------------------------------------

		for j, qt := range tc.QueryTagsByOwnerSteps {

			fmt.Println("." + strconv.Itoa(j))

			if qt.ExpectedErr != nil {
				_, err := tc.Operator.QueryTagsByOwner(qt.IS, qt.ROwner)
				require.Error(t, err, "where is an error on .QueryTags(%#v, %#s)?", qt.IS, qt.ROwner)
				continue
			}

			if qt.ISBad != nil {
				_, err := tc.Operator.QueryTagsByOwner(*qt.ISBad, (*qt.ISBad))
				require.Error(t, err, "where is an error on .QueryTags(%#v, %#s)?", *qt.ISBad, *qt.ISBad)
			}

			tagInfo, err := tc.Operator.QueryTagsByOwner(qt.IS, qt.ROwner)
			require.NoError(t, err, "what is an error on .QueryTags(%#v, %s)?", qt.IS, qt.ROwner)
			require.Equal(t, len(qt.ExpectedTagInfo), len(tagInfo), "len(tc.ExpectedTagInfo = %#v) != len(tagInfo = %#v)", qt.ExpectedTagInfo, tagInfo)

			exp := mapTagInfo(qt.ExpectedTagInfo)
			act := mapTagInfo(tagInfo)

			for tag, count := range exp {
				require.Equal(t, count, act[tag], "tagInfo[%s] isn't correct (%#v)", tag, tagInfo)
			}
		}

	}
}

func QueryTagsTestCases(linksOp Operator, userISToSet, userISToSetAnother common.ID) []QueryTagsTestCase {
	return []QueryTagsTestCase{
		// 0 all ok for userISToSet
		{
			Operator:      linksOp,
			SetLinksSteps: setLinkSteps(userISToSet, userISToSetAnother),
			QueryTagsSteps: []QueryTagsStep{
				{
					IS:       userISToSet,
					Selector: nil,
					ExpectedTagInfo: []TagInfo{
						{Tag: "linkName1", Count: 5},
						{Tag: "linkName2", Count: 3},
						{Tag: "linkName3", Count: 1},
					},
				},
			},
		},

		// 1 all ok for userISToSetAnother
		{
			Operator:      linksOp,
			SetLinksSteps: setLinkSteps(userISToSet, userISToSetAnother),
			QueryTagsSteps: []QueryTagsStep{
				{
					IS:       userISToSetAnother,
					Selector: nil,
					ExpectedTagInfo: []TagInfo{
						{Tag: "linkName1", Count: 4},
						{Tag: "linkName2", Count: 2},
						{Tag: "linkName3", Count: 1},
					},
				},
			},
		},

		// 2 all ok for userISNil
		{
			Operator:      linksOp,
			SetLinksSteps: setLinkSteps(userISToSet, userISToSetAnother),
			QueryTagsSteps: []QueryTagsStep{
				{
					IS:       userISNil,
					Selector: nil,
					ExpectedTagInfo: []TagInfo{
						{Tag: "linkName1", Count: 4},
						{Tag: "linkName2", Count: 2},
						{Tag: "linkName3", Count: 1},
					},
				},
			},
		},
	}
}

func QueryByTagTestCases(linksOp Operator, userISToSet, userISToSetAnother common.ID) []QueryByTagTestCase {
	return []QueryByTagTestCase{
		// 0 all ok for userISToSet
		{
			Operator:      linksOp,
			SetLinksSteps: setLinkSteps(userISToSet, userISToSetAnother),
			IS:            userISToSet,
			Tag:           linkName1,
			ExpectedLinked: []Linked{
				{LinkedType: linkedType1, LinkedID: linkedID1, Type: TypeTag, Tag: linkName1, ObjectID: objectID1},
				{LinkedType: linkedType1, LinkedID: linkedID2, Type: TypeTag, Tag: linkName1, ObjectID: objectID1},
				{LinkedType: linkedType1, LinkedID: linkedID2, Type: TypeTag, Tag: linkName1, ObjectID: objectID3},

				{LinkedType: linkedType2, LinkedID: linkedID1, Type: TypeTag, Tag: linkName1, ObjectID: objectID1},
				{LinkedType: linkedType2, LinkedID: linkedID1, Type: TypeTag, Tag: linkName1, ObjectID: objectID2},
			},
		},

		// 1 all ok for userISToSetAnother
		{
			Operator:      linksOp,
			SetLinksSteps: setLinkSteps(userISToSet, userISToSetAnother),
			IS:            userISToSet,
			Tag:           linkName2,
			ExpectedLinked: []Linked{
				{LinkedType: linkedType1, LinkedID: linkedID1, Type: TypeTag, Tag: linkName2, ObjectID: objectID2},
				{LinkedType: linkedType1, LinkedID: linkedID2, Type: TypeTag, Tag: linkName2, ObjectID: objectID4},

				{LinkedType: linkedType2, LinkedID: linkedID1, Type: TypeTag, Tag: linkName2, ObjectID: objectID4},
			},
		},

		// 2 all ok for userISNil
		{
			Operator:      linksOp,
			SetLinksSteps: setLinkSteps(userISToSet, userISToSetAnother),
			IS:            userISNil,
			Tag:           linkName3,
			ExpectedLinked: []Linked{
				{LinkedType: linkedType1, LinkedID: linkedID1, Type: "", Tag: linkName3, ObjectID: objectID2},
			},
		},
	}
}

func QueryByObjectIDTestCases(linksOp Operator, userISToSet, userISToSetAnother common.ID) []QueryByObjectIDTestCase {
	return []QueryByObjectIDTestCase{
		// 0 all ok for userISToSet
		{
			Operator:      linksOp,
			SetLinksSteps: setLinkSteps(userISToSet, userISToSetAnother),
			IS:            userISToSet,
			ObjectID:      objectID1,
			ExpectedLinked: []Linked{
				{LinkedType: linkedType1, LinkedID: linkedID1, Type: TypeTag, Tag: linkName1, ObjectID: objectID1},
				{LinkedType: linkedType1, LinkedID: linkedID2, Type: TypeTag, Tag: linkName1, ObjectID: objectID1},
				{LinkedType: linkedType2, LinkedID: linkedID1, Type: TypeTag, Tag: linkName1, ObjectID: objectID1},
			},
		},

		// 1 all ok for userISToSetAnother
		{
			Operator:      linksOp,
			SetLinksSteps: setLinkSteps(userISToSet, userISToSetAnother),
			IS:            userISToSet,
			ObjectID:      objectID2,
			ExpectedLinked: []Linked{
				{LinkedType: linkedType1, LinkedID: linkedID1, Type: TypeTag, Tag: linkName2, ObjectID: objectID2},
				{LinkedType: linkedType1, LinkedID: linkedID1, Type: "", Tag: linkName3, ObjectID: objectID2},
				{LinkedType: linkedType2, LinkedID: linkedID1, Type: TypeTag, Tag: linkName1, ObjectID: objectID2},
			},
		},

		// 2 all ok for userISNil
		{
			Operator:      linksOp,
			SetLinksSteps: setLinkSteps(userISToSet, userISToSetAnother),
			IS:            userISNil,
			ObjectID:      objectID4,
			ExpectedLinked: []Linked{
				{LinkedType: linkedType2, LinkedID: linkedID1, Type: TypeTag, Tag: linkName2, ObjectID: objectID4},

				{LinkedType: linkedType1, LinkedID: linkedID2, Type: TypeTag, Tag: linkName2, ObjectID: objectID4},
				// this is a private link though...
			},
		},
	}
}

func QueryTestCases(linksOp Operator, userISToSet, userISToSetAnother common.ID) []QueryTestCase {
	return []QueryTestCase{
		{
			Operator:      linksOp,
			SetLinksSteps: setLinkSteps(userISToSet, userISToSetAnother),
			QuerySteps:    []QueryStep{

				// 0 all ok for userISToSet
				//{
				//	IS:       userISToSet,
				//	Selector: selector.FieldEqual(FieldTag, linkName1),
				//	ExpectedLinked: []Linked{
				//		{LinkedType: linkedType1, LinkedID: linkedID1, Type: TypeTag, Tag: linkName1, ObjectID: objectID1},
				//		{LinkedType: linkedType1, LinkedID: linkedID2, Type: TypeTag, Tag: linkName1, ObjectID: objectID1},
				//		{LinkedType: linkedType2, LinkedID: linkedID1, Type: TypeTag, Tag: linkName1, ObjectID: objectID1},
				//		{LinkedType: linkedType2, LinkedID: linkedID1, Type: TypeTag, Tag: linkName1, ObjectID: objectID2},
				//		{LinkedType: linkedType1, LinkedID: linkedID2, Type: TypeTag, Tag: linkName1, ObjectID: objectID3},
				//	},
				//},
				//
				//// 1 all ok for userISToSetAnother
				//{
				//	IS:       userISToSet,
				//	Selector: selector.FieldEqual(FieldTag, linkName2),
				//	ExpectedLinked: []Linked{
				//		{LinkedType: linkedType1, LinkedID: linkedID1, Type: TypeTag, Tag: linkName2, ObjectID: objectID2},
				//		{LinkedType: linkedType1, LinkedID: linkedID2, Type: TypeTag, Tag: linkName2, ObjectID: objectID4},
				//		{LinkedType: linkedType2, LinkedID: linkedID1, Type: TypeTag, Tag: linkName2, ObjectID: objectID4},
				//	},
				//},
				//
				//// 2 all ok for userISNil
				//{
				//	IS:       userISNil,
				//	Selector: selector.FieldEqual(FieldTag, linkName3),
				//	ExpectedLinked: []Linked{
				//		{LinkedType: linkedType1, LinkedID: linkedID1, Type: "", Tag: linkName3, ObjectID: objectID2},
				//	},
				//},
			},
		},
	}
}

func mapLinkedInfo(linkedInfo []LinkedInfo) map[string]uint {
	res := map[string]uint{}
	for _, v := range linkedInfo {
		res[v.ObjectID] = v.CountLinked
	}

	return res
}

func mapTagInfo(tagInfo []TagInfo) map[string]uint64 {
	res := map[string]uint64{}
	for _, v := range tagInfo {
		res[v.Tag] = v.Count
	}

	return res
}

func Hash(li Linked) string {
	return li.ObjectID + " " + li.Type + " " + li.LinkedID + " " + li.LinkedType + " " + li.Tag
}

type byLinked []Linked

func (s byLinked) Len() int {
	return len(s)
}
func (s byLinked) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byLinked) Less(i, j int) bool {
	return Hash(s[i]) < Hash(s[j])
}

func prepareTest(t *testing.T, operator Operator, setLinksSteps []SetLinksStep) {
	// ClearDatabase ------------------------------------------------------------------------------------

	//err := operator.go.Clean()
	//require.NoError(t, err, "what is the error on .Clean()?")

	// test SetLinks --------------------------------------------------------------------------------------

	for j, sl := range setLinksSteps {
		fmt.Println("." + strconv.Itoa(j))

		if sl.ExpectedErr != nil {
			_, err := operator.SetLinks(sl.IS, sl.LinkedType, sl.LinkedID, sl.Links)
			require.Error(t, err, "where is an error on .SelLinks(%#v, %s, %s, %#v)?", sl.IS, sl.LinkedType, sl.LinkedID, sl.Links)
			continue
		}

		if sl.ISBad != nil {
			_, err := operator.SetLinks(*sl.ISBad, sl.LinkedType, sl.LinkedID, sl.Links)
			require.Error(t, err, "where is an error on .SelLinks(%#v, %s, %s, %#v)?", *sl.ISBad, sl.LinkedType, sl.LinkedID, sl.Links)
		}

		res, err := operator.SetLinks(sl.IS, sl.LinkedType, sl.LinkedID, sl.Links)
		require.NoError(t, err, "what is an error on .SelLinks(%#v, %s, %s, %#v)?", sl.IS, sl.LinkedType, sl.LinkedID, sl.Links)
		require.Equal(t, len(sl.ExpectedLinkedInfo), len(res), "len(sl.ExpectedLinkedInfo = %#v) != len(res = %#v)", sl.ExpectedLinkedInfo, res)

		exp := mapLinkedInfo(sl.ExpectedLinkedInfo)
		act := mapLinkedInfo(res)
		for objectID, countLinked := range exp {
			require.Equal(t, countLinked, act[objectID], "linkedInfo[%s] isn't correct (%#v)", objectID, res)
		}

	}

}

func QueryTagsTest(t *testing.T, testCases []QueryTagsTestCase) {
	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	for i, tc := range testCases {
		fmt.Println("QueryTagsTest: ", i)

		prepareTest(t, tc.Operator, tc.SetLinksSteps)

		// test QueryTags --------------------------------------------------------------------------------------

		for j, qt := range tc.QueryTagsSteps {
			fmt.Println("." + strconv.Itoa(j))

			if qt.ExpectedErr != nil {
				_, err := tc.Operator.QueryTags(qt.IS, qt.Selector)
				require.Error(t, err, "where is an error on .QueryTags(%#v, %#v)?", qt.IS, qt.Selector)
				continue
			}

			if qt.ISBad != nil {
				_, err := tc.Operator.QueryTags(*qt.ISBad, qt.Selector)
				require.Error(t, err, "where is an error on .QueryTags(%#v, %#v)?", *qt.ISBad, qt.Selector)
			}

			tagInfo, err := tc.Operator.QueryTags(qt.IS, qt.Selector)
			require.NoError(t, err, "what is an error on .QueryTags(%#v, %s)?", qt.IS, qt.Selector)
			require.Equal(t, len(qt.ExpectedTagInfo), len(tagInfo), "len(tc.ExpectedTagInfo = %#v) != len(tagInfo = %#v)", qt.ExpectedTagInfo, tagInfo)

			exp := mapTagInfo(qt.ExpectedTagInfo)
			act := mapTagInfo(tagInfo)

			for tag, count := range exp {
				require.Equal(t, count, act[tag], "tagInfo[%s] isn't correct (%#v)", tag, tagInfo)
			}
		}

	}
}

func QueryByObjectIDTest(t *testing.T, testCases []QueryByObjectIDTestCase) {
	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	for i, tc := range testCases {
		fmt.Println("QueryByObjectIDTest: ", i)

		prepareTest(t, tc.Operator, tc.SetLinksSteps)

		// test QueryByObjectID --------------------------------------------------------------------------------------

		if tc.ExpectedErr != nil {
			_, err := tc.Operator.QueryByObjectID(tc.IS, tc.ObjectID)
			require.Error(t, err, "where is an error on .QueryByObjectID(%#v, %s)?", tc.IS, tc.ObjectID)
			continue
		}

		if tc.ISBad != nil {
			_, err := tc.Operator.QueryByObjectID(*tc.ISBad, tc.ObjectID)
			require.Error(t, err, "where is an error on .QueryByObjectID(%#v, %s)?", *tc.ISBad, tc.ObjectID)
		}

		linked, err := tc.Operator.QueryByObjectID(tc.IS, tc.ObjectID)
		require.NoError(t, err, "what is an error on .QueryByObjectID(%#v, %s)?", tc.IS, tc.ObjectID)
		require.Equal(t, len(tc.ExpectedLinked), len(linked), "len(tc.ExpectedLinked = %#v) != len(linked = %#v)", tc.ExpectedLinked, linked)

		sort.Sort(byLinked(tc.ExpectedLinked))
		sort.Sort(byLinked(linked))

		for i, li := range tc.ExpectedLinked {
			require.Equal(t, Hash(li), Hash(linked[i]), "linked[%d] isn't correct", i)
		}

	}
}

func QueryByTagTest(t *testing.T, testCases []QueryByTagTestCase) {
	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	for i, tc := range testCases {
		fmt.Println("QueryByTagTest: ", i)

		prepareTest(t, tc.Operator, tc.SetLinksSteps)

		// test QueryByTag --------------------------------------------------------------------------------------

		if tc.ExpectedErr != nil {
			_, err := tc.Operator.QueryByTag(tc.IS, tc.Tag)
			require.Error(t, err, "where is an error on .QueryByTag(%#v, %s)?", tc.IS, tc.Tag)
			continue
		}

		if tc.ISBad != nil {
			_, err := tc.Operator.QueryByTag(*tc.ISBad, tc.Tag)
			require.Error(t, err, "where is an error on .QueryByTag(%#v, %s)?", *tc.ISBad, tc.Tag)
		}

		linked, err := tc.Operator.QueryByTag(tc.IS, tc.Tag)
		require.NoError(t, err, "what is an error on .QueryByTag(%#v, %s)?", tc.IS, tc.Tag)
		require.Equal(t, len(tc.ExpectedLinked), len(linked), "len(tc.ExpectedLinked = %#v) != len(linked = %#v)", tc.ExpectedLinked, linked)

		sort.Sort(byLinked(tc.ExpectedLinked))
		sort.Sort(byLinked(linked))

		for i, li := range tc.ExpectedLinked {
			require.Equal(t, Hash(li), Hash(linked[i]), "linked[%d] isn't correct", i)
		}

	}
}

func QueryTest(t *testing.T, testCases []QueryTestCase) {
	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	for i, tc := range testCases {
		fmt.Println("QueryTest: ", i)

		prepareTest(t, tc.Operator, tc.SetLinksSteps)

		for j, qt := range tc.QuerySteps {
			fmt.Println("." + strconv.Itoa(j))

			if qt.ExpectedErr != nil {
				_, err := tc.Operator.Query(qt.IS, qt.Selector)
				require.Error(t, err, "where is an error on .Query(%#v, %#v)?", qt.IS, qt.Selector)
				continue
			}

			if qt.ISBad != nil {
				_, err := tc.Operator.Query(*qt.ISBad, qt.Selector)
				require.Error(t, err, "where is an error on .QueryTags(%#v, %#v)?", *qt.ISBad, qt.Selector)
			}

			linked, err := tc.Operator.Query(qt.IS, qt.Selector)
			require.NoError(t, err, "what is an error on .Query(%#v, %#v)?", qt.IS, qt.Selector)
			require.Equal(t, len(qt.ExpectedLinked), len(linked), "len(qt.ExpectedLinked = %#v) != len(linked = %#v)", qt.ExpectedLinked, linked)

			sort.Sort(byLinked(qt.ExpectedLinked))
			sort.Sort(byLinked(linked))

			for i, li := range qt.ExpectedLinked {
				require.Equal(t, Hash(li), Hash(linked[i]), "linked[%d] isn't correct", i)
			}
		}
	}
}
