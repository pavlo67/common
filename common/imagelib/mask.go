package imagelib

import (
	"image"
	"image/color"

	"golang.org/x/image/draw"

	"github.com/pavlo67/common/common"

	"github.com/pavlo67/common/common/imagelib/coloring"
)

// GetMask ------------------------------------------------------------------------------

type GetMask interface {
	Color() *coloring.ColorNamed
	Mask(color.Color, common.Map) MasksOneColor
	Info(colorNamed coloring.ColorNamed) string
}

type MaskOneColor struct {
	Color  color.Color
	Points []image.Point
	Marker
}

type MasksOneColor []MaskOneColor

func (moc MasksOneColor) ShowOnRGBA(rgb image.RGBA) {
	for _, maskOneColor := range moc {
		for _, p := range maskOneColor.Points {
			rgb.Set(p.X, p.Y, maskOneColor.Color)
		}
	}
}

func (moc MasksOneColor) ShowOn(img image.Image) {
	drawImg, _ := img.(draw.Image)
	if drawImg != nil {
		for _, maskOneColor := range moc {
			for _, p := range maskOneColor.Points {
				drawImg.Set(p.X, p.Y, maskOneColor.Color)
			}
			if maskOneColor.Marker != nil {
				maskOneColor.Marker.Mark(drawImg, maskOneColor.Color)
			}
		}
	}
}

// PointsGetMask ------------------------------------------------------------------------

var _ GetMask = &PointsGetMask{}
var _ GetMask = PointsGetMask{}

type PointsGetMask struct {
	Points []image.Point
	*coloring.ColorNamed
	//PointSize    int
	//AddInfo      bool
	//DetailedInfo bool
	//FontFile     string
	//Title        string
}

func (pointsGetMask PointsGetMask) Color() *coloring.ColorNamed {
	if pointsGetMask.ColorNamed != nil {
		return pointsGetMask.ColorNamed
	}
	return nil
}

func (pointsGetMask PointsGetMask) Mask(clr color.Color, opts common.Map) MasksOneColor {
	return MasksOneColor{{Color: clr, Points: pointsGetMask.Points}}
}

func (pointsGetMask PointsGetMask) Info(colorNamed coloring.ColorNamed) string {
	return ""
	//if !pointsGetMask.AddInfo {
	//	return ""
	//}
	//
	//fiber := pointsGetMask.Fiber
	//
	//var title, details string
	//
	//if pointsGetMask.Title != "" {
	//	title = "\n" + pointsGetMask.Title
	//}
	//if pointsGetMask.DetailedInfo {
	//	details = fmt.Sprintf("%v", fiber.PolyChain)
	//}
	//
	//return title + fmt.Sprintf(
	//	"\nfiber#%d (%s, %s, %d, from = %v, to = %v, lengthPix = %f, widthAvgPix = %f, variationAvg = %f, points = %d (%s), color range = %v)",
	//	fiber.N, colorNamed.Name, fiber.Key, fiber.Direction.N, fiber.PolyChain[0], fiber.PolyChain[len(fiber.PolyChain)-1], fiber.Length,
	//	fiber.WidthAvg, fiber.VariationAvg, len(fiber.PolyChain), details, fiber.ColorRange)
	//
}
