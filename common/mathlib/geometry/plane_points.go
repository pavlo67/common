package geometry

import (
	"image"
	"math"
)

type Point2 struct {
	X, Y float64
}

func (p2 Point2) Point() image.Point {
	return image.Point{int(math.Round(p2.X)), int(math.Round(p2.Y))}
}

func Vector(p0, p1 Point2) Point2 {
	return Point2{X: p1.X - p0.X, Y: p1.Y - p0.Y}
}

// TODO!!!! be careful: all single point angles are calculated in range -pi < angle <= pi

func Angle(p Point2) float64 {
	if p.X == 0 {
		if p.Y > 0 {
			return math.Pi / 2
		} else {
			return -math.Pi / 2
		}
	} else if p.X >= 0 {
		return math.Atan(p.Y / p.X)
	} else if p.Y >= 0 {
		return math.Atan(p.Y/p.X) + math.Pi
	} else {
		return math.Atan(p.Y/p.X) - math.Pi
	}
}

// !!! is equal to Angle() but more complicate
//func Angle1(p Point2) float64 {
//	x, y := p.Position, p.Y
//	if x == 0 {
//		if y == 0 {
//			return math.NaN()
//		} else if y > 0 {
//			return math.Pi / 2
//		}
//		return -math.Pi / 2
//	}
//
//	yx := y / x
//	if math.IsInf(yx, 1) {
//		return math.Pi / 2
//	} else if math.IsInf(yx, -1) {
//		return -math.Pi / 2
//	}
//
//	return math.Atan2(y, x)
//}

func Angle2(v0, v1 Point2) float64 {
	angle := Angle(v1) - Angle(v0)
	if angle > math.Pi {
		return angle - 2*math.Pi
	} else if angle <= -math.Pi {
		return angle + 2*math.Pi
	}
	return angle
}

func Distance(p0, p1 Point2) float64 {
	return math.Sqrt((p0.X-p1.X)*(p0.X-p1.X) + (p0.Y-p1.Y)*(p0.Y-p1.Y))
}

func DistanceSquare(p0, p1 Point2) float64 {
	return (p0.X-p1.X)*(p0.X-p1.X) + (p0.Y-p1.Y)*(p0.Y-p1.Y)
}

func Center(points ...Point2) Point2 {
	if len(points) < 1 {
		return Point2{math.NaN(), math.NaN()}
	}
	var x, y float64
	for _, element := range points {
		x += element.X
		y += element.Y
	}

	n := float64(len(points))

	return Point2{X: x / n, Y: y / n}
}

func TurnAroundAxis(axis LineSegment, p2 Point2) Point2 {
	axisDX, axisDY := axis[1].X-axis[0].X, axis[1].Y-axis[0].Y
	var axisDerivative float64

	if axisDX == 0 {
		if axisDY == 0 {
			// TODO!!! be careful, it's a very non-standard case
			return Point2{axis[0].X*2 - p2.X, axis[0].Y*2 - p2.Y}
		}
		return Point2{axis[0].X*2 - p2.X, p2.Y}
	} else if axisDerivative = axisDY / axisDX; math.IsInf(axisDerivative, 0) {
		return Point2{axis[0].X*2 - p2.X, p2.Y}
	} else if axisDY == 0 || math.IsInf(1/axisDerivative, 0) {
		return Point2{p2.X, axis[0].Y*2 - p2.Y}
	}

	pIntersection := LinesIntersection(LineSegment{p2, Point2{p2.X + 1, p2.Y - 1/axisDerivative}}, axis)
	return Point2{pIntersection.X*2 - p2.X, pIntersection.Y*2 - p2.Y}

}

func TurnAroundAxisMultiple(axis LineSegment, p2s ...Point2) []Point2 {
	axisDX, axisDY := axis[1].X-axis[0].X, axis[1].Y-axis[0].Y
	var axisDerivative float64

	p2sTurned := make([]Point2, len(p2s))

	if axisDX == 0 {
		if axisDY == 0 {
			// TODO!!! be careful, it's a very non-standard case
			for i, p2 := range p2s {
				p2sTurned[i] = Point2{axis[0].X*2 - p2.X, axis[0].Y*2 - p2.Y}
			}
		} else {
			for i, p2 := range p2s {
				p2sTurned[i] = Point2{axis[0].X*2 - p2.X, p2.Y}
			}
		}
	} else if axisDerivative = axisDY / axisDX; math.IsInf(axisDerivative, 0) {
		for i, p2 := range p2s {
			p2sTurned[i] = Point2{axis[0].X*2 - p2.X, p2.Y}
		}
	} else if axisDY == 0 || math.IsInf(1/axisDerivative, 0) {
		for i, p2 := range p2s {
			p2sTurned[i] = Point2{p2.X, axis[0].Y*2 - p2.Y}
		}
	} else {
		for i, p2 := range p2s {
			pIntersection := LinesIntersection(LineSegment{p2, Point2{p2.X + 1, p2.Y - 1/axisDerivative}}, axis)
			p2sTurned[i] = Point2{pIntersection.X*2 - p2.X, pIntersection.Y*2 - p2.Y}
		}
	}

	return p2sTurned
}

func RotateByAngle(p Point2, addAngle float64) Point2 {
	angle := Angle(p)
	r := math.Sqrt(p.X*p.X + p.Y*p.Y)

	return Point2{r * math.Cos(angle+addAngle), r * math.Sin(angle+addAngle)}
}

func RotateWithRatio(p Point2, ratio float64) Point2 {
	angle := Angle(p)
	r := math.Sqrt(p.X*p.X + p.Y*p.Y)

	return Point2{r * math.Cos(angle*ratio), r * math.Sin(angle*ratio)}
}

// TODO!!!! be careful: axis angles are calculated in range -pi/2 < angle <= pi/2
// TODO!!!! be careful:	if math.Abs(angle - math.Pi / 2) <= eps then deltaFromCenter == DX, else == DY

//func CenterAxis(axis LineSegment) (angle, deltaFromCenter float64) {
//	axisDX, axisDY := axis.End.Position-axis.Begin.Position, axis.End.Y-axis.Begin.Y
//
//	if axisDX == 0 {
//		if axisDY == 0 {
//			// TODO!!! be careful: it's a convention only
//			return 0, axis.Begin.Y
//		}
//		return math.Pi / 2, axis.Begin.Position
//	}
//
//	axisDerivative := axisDY / axisDX
//	if math.IsInf(axisDerivative, 0) {
//		return math.Pi / 2, axis.Begin.Position
//	}
//
//	angle = math.Atan(axisDerivative)
//	if math.Abs(angle-math.Pi/2) <= eps {
//		return angle, axis.Begin.Position - axis.Begin.Y/axisDerivative
//	}
//
//	return angle, axis.Begin.Y - axis.Begin.Position*axisDerivative
//}
//
//func SetAlongOx(points []Point2, axis LineSegment) ([]Point2, image.Rect) {
//	angle, deltaFromCenter := CenterAxis(axis)
//	pointsMoved := make([]Point2, len(points))
//
//	if math.Abs(angle-math.Pi/2) <= eps {
//		for i, p := range points {
//			pointsMoved[i] = RotateByAngle(Point2{p.Position - deltaFromCenter, p.Y}, -angle)
//		}
//	} else {
//		for i, p := range points {
//			pointsMoved[i] = RotateByAngle(Point2{p.Position, p.Y - deltaFromCenter}, -angle)
//		}
//	}
//
//	return pointsMoved
//}

func TriangleArea(p0, p1, p2 Point2) float64 {
	a, b, c := Distance(p0, p1), Distance(p1, p2), Distance(p2, p0)

	//log.Printf("p0 (%v), p1 (%v), p2 (%v) --> a (%f), b (%f), c (%f)", p0, p1, p2, a, b, c)

	p := (a + b + c) / 2

	return math.Sqrt(p * (p - a) * (p - b) * (p - c))
}

func DistanceToLine(p Point2, line LineSegment) float64 {
	p2 := line.Vector()
	pIntersect := LinesIntersection(LineSegment{p, Point2{p.X - p2.Y, p.Y + p2.X}}, line)
	if pIntersect == nil {
		// TODO!!! be careful: it's impossible
		return math.NaN()
	}

	return math.Sqrt((pIntersect.X-p.X)*(pIntersect.X-p.X) + (pIntersect.Y-p.Y)*(pIntersect.Y-p.Y))
}
