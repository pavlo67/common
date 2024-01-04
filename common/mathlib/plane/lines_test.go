package plane

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type IntersectionTestCase struct {
	S1, S2       Segment
	Intersection *Point2
}

func TestLineSegmentsIntersection(t *testing.T) {
	testCases := []IntersectionTestCase{
		{Segment{Point2{1, 0}, Point2{1, 1}},
			Segment{Point2{0, 1.5}, Point2{2, 1.5}}, nil},
		{Segment{Point2{1, 2}, Point2{1, 1}},
			Segment{Point2{1.25, 1.5}, Point2{2, 1.5}}, nil},

		{Segment{Point2{1, 1}, Point2{1, 2}},
			Segment{Point2{3, 0}, Point2{2, 1}}, nil},

		{Segment{Point2{2, 2}, Point2{-1, -1}},
			Segment{Point2{-3, 1}, Point2{3, -1}}, &Point2{0, 0}},
		{Segment{Point2{2, 2}, Point2{3, 3}},
			Segment{Point2{-3, 1}, Point2{3, -1}}, nil},
		{Segment{Point2{3, 3}, Point2{2, 2}},
			Segment{Point2{-3, 1}, Point2{3, -1}}, nil},
		{Segment{Point2{0, 0}, Point2{1, 1}},
			Segment{Point2{1, 1}, Point2{2, 2}}, nil},

		{Segment{Point2{0, 0}, Point2{0, 1}},
			Segment{Point2{0, 0.1}, Point2{0, 2}}, &Point2{0, 0.1}},

		{Segment{Point2{0, 0}, Point2{1, 0}},
			Segment{Point2{-2, 0}, Point2{-1, 0}}, nil},
		{Segment{Point2{1, 1}, Point2{1, 2}},
			Segment{Point2{1, 1}, Point2{2, 1}}, &Point2{1, 1}},
		{Segment{Point2{1, 1}, Point2{1, 2}},
			Segment{Point2{2, 1}, Point2{2, 2}}, nil},
		{Segment{Point2{1, 1}, Point2{1, 2}},
			Segment{Point2{0, 1.5}, Point2{2, 1.5}}, &Point2{1, 1.5}},
	}

	for i, testCase := range testCases {
		t.Log(i)
		// t.Logf("%#v", testCase)
		intersection := SegmentsIntersection(testCase.S1, testCase.S2)
		if testCase.Intersection == nil {
			require.Nil(t, intersection)
		} else {
			require.NotNil(t, intersection)
			require.Equal(t, *testCase.Intersection, *intersection)
		}
	}
}

func TestLinesIntersection(t *testing.T) {
	testCases := []IntersectionTestCase{
		{Segment{Point2{1, 1}, Point2{1, 2}},
			Segment{Point2{1, 1}, Point2{2, 1}}, &Point2{1, 1}},

		{Segment{Point2{1, 1}, Point2{1, 2}},
			Segment{Point2{2, 1}, Point2{2, 2}}, nil},

		{Segment{Point2{1, 1}, Point2{1, 2}},
			Segment{Point2{0, 1.5}, Point2{2, 1.5}}, &Point2{1, 1.5}},

		{Segment{Point2{1, 2}, Point2{1, 1}},
			Segment{Point2{1.25, 1.5}, Point2{2, 1.5}}, &Point2{1, 1.5}},

		{Segment{Point2{1, 1}, Point2{1, 2}},
			Segment{Point2{3, 0}, Point2{2, 1}}, &Point2{1, 2}},

		{Segment{Point2{2, 2}, Point2{-1, -1}},
			Segment{Point2{-3, 1}, Point2{3, -1}}, &Point2{0, 0}},

		{Segment{Point2{2, 2}, Point2{3, 3}},
			Segment{Point2{-3, 1}, Point2{3, -1}}, &Point2{0, 0}},

		{Segment{Point2{3, 3}, Point2{2, 2}},
			Segment{Point2{-3, 1}, Point2{3, -1}}, &Point2{0, 0}},

		{Segment{Point2{0, 0}, Point2{1, 1}},
			Segment{Point2{1, 1}, Point2{2, 2}}, nil},

		{Segment{Point2{0, 0}, Point2{0, 1}},
			Segment{Point2{0, 0.1}, Point2{0, 2}}, nil},

		{Segment{Point2{0, 0}, Point2{1, 0}},
			Segment{Point2{-2, 0}, Point2{-1, 0}}, nil},
	}

	for i, testCase := range testCases {
		t.Log(i)
		intersection := testCase.S1.LinesIntersection(testCase.S2)
		if testCase.Intersection == nil {
			require.Nil(t, intersection)
		} else {
			require.NotNil(t, intersection)
			require.Equal(t, *testCase.Intersection, *intersection)
		}
	}
}
