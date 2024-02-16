package plane

import "math"

// LeftAngle is angle in radians (-Pi < rotation <= Pi), rotation from Ox to Oy (counter clockwise) has a positive angle
type LeftAngle float64

// Position moves geometry shape to Point2 and (after!) rotates it with LeftAngle angle around itself
type Position struct {
	Point2
	LeftAngle
}

func (leftAngle LeftAngle) Canon() LeftAngle {
	if leftAngle > math.Pi {
		return leftAngle - 2*math.Pi
	} else if leftAngle <= -math.Pi {
		return leftAngle + 2*math.Pi
	}

	return leftAngle
}

func (leftAngle LeftAngle) Point2(radius float64) Point2 {
	return Point2{radius * math.Cos(float64(leftAngle)), radius * math.Sin(float64(leftAngle))}
}
