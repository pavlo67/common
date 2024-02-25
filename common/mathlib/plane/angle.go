package plane

import "math"

// XToYAngle is angle in radians (-Pi < rotation <= Pi), rotation from Ox to Oy (counter clockwise) has a positive angle
type XToYAngle float64

// Position moves geometry shape to Point2 and (after!) rotates it with XToYAngle angle around itself
type Position struct {
	Point2
	XToYAngle
}

func (xToYAngle XToYAngle) Canon() XToYAngle {
	if xToYAngle > math.Pi {
		return xToYAngle - 2*math.Pi
	} else if xToYAngle <= -math.Pi {
		return xToYAngle + 2*math.Pi
	}

	return xToYAngle
}

func (xToYAngle XToYAngle) Point2(radius float64) Point2 {
	return Point2{radius * math.Cos(float64(xToYAngle)), radius * math.Sin(float64(xToYAngle))}
}
