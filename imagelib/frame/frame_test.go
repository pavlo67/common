package frame

import (
	"image"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/imagelib"
	"github.com/pavlo67/common/common/mathlib"
	"github.com/pavlo67/common/common/mathlib/plane"
)

func TestFrame_PolyChainToOuter(t *testing.T) {
	tests := []struct {
		name             string
		Position         plane.Position
		DPM              float64
		RectInner        image.Rectangle
		pChInner         plane.PolyChain
		pChOuterExpected plane.PolyChain
	}{
		{
			name:             "",
			Position:         plane.Position{Point2: plane.Point2{1, 1}, XToYAngle: math.Pi},
			DPM:              1,
			RectInner:        image.Rectangle{Max: image.Point{5, 5}},
			pChInner:         plane.PolyChain{{0, 0}},
			pChOuterExpected: plane.PolyChain{{3, -1}},
		},
		{
			name:             "",
			Position:         plane.Position{Point2: plane.Point2{1, 1}, XToYAngle: math.Pi},
			DPM:              2,
			RectInner:        image.Rectangle{Max: image.Point{5, 5}},
			pChInner:         plane.PolyChain{{0, 0}},
			pChOuterExpected: plane.PolyChain{{2, 0}},
		},
		{
			name:             "",
			Position:         plane.Position{Point2: plane.Point2{1, 1}, XToYAngle: math.Pi},
			DPM:              2,
			RectInner:        image.Rectangle{Max: image.Point{5, 5}},
			pChInner:         plane.PolyChain{{2, 2}},
			pChOuterExpected: plane.PolyChain{{1, 1}},
		},
		{
			name:             "",
			Position:         plane.Position{Point2: plane.Point2{1, 1}, XToYAngle: math.Pi},
			DPM:              2,
			RectInner:        image.Rectangle{Max: image.Point{5, 3}},
			pChInner:         plane.PolyChain{{0, 0}},
			pChOuterExpected: plane.PolyChain{{2, 0.5}},
		},
		{
			name:             "",
			Position:         plane.Position{Point2: plane.Point2{1, 1}, XToYAngle: math.Pi / 2},
			DPM:              2,
			RectInner:        image.Rectangle{Max: image.Point{5, 3}},
			pChInner:         plane.PolyChain{{0, 0}},
			pChOuterExpected: plane.PolyChain{{0.5, 0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr := Frame{
				LayerRGBA: LayerRGBA{
					RGBA:     image.RGBA{Rect: tt.RectInner},
					Settings: imagelib.Settings{DPM: tt.DPM},
				},
				Position: tt.Position,
			}
			pChOuter := fr.PointsToOuter(tt.pChInner...)

			for i, pOuter := range pChOuter {
				distance := pOuter.DistanceTo(tt.pChOuterExpected[i])
				if math.IsNaN(distance) || distance > mathlib.EPS {
					t.Errorf("PointsToOuter() = %v, pChOuterExpected %v", pChOuter, tt.pChOuterExpected)
				}

				pInner := fr.PointToInner(pOuter)
				distanceFinal := pInner.DistanceTo(tt.pChInner[i])
				if distanceFinal > mathlib.EPS {
					t.Errorf("pInner = %v, pOuter = %v, pOuterInner = %v", tt.pChInner[i], pOuter, pInner)
				}

			}
		})
	}
}

func TestFrame_MovingToOuter(t *testing.T) {
	tests := []struct {
		name                string
		layerRGBA           LayerRGBA
		position            plane.Position
		movingInner         plane.Point2
		movingOuterExpected plane.Point2
	}{
		{
			name: "1",
			layerRGBA: LayerRGBA{
				RGBA:     image.RGBA{Rect: image.Rectangle{Max: image.Point{100, 100}}},
				Settings: imagelib.Settings{DPM: 2},
			},
			movingInner:         plane.Point2{2, 2},
			movingOuterExpected: plane.Point2{1, -1},
		},
		{
			name: "2",
			layerRGBA: LayerRGBA{
				RGBA: image.RGBA{
					Rect: image.Rectangle{Max: image.Point{100, 100}},
				},
				Settings: imagelib.Settings{DPM: 2},
			},
			position:            plane.Position{Point2: plane.Point2{1, 1}, XToYAngle: math.Pi / 2},
			movingInner:         plane.Point2{2, 2},
			movingOuterExpected: plane.Point2{1, 1},
		},
		{
			name: "3",
			layerRGBA: LayerRGBA{
				RGBA: image.RGBA{
					Rect: image.Rectangle{Max: image.Point{100, 100}},
				},
				Settings: imagelib.Settings{DPM: 2},
			},
			position:            plane.Position{Point2: plane.Point2{10, 1}, XToYAngle: math.Pi},
			movingInner:         plane.Point2{2, 2},
			movingOuterExpected: plane.Point2{-1, 1},
		},
		{
			name: "4",
			layerRGBA: LayerRGBA{
				RGBA: image.RGBA{
					Rect: image.Rectangle{Max: image.Point{100, 100}},
				},
				Settings: imagelib.Settings{DPM: 2},
			},
			position:            plane.Position{Point2: plane.Point2{10, -10}, XToYAngle: math.Pi * 3 / 2},
			movingInner:         plane.Point2{2, 2},
			movingOuterExpected: plane.Point2{-1, -1},
		},
		{
			name: "5",
			layerRGBA: LayerRGBA{
				RGBA: image.RGBA{
					Rect: image.Rectangle{Max: image.Point{100, 100}},
				},
				Settings: imagelib.Settings{DPM: 2},
			},
			position:            plane.Position{Point2: plane.Point2{10, -10}, XToYAngle: math.Pi * 3 / 2},
			movingInner:         plane.Point2{0, 0},
			movingOuterExpected: plane.Point2{0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame := Frame{
				LayerRGBA: tt.layerRGBA,
				Position:  tt.position,
			}

			movingOuter := frame.MovingToOuter(tt.movingInner)
			if movingOuter.DistanceTo(tt.movingOuterExpected) > mathlib.EPS {
				t.Errorf("MovingToOuter() = %v, movingOuterExpected = %v, distance = %f",
					movingOuter, tt.movingOuterExpected, movingOuter.DistanceTo(tt.movingOuterExpected))
			}

			movingInner := frame.MovingToInner(movingOuter)
			if movingInner.DistanceTo(tt.movingInner) > mathlib.EPS {
				t.Errorf("MovingToInner() = %v, movingInner = %v, distance = %f",
					movingOuter, tt.movingOuterExpected, movingInner.DistanceTo(tt.movingInner))
			}

		})
	}
}

func TestFrame_ZeroMoving(t *testing.T) {
	frame := Frame{
		LayerRGBA: LayerRGBA{Settings: imagelib.Settings{DPM: 4.815940238327553}},
		Position:  plane.Position{XToYAngle: 2.8448274196093593},
	}

	movingOuter := plane.Point2{}
	movingInner := frame.MovingToInner(movingOuter)
	require.Truef(t, movingInner.Radius() < mathlib.EPS, "movingInner: %+v", movingInner)

	movingInner = plane.Point2{}
	movingOuter = frame.MovingToOuter(movingInner)
	require.Truef(t, movingOuter.Radius() < mathlib.EPS, "movingOuter: %+v", movingOuter)

}
