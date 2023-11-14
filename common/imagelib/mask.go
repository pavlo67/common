package imagelib

import (
	"image"
	"image/color"

	"golang.org/x/image/draw"
)

// GetMask ------------------------------------------------------------------------------

type GetMask interface {
	Color() *ColorNamed
	Mask(color.Color) Mask
	Info(colorNamed ColorNamed) string
}

type MaskOneColor struct {
	Color  color.Color
	Points []image.Point
	Marker
}

type Mask []MaskOneColor

func (mask Mask) ShowOnRGBA(rgb image.RGBA) {
	for _, maskOneColor := range mask {
		for _, p := range maskOneColor.Points {
			rgb.Set(p.X, p.Y, maskOneColor.Color)
		}
	}
}

func (mask Mask) ShowOn(img image.Image) {
	drawImg, _ := img.(draw.Image)
	if drawImg != nil {
		for _, maskOneColor := range mask {
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
	*ColorNamed
	//PointSize    int
	//AddInfo      bool
	//DetailedInfo bool
	//FontFile     string
	//Title        string
}

func (pointsGetMask PointsGetMask) Color() *ColorNamed {
	if pointsGetMask.ColorNamed != nil {
		return pointsGetMask.ColorNamed
	}
	return nil
}

func (pointsGetMask PointsGetMask) Mask(clr color.Color) Mask {
	return Mask{{Color: clr, Points: pointsGetMask.Points}}
}

func (pointsGetMask PointsGetMask) Info(colorNamed ColorNamed) string {
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
