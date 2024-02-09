package imagelib

import (
	"image"
	"math"

	"github.com/pavlo67/common/common/mathlib/plane"
)

// ...image.Point ----------------------------------------------------------

func Distance(el1, el2 image.Point) float64 {
	return math.Sqrt(float64((el1.X-el2.X)*(el1.X-el2.X) + (el1.Y-el2.Y)*(el1.Y-el2.Y)))
}

func ConvertImagePoints(points0 []image.Point, transpose bool, pMin image.Point, scale int) []image.Point {
	points := make([]image.Point, len(points0))

	for i, p := range points0 {
		if transpose {
			p = image.Point{p.Y, p.X}
		}
		points[i] = image.Point{pMin.X + scale*p.X, pMin.Y + scale*p.Y}
	}

	return points
}

func RectInternal(rect image.Rectangle, margin float64) image.Rectangle {
	marginPoint := image.Point{int(math.Ceil(margin)), int(math.Ceil(margin))}
	return image.Rectangle{Min: rect.Min.Add(marginPoint), Max: rect.Max.Sub(marginPoint)}
}

func RectangleAround(marginPix int, pts ...image.Point) image.Rectangle {
	if len(pts) < 1 {
		return image.Rectangle{}
	}

	if marginPix < 0 {
		marginPix = 0
	}

	minX, minY, maxX, maxY := pts[0].X, pts[0].Y, pts[0].X, pts[0].Y

	for _, p := range pts[1:] {
		if p.X < minX {
			minX = p.X
		} else if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		} else if p.Y > maxY {
			maxY = p.Y
		}
	}

	return image.Rectangle{
		Min: image.Point{minX - marginPix, minY - marginPix},
		Max: image.Point{maxX + 1 + marginPix, maxY + 1 + marginPix},
	}
}

// ...plane.Point2 ---------------------------------------------------------

func PointsFromPolyChain(points0 ...plane.Point2) []image.Point {
	points := make([]image.Point, len(points0))

	for i, p := range points0 {
		points[i] = image.Point{int(math.Round(p.X)), int(math.Round(p.Y))}
	}

	return points
}

func PointFramed(p plane.Point2, rect image.Rectangle) plane.Point2 {
	halfSideX, halfSideY := 0.5*float64(rect.Max.X-rect.Min.X), 0.5*float64(rect.Max.Y-rect.Min.Y)
	xImg, yImg := p.X-float64(rect.Min.X), p.Y-float64(rect.Min.Y)

	return plane.Point2{-halfSideX + xImg, halfSideY - yImg}
}

func Segment(p0, p1 image.Point) plane.Segment {
	return plane.Segment{{float64(p0.X), float64(p0.Y)}, {float64(p1.X), float64(p1.Y)}}

}

func PolyChain(points []image.Point) plane.PolyChain {
	polyChain := make(plane.PolyChain, len(points))
	for i, p := range points {
		polyChain[i].X, polyChain[i].Y = float64(p.X), float64(p.Y)
	}

	return polyChain
}

func Center(points ...image.Point) plane.Point2 {
	if len(points) < 1 {
		return plane.Point2{math.NaN(), math.NaN()}
	}
	var x, y float64
	for _, element := range points {
		x += float64(element.X)
		y += float64(element.Y)
	}

	n := float64(len(points))

	return plane.Point2{X: x / n, Y: y / n}
}

func CenterImage(points ...image.Point) image.Point {
	if len(points) < 1 {
		return image.Point{-1, -1}
	}
	var x, y int
	for _, element := range points {
		x += element.X
		y += element.Y
	}

	return image.Point{X: x / len(points), Y: y / len(points)}
}
