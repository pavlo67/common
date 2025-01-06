package convolution

//import (
//	"image"
//	"image/color"
//	"testing"
//
//	"github.com/stretchr/testify/require"
//
//	"github.com/pavlo67/common/common/imagelib"
//	"github.com/pavlo67/maps/mapping/imager/layer"
//)
//
//func TestMaskDirection(t *testing.T) {
//	tests := []struct {
//		name              string
//		points            []RGBAPoint
//		rect              image.Rectangle
//		side            []int
//		differenceMax     pix.ValueDelta
//		gapMaxPix         int
//		garbageRatioMax   float64
//		extensionRatioMin float64
//		x                 int
//		y                 int
//		want              layer.Value
//	}{
//		{
//			name: "1/full",
//			rect: image.Rect(0, 0, 10, 10),
//			points: []RGBAPoint{
//				{1, 2, color.RGBA{127, 127, 127, 255}},
//				{2, 1, color.RGBA{127, 127, 127, 255}},
//				{2, 3, color.RGBA{127, 127, 127, 255}},
//				{3, 2, color.RGBA{127, 127, 127, 255}},
//				//
//				{2, 2, color.RGBA{127, 127, 127, 255}},
//			},
//			side:            []int{1},
//			differenceMax:     20,
//			gapMaxPix:         1,
//			garbageRatioMax:   0,
//			extensionRatioMin: 0,
//			x:                 2,
//			y:                 2,
//			want:              0,
//		},
//		{
//			name: "1/2",
//			rect: image.Rect(0, 0, 10, 10),
//			points: []RGBAPoint{
//				{1, 2, color.RGBA{127, 127, 127, 255}},
//				{3, 2, color.RGBA{127, 127, 127, 255}},
//				//
//				{2, 2, color.RGBA{127, 127, 127, 255}},
//			},
//			side:            []int{1},
//			differenceMax:     20,
//			gapMaxPix:         0,
//			garbageRatioMax:   0,
//			extensionRatioMin: 0,
//			x:                 2,
//			y:                 2,
//			want:              1,
//		},
//		{
//			name: "2/full",
//			rect: image.Rect(0, 0, 10, 10),
//			points: []RGBAPoint{
//				{1, 0, color.RGBA{127, 127, 127, 255}},
//				{2, 0, color.RGBA{127, 127, 127, 255}},
//				{3, 0, color.RGBA{127, 127, 127, 255}},
//				{4, 1, color.RGBA{127, 127, 127, 255}},
//				{4, 2, color.RGBA{127, 127, 127, 255}},
//				{4, 3, color.RGBA{127, 127, 127, 255}},
//				{3, 4, color.RGBA{127, 127, 127, 255}},
//				{2, 4, color.RGBA{127, 127, 127, 255}},
//				{1, 4, color.RGBA{127, 127, 127, 255}},
//				{0, 3, color.RGBA{127, 127, 127, 255}},
//				{0, 2, color.RGBA{127, 127, 127, 255}},
//				{0, 1, color.RGBA{127, 127, 127, 255}},
//				//
//				{2, 2, color.RGBA{127, 127, 127, 255}},
//			},
//			side:            []int{2},
//			differenceMax:     20,
//			gapMaxPix:         0,
//			garbageRatioMax:   0,
//			extensionRatioMin: 0,
//			x:                 2,
//			y:                 2,
//			want:              0,
//		},
//		{
//			name: "2/2",
//			rect: image.Rect(0, 0, 10, 10),
//			points: []RGBAPoint{
//				{1, 0, color.RGBA{127, 127, 127, 255}},
//				{2, 0, color.RGBA{127, 127, 127, 255}},
//				{3, 0, color.RGBA{127, 127, 127, 255}},
//				{3, 4, color.RGBA{127, 127, 127, 255}},
//				{2, 4, color.RGBA{127, 127, 127, 255}},
//				{1, 4, color.RGBA{127, 127, 127, 255}},
//				//
//				{2, 2, color.RGBA{127, 127, 127, 255}},
//			},
//			side:            []int{2},
//			differenceMax:     20,
//			gapMaxPix:         0,
//			garbageRatioMax:   0,
//			extensionRatioMin: 0,
//			x:                 2,
//			y:                 2,
//			want:              2,
//		},
//		{
//			name: "3/1",
//			rect: image.Rect(0, 0, 10, 10),
//			points: []RGBAPoint{
//				{1, 0, color.RGBA{127, 127, 127, 255}},
//				{2, 0, color.RGBA{127, 127, 127, 255}},
//				{3, 0, color.RGBA{127, 127, 127, 255}},
//				{3, 4, color.RGBA{127, 127, 127, 255}},
//				{2, 4, color.RGBA{127, 127, 127, 255}},
//				{1, 4, color.RGBA{127, 127, 127, 255}},
//				//
//				{2, 2, color.RGBA{127, 127, 127, 255}},
//			},
//			side:            []int{3},
//			differenceMax:     20,
//			gapMaxPix:         0,
//			garbageRatioMax:   0,
//			extensionRatioMin: 0,
//			x:                 2,
//			y:                 2,
//			want:              0,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			dirMask, err := Direction(tt.side, tt.differenceMax, tt.gapMaxPix, tt.garbageRatioMax, tt.extensionRatioMin)
//			require.NoError(t, err)
//			require.NotNil(t, dirMask)
//
//			MaskTestScenario(t, dirMask, tt.rect, tt.points, tt.x, tt.y, tt.want)
//		})
//	}
//}
