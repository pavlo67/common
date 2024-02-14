package plane

import "math"

// LeftAngleFromOx is angle in radians (-Pi < rotation <= Pi), rotation from Ox to Oy (counter clockwise) has a positive angle
type LeftAngleFromOx float64

// Position moves geometry shape to Point2 and (after!) rotates it with LeftAngleFromOx angle around itself
type Position struct {
	Point2
	LeftAngleFromOx
}

func (r LeftAngleFromOx) Canon() LeftAngleFromOx {
	if r > math.Pi {
		return r - 2*math.Pi
	} else if r <= -math.Pi {
		return r + 2*math.Pi
	}

	return r
}

func (r LeftAngleFromOx) Point2(radius float64) Point2 {
	return Point2{radius * math.Cos(float64(r)), radius * math.Sin(float64(r))}
}
