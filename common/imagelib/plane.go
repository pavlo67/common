package imagelib

import (
	"image"

	"github.com/pavlo67/common/common/mathlib/plane"
)

func PointFramed(p plane.Point2, rect image.Rectangle) plane.Point2 {
	halfSideX, halfSideY := 0.5*float64(rect.Max.X-rect.Min.X), 0.5*float64(rect.Max.Y-rect.Min.Y)
	xImg, yImg := p.X-float64(rect.Min.X), p.Y-float64(rect.Min.Y)

	return plane.Point2{-halfSideX + xImg, halfSideY - yImg}
}

func Segment(p0, p1 image.Point) plane.Segment {
	return plane.Segment{{float64(p0.X), float64(p0.Y)}, {float64(p1.X), float64(p1.Y)}}

}
