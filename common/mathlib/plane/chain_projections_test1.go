package plane

import (
	"reflect"
	"testing"
)

func TestCutWithProjections(t *testing.T) {
	tests := []struct {
		name string
		pCh  PolyChain
		pr0  ProjectionOnPolyChain
		pr1  ProjectionOnPolyChain
		want PolyChain
	}{
		{
			name: "",
			pCh:  PolyChain{{X: 2.3, Y: 458}, {X: 91, Y: 427.5}},
			pr0:  ProjectionOnPolyChain{N: 1, Position: 0, Point2: Point2{X: 91, Y: 427.5}},
			pr1:  ProjectionOnPolyChain{N: 0, Position: 9.7, Point2: Point2{X: 11.5, Y: 455}},
			want: PolyChain{{X: 91, Y: 427.5}, {11.5, 455}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CutWithProjections(tt.pCh, tt.pr0, tt.pr1); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CutWithProjections() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDistanceToSegment(t *testing.T) {

	tests := []struct {
		name    string
		p       Point2
		segment Segment
		want    float64
	}{
		{
			name:    "",
			p:       Point2{203.77788799006555, 564.5811856102348},
			segment: Segment{{230.79193808962214, 591.2694925293472}, {145.99999999999994, 507.50000000000006}},
			want:    0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.p.DistanceToSegment(tt.segment)
			if got != tt.want {
				t.Errorf("DistanceToSegment() got = %v, want %v", got, tt.want)
			}
		})
	}
}
