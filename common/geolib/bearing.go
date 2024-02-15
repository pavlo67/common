package geolib

import (
	"math"

	"github.com/pavlo67/common/common/mathlib/plane"
)

// Bearing if a "right-angle" in degrees
type Bearing Degrees

func (bearing Bearing) Canon() Bearing {
	for bearing >= 360 {
		bearing -= 360
	}
	for bearing < 0 {
		bearing += 360
	}

	return bearing
}

func BearingFromPoint(point plane.Point2) Bearing {
	bearingDegrees := 90 - (180 * point.LeftAngleFromOx() / math.Pi)

	for bearingDegrees >= 360 {
		bearingDegrees -= 360
	}
	for bearingDegrees < 0 {
		bearingDegrees += 360
	}

	return Bearing(bearingDegrees)
}

func BearingFromLeftAngle(rotation plane.LeftAngle) Bearing {
	bearingDegrees := -(180 * rotation / math.Pi)

	for bearingDegrees >= 360 {
		bearingDegrees -= 360
	}
	for bearingDegrees < 0 {
		bearingDegrees += 360
	}

	return Bearing(bearingDegrees)
}

// LeftAngle is measured from Oy (as well as Bearing itself)
func (bearing Bearing) LeftAngle() plane.LeftAngle {
	angle := plane.LeftAngle(-bearing * math.Pi / 180)
	if angle <= -math.Pi {
		return angle + 2*math.Pi
	} else if angle > math.Pi {
		return angle - 2*math.Pi
	}

	return angle
}

func (bearing Bearing) Point(radius float64) plane.Point2 {
	return (math.Pi * (0.5 - plane.LeftAngle(bearing)/180)).Point2(radius)
}
