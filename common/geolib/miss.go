package geolib

import (
	"github.com/pavlo67/common/common/mathlib/plane"
)

func GeoError(geoPointReal, geoPoint Point) (missPoint plane.Point2, miss float64, bearing Bearing) {
	missDirection := geoPointReal.DirectionTo(geoPoint)

	return missDirection.Moving(), missDirection.Distance, missDirection.Bearing
}
