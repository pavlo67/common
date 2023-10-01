package plane

import "math"

// Rotation is angle in radians (-Pi < rotation <= Pi), rotation from Ox to Oy (counter clockwise) has a positive angle
type Rotation float64

// Position moves geometry shape to Point2 and (after!) rotates it with Rotation angle around itself
type Position struct {
	Point2
	Rotation
}

func (r Rotation) Sub(r1 Rotation) Rotation {
	sub := r - r1
	if sub > math.Pi {
		return sub - 2*math.Pi
	} else if sub <= -math.Pi {
		return sub + 2*math.Pi
	}

	return sub
}
