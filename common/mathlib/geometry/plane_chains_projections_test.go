package geometry

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProjectionsOnPolyChain(tt.polyChain, tt.p, tt.distanceMax)
			CheckProjections(t, tt.expected, got)
		})
	}
}

func CheckProjections(t *testing.T, expected, got []ProjectionOnPolyChainDirected) {
	require.Equalf(t, len(expected), len(got), "expected: %v, got: %v", expected, got)
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
