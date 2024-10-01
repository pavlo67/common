package imagelib

import (
	"math"
	"testing"

	"github.com/pavlo67/common/common/mathlib"
)

func TestAngle(t *testing.T) {
	tests := []struct {
		name     string
		v        float64
		vMax     float64
		angleMax float64
		want     float64
	}{
		{
			name:     "",
			v:        0.5 / math.Cos(math.Pi/6),
			vMax:     1,
			angleMax: math.Pi / 4,
			want:     math.Pi / 6,
		},
		{
			name:     "",
			v:        1,
			vMax:     0.5 / math.Cos(math.Pi/6),
			angleMax: math.Pi / 6,
			want:     math.Pi / 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Angle(tt.v, tt.vMax, tt.angleMax); math.Abs(got-tt.want) > mathlib.EPS {
				t.Errorf("Angle() = %v, want %v", got, tt.want)
			}
		})
	}
}
