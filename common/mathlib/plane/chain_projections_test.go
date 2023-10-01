package plane

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/mathlib"
)

func TestProjectionsOnPolyChain(t *testing.T) {
	tests := []struct {
		name        string
		polyChain   PolyChain
		p           Point2
		distanceMax float64
		expected    []ProjectionOnPolyChainDirected
	}{
		{
			name:        "",
			polyChain:   PolyChain{{0, 2}, {1, 1}, {2, 2}}, // , {3, 1}
			p:           Point2{1, 0},
			distanceMax: 10,
			expected: []ProjectionOnPolyChainDirected{
				{Distance: 1, ProjectionOnPolyChain: ProjectionOnPolyChain{N: 1, Point2: Point2{1, 1}}}},
		},
		{
			name:        "",
			polyChain:   PolyChain{{0, 2}, {1, 1}, {2, 2}, {3, 1}}, //
			p:           Point2{1, 0},
			distanceMax: 10,
			expected: []ProjectionOnPolyChainDirected{
				{Distance: 1, ProjectionOnPolyChain: ProjectionOnPolyChain{N: 1, Point2: Point2{1, 1}}},
				{Distance: 3 * math.Sqrt(2) / 2,
					ProjectionOnPolyChain: ProjectionOnPolyChain{N: 2, Position: math.Sqrt(2) / 2, Point2: Point2{2.5, 1.5}}},
			},
		},
		{
			name:        "",
			polyChain:   PolyChain{{X: 567.758571909734, Y: 327.33667901650387}, {X: 588, Y: 381}}, //
			p:           Point2{X: 574, Y: 355},
			distanceMax: 3,
			expected:    []ProjectionOnPolyChainDirected{},
		},
		{
			name:        "",
			polyChain:   PolyChain{{X: 416.5, Y: 420}, {X: 422, Y: 413}, {X: 414.5, Y: 402}}, //
			p:           Point2{X: 427, Y: 412},
			distanceMax: 3,
			expected:    []ProjectionOnPolyChainDirected{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProjectionsOnPolyChain(tt.polyChain, tt.p, tt.distanceMax)
			CheckProjections(t, tt.expected, got)
		})
	}
}

// PolyChain{Point2{X:556, Y:355}, Point2{X:559, Y:355}, Point2{X:562, Y:355}, Point2{X:565, Y:355}, Point2{X:568,
//	Y:355}, Point2{X:571, Y:355}, Point2} / 3)

func CheckProjections(t *testing.T, expected, got []ProjectionOnPolyChainDirected) {
	require.Equalf(t, len(expected), len(got), "expected: %v, got: %#v", expected, got)
	for i, e := range expected {
		g := got[i]
		require.Equalf(t, e.N, g.N, "#d: %v vs %v", i, e, g)
		require.Truef(t, math.Abs(e.Position-g.Position) <= mathlib.Eps, "#d: %v vs %v", i, e, g)
		require.Truef(t, math.Abs(e.Distance-g.Distance) <= mathlib.Eps, "#d: %v vs %v", i, e, g)
		require.Truef(t, math.Abs(e.Angle-g.Angle) <= mathlib.Eps, "#d: %v vs %v", i, e, g)
		require.Truef(t, math.Abs(e.X-g.X) <= mathlib.Eps, "#d: %v vs %v", i, e, g)
		require.Truef(t, math.Abs(e.Y-g.Y) <= mathlib.Eps, "#d: %v vs %v", i, e, g)

	}
}
