package geolib

import (
	"math"

	"github.com/pavlo67/common/common/mathlib/plane"
)

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

func PointBearing(point plane.Point2) Bearing {
	bearingDegrees := 90 - (180 * point.LeftAngleFromOx() / math.Pi)

	for bearingDegrees >= 360 {
		bearingDegrees -= 360
	}
	for bearingDegrees < 0 {
		bearingDegrees += 360
	}

	return Bearing(bearingDegrees)
}

func PlaneBearing(rotation plane.LeftAngleFromOx) Bearing {
	bearingDegrees := -(180 * rotation / math.Pi)

	for bearingDegrees >= 360 {
		bearingDegrees -= 360
	}
	for bearingDegrees < 0 {
		bearingDegrees += 360
	}

	return Bearing(bearingDegrees)
}

// OYLeftAngle was previously named LeftAngleFromOx()
func (bearing Bearing) OYLeftAngle() float64 {
	angle := float64(-bearing * math.Pi / 180)
	if angle <= -math.Pi {
		return angle + 2*math.Pi
	} else if angle > math.Pi {
		return angle - 2*math.Pi
	}

	return angle
}

// OXLeftAngle was previously named OxyAngle()
func (bearing Bearing) OXLeftAngle() float64 {
	angle := math.Pi * (0.5 - float64(bearing)/180)
	if angle <= -2*math.Pi {
		return angle + 2*math.Pi
	} else if angle > 2*math.Pi {
		return angle - 2*math.Pi
	}

	return angle
}
