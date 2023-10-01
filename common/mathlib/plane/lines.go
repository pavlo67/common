package plane

import (
	"math"

	"github.com/pavlo67/common/common/mathlib"
)

type Intersection struct {
	Angle  float64
	X0, X1 float64
}

func Cross(p1, p2 Point2) float64 {
	return p1.X*p2.Y - p1.Y*p2.X
}

func (s Segment) TurnAroundAxis(p2 Point2) Point2 {
	axisDX, axisDY := s[1].X-s[0].X, s[1].Y-s[0].Y
	var axisDerivative float64

	if axisDX == 0 {
		if axisDY == 0 {
			// TODO!!! be careful, it's a very non-standard case
			return Point2{s[0].X*2 - p2.X, s[0].Y*2 - p2.Y}
		}
		return Point2{s[0].X*2 - p2.X, p2.Y}
	} else if axisDerivative = axisDY / axisDX; math.IsInf(axisDerivative, 0) {
		return Point2{s[0].X*2 - p2.X, p2.Y}
	} else if axisDY == 0 || math.IsInf(1/axisDerivative, 0) {
		return Point2{p2.X, s[0].Y*2 - p2.Y}
	}

	pIntersection := s.LinesIntersection(Segment{p2, Point2{p2.X + 1, p2.Y - 1/axisDerivative}})
	return Point2{pIntersection.X*2 - p2.X, pIntersection.Y*2 - p2.Y}

}

func (s Segment) TurnAroundAxisMultiple(p2s ...Point2) []Point2 {
	axisDX, axisDY := s[1].X-s[0].X, s[1].Y-s[0].Y
	var axisDerivative float64

	p2sTurned := make([]Point2, len(p2s))

	if axisDX == 0 {
		if axisDY == 0 {
			// TODO!!! be careful, it's a very non-standard case
			for i, p2 := range p2s {
				p2sTurned[i] = Point2{s[0].X*2 - p2.X, s[0].Y*2 - p2.Y}
			}
		} else {
			for i, p2 := range p2s {
				p2sTurned[i] = Point2{s[0].X*2 - p2.X, p2.Y}
			}
		}
	} else if axisDerivative = axisDY / axisDX; math.IsInf(axisDerivative, 0) {
		for i, p2 := range p2s {
			p2sTurned[i] = Point2{s[0].X*2 - p2.X, p2.Y}
		}
	} else if axisDY == 0 || math.IsInf(1/axisDerivative, 0) {
		for i, p2 := range p2s {
			p2sTurned[i] = Point2{p2.X, s[0].Y*2 - p2.Y}
		}
	} else {
		for i, p2 := range p2s {
			pIntersection := s.LinesIntersection(Segment{p2, Point2{p2.X + 1, p2.Y - 1/axisDerivative}})
			p2sTurned[i] = Point2{pIntersection.X*2 - p2.X, pIntersection.Y*2 - p2.Y}
		}
	}

	return p2sTurned
}

func (s Segment) LinesIntersection(s1 Segment) *Point2 {
	// https://stackoverflow.com/questions/7446126/opencv-2d-line-intersection-helper-function
	// https://stackoverflow.com/questions/563198/how-do-you-detect-where-two-line-segments-intersect

	r := Point2{s[1].X - s[0].X, s[1].Y - s[0].Y}
	l := Point2{s1[1].X - s1[0].X, s1[1].Y - s1[0].Y}

	cr := Cross(r, l)
	if cr > -mathlib.Eps && cr < mathlib.Eps {
		// lines are about parallel (they may be collinear!)
		return nil
	}

	q := Point2{s1[0].X - s[0].X, s1[0].Y - s[0].Y}
	t := Cross(q, l) / cr

	return &Point2{s[0].X + t*r.X, s[0].Y + t*r.Y}

}

func (s Segment) DivideByLine(p0, p1 Point2) *Point2 {
	pIntersect := s.LinesIntersection(Segment{p0, p1})
	if pIntersect == nil {
		return nil
	}

	if p1.X > p0.X {
		if pIntersect.X > p0.X && pIntersect.X < p1.X {
			return pIntersect
		}
	} else if p1.X < p0.X {
		if pIntersect.X > p1.X && pIntersect.X < p0.X {
			return pIntersect
		}
	} else if p1.Y > p0.Y {
		if pIntersect.Y > p0.Y && pIntersect.Y < p1.Y {
			return pIntersect
		}
	} else {
		if pIntersect.Y > p1.Y && pIntersect.Y < p0.Y {
			return pIntersect
		}
	}

	// log.Print(p0, p1, pIntersect)

	return nil
}
