package geometry

import (
	"math"
)

type LineSegment [2]Point2

func (lineSegment *LineSegment) Vector() Point2 {
	if lineSegment == nil {
		return Point2{}
	}
	return Point2{lineSegment[1].X - lineSegment[0].X, lineSegment[1].Y - lineSegment[0].Y}
}

func (lineSegment LineSegment) Middle() Point2 {
	return Point2{(lineSegment[1].X + lineSegment[0].X) / 2, (lineSegment[1].Y + lineSegment[0].Y) / 2}
}

const eps = 1e-8

type Intersection struct {
	Angle  float64
	X0, X1 float64
}

func LinesIntersection(s0, s1 LineSegment) *Point2 {
	// https://stackoverflow.com/questions/7446126/opencv-2d-line-intersection-helper-function
	// https://stackoverflow.com/questions/563198/how-do-you-detect-where-two-line-segments-intersect

	r := Point2{s0[1].X - s0[0].X, s0[1].Y - s0[0].Y}
	s := Point2{s1[1].X - s1[0].X, s1[1].Y - s1[0].Y}

	cr := cross(r, s)
	if cr > -eps && cr < eps {
		// lines are about parallel (they may be collinear!)
		return nil
	}

	q := Point2{s1[0].X - s0[0].X, s1[0].Y - s0[0].Y}
	t := cross(q, s) / cr

	return &Point2{s0[0].X + t*r.X, s0[0].Y + t*r.Y}

}

func LineSegmentsIntersection(s0, s1 LineSegment) *Point2 {
	if s0[1].X < s0[0].X {
		s0 = LineSegment{s0[1], s0[0]}
	}
	if s1[1].X < s1[0].X {
		s1 = LineSegment{s1[1], s1[0]}
	}
	if s1[0].X < s0[0].X {
		s0, s1 = s1, s0
	}
	if s1[0].X-s0[1].X >= eps {
		return nil
	}

	r := Point2{s0[1].X - s0[0].X, s0[1].Y - s0[0].Y}
	s := Point2{s1[1].X - s1[0].X, s1[1].Y - s1[0].Y}

	cr := cross(r, s)
	if cr > -eps && cr < eps {
		// vertical segments
		if math.Abs(cross(r, Point2{0, 1})) < eps {
			if s0[1].X-s1[0].X >= eps {
				return nil
			}

			if s0[1].Y < s0[0].Y {
				s0 = LineSegment{s0[1], s0[0]}
			}
			if s1[1].Y < s1[0].Y {
				s1 = LineSegment{s1[1], s1[0]}
			}

			if s0[0].Y < s1[0].Y {
				if s0[1].Y >= s1[0].Y {
					return &s1[0]
				}
			} else if s0[0].Y > s1[0].Y {
				if s1[1].Y >= s0[0].Y {
					return &s0[0]
				}
			} else {
				return &s0[0]
			}
			return nil
		}

		// compare s1[0].Y and corresponding point on s0
		k := r.Y / r.X
		s01Y := s0[0].Y + k*(s1[0].X-s0[0].X)
		if math.Abs(s01Y/math.Sqrt(1+k*k)) >= eps {
			return nil
		}

		return &s1[0]
	}

	q := Point2{s1[0].X - s0[0].X, s1[0].Y - s0[0].Y}
	t := cross(q, s) / cr
	xIntersect := s0[0].X + t*r.X

	if xIntersect < s0[0].X || xIntersect > s0[1].X || xIntersect < s1[0].X || xIntersect > s1[1].X {
		return nil
	}

	return &Point2{xIntersect, s0[0].Y + t*r.Y}
}

func cross(p1, p2 Point2) float64 {
	return p1.X*p2.Y - p1.Y*p2.X
}

func DividedByLine(p0, p1 Point2, axis LineSegment) *Point2 {
	pIntersect := LinesIntersection(LineSegment{p0, p1}, axis)
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

func DistanceToLineSegment(p Point2, ls LineSegment) (distance, projectionPosition float64) {
	d0, d1, d := DistanceSquare(p, ls[0]), DistanceSquare(p, ls[1]), DistanceSquare(ls[0], ls[1])
	var reversed bool
	if d1 < d0 {
		d0, d1 = d1, d0
		reversed = true
	}
	if d0+d <= d1 {
		if reversed {
			return math.Sqrt(d0), math.Sqrt(d)
		} else {
			return math.Sqrt(d0), 0
		}
	}

	c0 := (d0 + d - d1) / (2 * math.Sqrt(d))

	if reversed {
		return math.Sqrt(d0 - c0*c0), math.Sqrt(d) - c0
	} else {
		return math.Sqrt(d0 - c0*c0), c0
	}
}
