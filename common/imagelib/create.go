package imagelib

import (
	"image"
	"math"

	"golang.org/x/image/colornames"

	"github.com/pavlo67/common/common/mathlib/geometry"
)

// ----------------------------------------------------------------

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

func GrayFromPoints2(points2 []geometry.Point2, rect *image.Rectangle) image.Gray {
	if len(points2) < 1 {
		return image.Gray{}
	}

	var dX, dY int
	if rect == nil {
		p0 := points2[0].Point()
		rect = &image.Rectangle{Min: p0, Max: image.Point{p0.X + 1, p0.Y + 1}}
		for _, p2 := range points2[1:] {
			p := p2.Point()
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