package convolution

//import (
//	"image"
//	"testing"
//
//	"github.com/pavlo67/maps/mapping/imager"
//	"github.com/pavlo67/maps/mapping/imager/layer"
//
//	"github.com/stretchr/testify/require"
//
//	"github.com/pavlo67/common/common/logger"
//	"github.com/pavlo67/common/common/logger/logger_zap"
//)
//
//func TestMaskAveraging(t *testing.T) {
//	l, err := logger_zap.New(logger.Config{SaveFiles: true})
//	require.NoError(t, err)
//	require.NotNil(t, l)
//
//	xWidth, yHeight := 16, 16
//	ls := imager.PreparationSettings{DPM: 5}
//
//	imgRGB := layer.Layer{
//		Pix:   make([]layer.Value, xWidth*yHeight),
//		Stride:   xWidth,
//		Rect:     image.Rectangle{Max: image.Contour{xWidth, yHeight}},
//		PreparationSettings: ls,
//	}
//
//	for y := 0; y < yHeight; y++ {
//		for x := 4; x < 8; x++ {
//			imgRGB.Pix[imgRGB.Offset(x, y)] = layer.ValueMax
//		}
//	}
//
//	for x := 0; x < xWidth; x++ {
//		for y := 3; y < 5; y++ {
//			imgRGB.Pix[imgRGB.Offset(x, y)] = layer.ValueMax
//		}
//	}
//
//	l.Image("imgRGB.png", &imgRGB)
//
//	averagingStep := 4
//
//	aggregateFeature, err := Erosion(averagingStep, 0)
//	require.NoError(t, err)
//	require.NotNil(t, aggregateFeature)
//
//	lyrConvolved, err := Layer(&imgRGB, ls.DPM, aggregateFeature, averagingStep, false)
//	require.NoError(t, err)
//	require.NotNil(t, lyrConvolved)
//
//	l.Image("lyr_convolved.png", lyrConvolved)
//
//}
