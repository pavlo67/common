package plane

import (
	"math"

	"github.com/pavlo67/common/common/mathlib"
)

type Segment [2]Point2

func (s *Segment) Vector() Point2 {
	if s == nil {
		return Point2{}
	}
	return Point2{s[1].X - s[0].X, s[1].Y - s[0].Y}
}

func (s Segment) Middle() Point2 {
	return Point2{(s[1].X + s[0].X) / 2, (s[1].Y + s[0].Y) / 2}
}

func SegmentsIntersection(s, s1 Segment) (pCross *Point2, atEnd bool) {
	if s[1].X < s[0].X {
		s = Segment{s[1], s[0]}
	}
	if s1[1].X < s1[0].X {
		s1 = Segment{s1[1], s1[0]}
	}
	if s1[0].X < s[0].X {
		s, s1 = s1, s
	}
	if s1[0].X-s[1].X >= mathlib.Eps {
		return nil, false
	}

	r := Point2{s[1].X - s[0].X, s[1].Y - s[0].Y}
	l := Point2{s1[1].X - s1[0].X, s1[1].Y - s1[0].Y}

	cr := Cross(r, l)
	if cr > -mathlib.Eps && cr < mathlib.Eps {
		// vertical segments
		if math.Abs(Cross(r, Point2{0, 1})) < mathlib.Eps {
			if s[1].X-s1[0].X >= mathlib.Eps {
				return nil, false
			}

			if s[1].Y < s[0].Y {
				s = Segment{s[1], s[0]}
			}
			if s1[1].Y < s1[0].Y {
				s1 = Segment{s1[1], s1[0]}
			}

			if s[0].Y < s1[0].Y {
				if s[1].Y >= s1[0].Y {
					return &s1[0], true
				}
			} else if s[0].Y > s1[0].Y {
				if s1[1].Y >= s[0].Y {
					return &s[0], true
				}
			} else {
				return &s[0], true
			}
			return nil, false
		}

		// compare s1[0].Y and corresponding point on s
		k := r.Y / r.X
		s01Y := s[0].Y + k*(s1[0].X-s[0].X)
		if math.Abs(s01Y/math.Sqrt(1+k*k)) >= mathlib.Eps {
			return nil, false
		}

		return &s1[0], true
	} else if math.Abs(s[0].X-s[1].X) < mathlib.Eps {
		s, s1 = s1, s
		r = Point2{s[1].X - s[0].X, s[1].Y - s[0].Y}
		l = Point2{s1[1].X - s1[0].X, s1[1].Y - s1[0].Y}
		cr = Cross(r, l)
	}

	q := Point2{s1[0].X - s[0].X, s1[0].Y - s[0].Y}
	t := Cross(q, l) / cr
	xIntersect := s[0].X + t*r.X

	if xIntersect < s[0].X || xIntersect > s[1].X || xIntersect < s1[0].X || xIntersect > s1[1].X {
		return nil, false
	} else if xIntersect == s[0].X || xIntersect == s[1].X || xIntersect == s1[0].X || xIntersect == s1[1].X {
		//} else if xIntersect-s[0].X <= Eps || s[1].X-xIntersect <= Eps || xIntersect-s1[0].X <= Eps || s1[1].X-xIntersect <= Eps {
		return &Point2{xIntersect, s[0].Y + t*r.Y}, true
	}

	return &Point2{xIntersect, s[0].Y + t*r.Y}, false
}

func (s Segment) DistanceTo(s1 Segment) float64 {
	if pCross, _ := SegmentsIntersection(s, s1); pCross != nil {
		return 0
	}

	d00, _ := s[0].DistanceToSegment(s1)
	d01, _ := s[1].DistanceToSegment(s1)
	d10, _ := s1[0].DistanceToSegment(s)
	d11, _ := s1[1].DistanceToSegment(s)

	return min(d00, d01, d10, d11)
}

func (s Segment) AngleAbs(s1 Segment) float64 {
	angle := math.Abs(s.Vector().AnglesDelta(s1.Vector()))
	for angle > math.Pi {
		angle -= math.Pi
	}
	if angle <= 0.5*math.Pi {
		return angle
	}

	return math.Pi - angle
}
