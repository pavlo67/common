package geometry

import (
	"reflect"
	"testing"
)

func TestAveragePolyChains(t *testing.T) {
	tests := []struct {
		name                   string
		polyChain0             PolyChain
		polyChain1             PolyChain
		minDistance            float64
		wantOk                 bool
		wantPolyChain0Averaged PolyChain
		wantPolyChain1Rest     []PolyChain
	}{
		{
			name:                   "",
			polyChain0:             PolyChain{{0, 0}, {0, 2}, {0, 4}},
			polyChain1:             PolyChain{{0, 4}, {0, 6}},
			minDistance:            0,
			wantOk:                 true,
			wantPolyChain0Averaged: PolyChain{{0, 0}, {0, 2}, {0, 4}},
			wantPolyChain1Rest:     []PolyChain{{{0, 4}, {0, 6}}},
		},
		{
			name:                   "",
			polyChain0:             PolyChain{{0, 0}, {0, 2}, {0, 4}},
			polyChain1:             PolyChain{{0, 3}, {0, 4}, {0, 6}},
			minDistance:            0,
			wantOk:                 true,
			wantPolyChain0Averaged: PolyChain{{0, 0}, {0, 2}, {0, 3}, {0, 4}},
			wantPolyChain1Rest:     []PolyChain{{{0, 4}, {0, 6}}},
		},
		{
			name:                   "",
			polyChain0:             PolyChain{{0, 0}, {0.05, 2}, {0, 4}},
			polyChain1:             PolyChain{{0.1, 3}, {0.1, 4}, {0.1, 6}},
			minDistance:            0.1,
			wantOk:                 true,
			wantPolyChain0Averaged: PolyChain{{0, 0}, {0.05, 2}, {0.07500000000000001, 3}, {0.05, 4}},
			wantPolyChain1Rest:     []PolyChain{{{0.05, 4}, {0.1, 6}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOk, gotAveraged, gotRest := AveragePolyChains(tt.polyChain0, tt.polyChain1, tt.minDistance)

			if gotOk != tt.wantOk {
				t.Errorf("Averagegeometry.PolyChains() gotOk = %t, wantOk %t", gotOk, tt.wantOk)
			}
			if !reflect.DeepEqual(gotAveraged, tt.wantPolyChain0Averaged) {
				t.Errorf("Averagegeometry.PolyChains() gotAveraged = %v, wantAveraged %v", gotAveraged, tt.wantPolyChain0Averaged)
			}
			if !reflect.DeepEqual(gotRest, tt.wantPolyChain1Rest) {
				t.Errorf("Averagegeometry.PolyChains() gotRest = %v, wantRest %v", gotRest, tt.wantPolyChain1Rest)
			}
		})
	}
}
