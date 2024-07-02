package imagelib

import (
	"image"
	"image/color"
	"math"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/mathlib/plane"
	"golang.org/x/image/colornames"

	"github.com/pavlo67/common/common/errors"
)

func GrayFromPoints(points []image.Point, rect *image.Rectangle) image.Gray {
	if len(points) < 1 {
		return image.Gray{}
	}

	if rect == nil {
		rect = &image.Rectangle{Min: points[0], Max: image.Point{points[0].X + 1, points[0].Y + 1}}
		for _, p := range points[1:] {
			if p.X >= rect.Max.X {
				rect.Max.X = p.X + 1
			} else if p.X < rect.Min.X {
				rect.Min.X = p.X
			}
			if p.Y >= rect.Max.Y {
				rect.Max.Y = p.Y + 1
			} else if p.Y < rect.Min.Y {
				rect.Min.Y = p.Y
			}
		}
	}

	gray := image.Gray{
		Pix:    make([]uint8, (rect.Max.X-rect.Min.X)*(rect.Max.Y-rect.Min.Y)),
		Stride: rect.Max.X - rect.Min.X,
		Rect:   *rect,
	}
	for _, p := range points {
		gray.Set(p.X, p.Y, colornames.White)
	}

	return gray
}

func GrayFromPoints2(points2 []plane.Point2, rect *image.Rectangle) image.Gray {
	if len(points2) < 1 {
		return image.Gray{}
	}

	var dX, dY int
	if rect == nil {
		p0 := points2[0].ImagePoint()
		rect = &image.Rectangle{Min: p0, Max: image.Point{p0.X + 1, p0.Y + 1}}
		for _, p2 := range points2[1:] {
			p := p2.ImagePoint()
			if p.X >= rect.Max.X {
				rect.Max.X = p.X + 1
			} else if p.X < rect.Min.X {
				rect.Min.X = p.X
			}
			if p.Y >= rect.Max.Y {
				rect.Max.Y = p.Y + 1
			} else if p.Y < rect.Min.Y {
				rect.Min.Y = p.Y
			}
		}
		if rect.Min.X < 0 {
			dX = -rect.Min.X
			rect.Min.X, rect.Max.X = 0, rect.Max.X+dX
		}
		if rect.Min.Y < 0 {
			dY = -rect.Min.Y
			rect.Min.Y, rect.Max.Y = 0, rect.Max.Y+dY
		}
	}

	gray := image.Gray{
		Pix:    make([]uint8, (rect.Max.X-rect.Min.X)*(rect.Max.Y-rect.Min.Y)),
		Stride: rect.Max.X - rect.Min.X,
		Rect:   *rect,
	}
	for _, p := range points2 {
		gray.Set(int(math.Round(p.X))+dX, int(math.Round(p.Y))+dY, colornames.White)
	}

	return gray
}

type PointsImage struct {
	Points []image.Point
	image.Rectangle
}

var _ logger.GetImage = &PointsImage{}

func (imageOp *PointsImage) Bounds() image.Rectangle {
	if imageOp == nil {
		return image.Rectangle{}
	}

	return imageOp.Rectangle
}

func (imageOp *PointsImage) Image(opts common.Map) (image.Image, string, error) {
	if imageOp == nil {
		return nil, "", errors.New("*PointsImage = nil")
	}

	return PointsToGrayscale(imageOp.Points, imageOp.Rectangle), "", nil
}

// PointsToGrayscale returns Gray_ structure instead of image.Gray because the structure implements Imager interface required for show.Demo()
func PointsToGrayscale(points []image.Point, rect image.Rectangle) image.Image {
	xWidth := rect.Max.X - rect.Min.X
	yHeight := rect.Max.Y - rect.Min.Y

	gray := image.Gray{
		Pix:    make([]uint8, xWidth*yHeight),
		Stride: xWidth,
		Rect:   rect,
	}

	for _, p := range points {
		gray.Set(p.X, p.Y, color.White)
	}

	return &gray
}

func PointsToGrayscalePng(points []image.Point, rect image.Rectangle, path string) error {
	return SavePNG(PointsToGrayscale(points, rect), path)
}
