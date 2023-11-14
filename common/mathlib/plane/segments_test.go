package plane

import (
	"reflect"
	"testing"
)

func TestSegmentGoOutCircle(t *testing.T) {
	tests := []struct {
		name string
		s    Segment
		p    Point2
		r    float64
		want *Point2
	}{
		{
			name: "",
			s:    Segment{{3, 0}, {3, 10}},
			p:    Point2{0, 0},
			r:    5,
			want: &Point2{3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SegmentGoOutCircle(tt.s, tt.p, tt.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SegmentGoOutCircle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSegmentsIntersection(t *testing.T) {
	tests := []struct {
		name string
		s    Segment
		s1   Segment
		want *Point2
	}{
		{"", Segment{{1, 0}, {1, 1}}, Segment{{0, 1.5}, {2, 1.5}}, nil},
		{"", Segment{{1, 2}, {1, 1}}, Segment{{1.25, 1.5}, {2, 1.5}}, nil},
		{"", Segment{{1, 1}, {1, 2}}, Segment{{3, 0}, {2, 1}}, nil},
		{"", Segment{{2, 2}, {-1, -1}}, Segment{{-3, 1}, {3, -1}}, &Point2{0, 0}},
		{"", Segment{{2, 2}, {3, 3}}, Segment{{-3, 1}, {3, -1}}, nil},
		{"", Segment{{3, 3}, {2, 2}}, Segment{{-3, 1}, {3, -1}}, nil},
		{"", Segment{{0, 0}, {1, 1}}, Segment{{1, 1}, {2, 2}}, nil},
		{"", Segment{{0, 0}, {0, 1}}, Segment{{0, 0.1}, {0, 2}}, &Point2{0, 0.1}},
		{"", Segment{{0, 0}, {1, 0}}, Segment{{-2, 0}, {-1, 0}}, nil},
		{"", Segment{{1, 1}, {1, 2}}, Segment{{1, 1}, {2, 1}}, &Point2{1, 1}},
		{"", Segment{{1, 1}, {1, 2}}, Segment{{2, 1}, {2, 2}}, nil},
		{"", Segment{{1, 1}, {1, 2}}, Segment{{0, 1.5}, {2, 1.5}}, &Point2{1, 1.5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPCross := SegmentsIntersection(tt.s, tt.s1)
			if !reflect.DeepEqual(gotPCross, tt.want) {
				t.Errorf("SegmentsIntersection() gotPCross = %v, want %v", gotPCross, tt.want)
			}
		})
	}
}
