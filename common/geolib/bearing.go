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

func BearingFromGeometry(rotation plane.Rotation) Bearing {
	//angle := float64(rotation.Rotation + rotation.Rotation)

	bearingDegrees := 90 - (180 * rotation / math.Pi)

	for bearingDegrees >= 360 {
		bearingDegrees -= 360
	}
	for bearingDegrees < 0 {
		bearingDegrees += 360
	}

	return Bearing(bearingDegrees)
}

func (bearing Bearing) Rotation() plane.Rotation {
	angle := plane.Rotation(-bearing * math.Pi / 180)
	if angle <= -math.Pi {
		return angle + 2*math.Pi
	} else if angle > math.Pi {
		return angle - 2*math.Pi
	}

	return angle
}

func (bearing Bearing) OxyAngle() float64 {
	angle := math.Pi * (0.5 - float64(bearing)/180)
	if angle <= -2*math.Pi {
		return angle + 2*math.Pi
	} else if angle > 2*math.Pi {
		return angle - 2*math.Pi
	}

	return angle
}
