package imagelib

import (
	"image"
	"math"

	"github.com/pavlo67/common/common/mathlib/geometry"
)

type Bounded interface {
	Bounds() image.Rectangle
}

func RectangleAround(rect image.Rectangle, marginPix float64, pts ...geometry.Point2) image.Rectangle {
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

	return rect.Intersect(image.Rectangle{
		Min: geometry.Point2{minX - marginPix, minY - marginPix}.Point(),
		Max: geometry.Point2{maxX + marginPix, maxY + marginPix}.Point()})
}

func PolyChain(points []image.Point) geometry.PolyChain {
	polyChain := make(geometry.PolyChain, len(points))
	for i, p := range points {
		polyChain[i].X, polyChain[i].Y = float64(p.X), float64(p.Y)
	}

	return polyChain
}

func Distance(el1, el2 image.Point, dpm float64) float64 {
	return math.Sqrt(float64((el1.X-el2.X)*(el1.X-el2.X)+(el1.Y-el2.Y)*(el1.Y-el2.Y))) / dpm
}

func Direction(el1, el2 image.Point) float64 {
	dx := float64(el2.X) - float64(el1.X)
	dy := float64(el2.Y) - float64(el1.Y)

	if dx == 0 {
		if dy > 0 {
			return 0
		} else if dy < 0 {
			return 180
		} else {
			math.NaN()
		}
	}

	direction_ := 180 * math.Atan(dy/dx) / math.Pi
	if dx > 0 {
		return direction_
	} else if dy > 0 {
		return 180 + direction_
	}

	return -180 + direction_
}

func Center(points ...image.Point) geometry.Point2 {
	if len(points) < 1 {
		return geometry.Point2{math.NaN(), math.NaN()}
	}
	var x, y float64
	for _, element := range points {
		x += float64(element.X)
		y += float64(element.Y)
	}

	n := float64(len(points))

	return geometry.Point2{X: x / n, Y: y / n}
}

//func AverageAlongOx(points2 []numlib.Point2) ([]image.Point, image.Rect) {
//	if len(points2) < 1 {
//		return nil, image.Rect{}
//	}
//
//	sort.Slice(points2, func(i, j int) bool { return points2[i].Position < points2[j].Position })
//	pX := int(points2[0].Position)
//	if points2[0].Position < 0 {
//		pX--
//	}
//	xBase := -pX
//	yBase := 0
//	var yBaseI int
//	for _, p := range points2 {
//		if p.Y >= 0 {
//			yBaseI = int(p.Y)
//		} else {
//			yBaseI = -int(p.Y)
//		}
//
//		if yBaseI > yBase {
//			yBase = yBaseI
//		}
//	}
//
//	points := make([]image.Point, len(points2))
//
//	var yPlus, yMinus []float64
//	for _, p := range points2 {
//		pXNext := int(p.Position)
//		if p.Position < 0 {
//			pXNext--
//		}
//		if pXNext != pX {
//			points = append(points, AveragedPoint(pX+xBase, yBase+1, yPlus, yMinus))
//			yPlus, yMinus, pX = nil, nil, pXNext
//		}
//
//		if p.Y > 0 {
//			yPlus = append(yPlus, p.Y)
//		} else if p.Y < 0 {
//			yMinus = append(yMinus, p.Y)
//		} else {
//			yPlus = append(yPlus, p.Y)
//			yMinus = append(yMinus, p.Y)
//		}
//	}
//	if len(yPlus)+len(yMinus) > 0 {
//		points = append(points, AveragedPoint(pX+xBase, yBase, yPlus, yMinus))
//	}
//	toX := points[len(points)-1].Position + 1
//
//	for i := len(points) - 1; i >= 0; i-- {
//		points = append(points, image.Point{Position: points[i].Position, Y: yBase - points[i].Y})
//	}
//
//	return points, image.Rect{Max: image.Point{toX, yBase * 2}}
//}
//
//func AveragedPoint(x, yBase int, yPlus, yMinus []float64) image.Point {
//	var yPlusAvg, yMinusAvg float64
//
//	if len(yPlus) > 0 {
//		for _, y := range yPlus {
//			yPlusAvg += y
//		}
//		yPlusAvg /= float64(len(yPlus))
//	}
//
//	if len(yMinus) > 0 {
//		for _, y := range yMinus {
//			yMinusAvg += y
//		}
//		yMinusAvg /= float64(len(yMinus))
//	}
//
//	return image.Point{x, yBase + int(math.Round((yPlusAvg-yMinusAvg)/2))}
//}
