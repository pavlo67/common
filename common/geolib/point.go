package geolib

import (
	"math"

	geo "github.com/kellydunn/golang-geo"
	"github.com/pavlo67/common/common/mathlib/plane"
	// geo "github.com/billups/golang-geo"
)

type Point struct {
	Lat, Lon Degrees
}

type Direction struct {
	Bearing
	Distance float64
}

func (dir Direction) Moving() plane.Point2 {
	localAngle := dir.Bearing.OXLeftAngle()

	return plane.Point2{dir.Distance * math.Cos(localAngle), dir.Distance * math.Sin(localAngle)}
}

// https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames

func (p Point) Tile(zoom int) Tile {
	// TODO!!! check zoom

	tile := Tile{Zoom: zoom}

	n := math.Exp2(float64(zoom))

	tile.X = int(math.Floor(float64((p.Lon+180.0)/360.0) * n))
	if float64(tile.X) >= n {
		tile.X = int(n - 1)
	}

	latAngle := float64(p.Lat) * math.Pi / 180
	tile.Y = int(math.Floor(n * (1 - math.Log(math.Tan(latAngle)+1/math.Cos(latAngle))/math.Pi) / 2))

	// log.Print(p, zoom, tile)

	return tile
}

func (p Point) Geo() geo.Point {
	return *geo.NewPoint(float64(p.Lat), float64(p.Lon))
}

func (p Point) MovedAt(point2 plane.Point2) Point {
	if point2.X == 0 && point2.Y == 0 {
		return p
	}

	dxKm, dyKm := point2.X*0.001, point2.Y*0.001

	bearing := DirectionBearingFromRotation(plane.Point2{dxKm, dyKm}.Rotation())

	geoPoint := p.Geo()

	// geoPoint.PointAtDistanceAndBearing() requires distance in kilometers
	geoPointNew := geoPoint.PointAtDistanceAndBearing(math.Sqrt(dxKm*dxKm+dyKm*dyKm), float64(bearing))

	return Point{Degrees(geoPointNew.Lat()), Degrees(geoPointNew.Lng())}
}

func (p Point) PointAtDirection(dir Direction) Point {
	geoPoint := p.Geo()

	// geoPoint.PointAtDistanceAndBearing() requires distance in kilometers
	geoPointNew := geoPoint.PointAtDistanceAndBearing(dir.Distance*0.001, float64(dir.Bearing))

	return Point{Degrees(geoPointNew.Lat()), Degrees(geoPointNew.Lng())}
}

func (p Point) BearingTo(p1 Point) Bearing {
	geoPoint, geoPoint1 := p.Geo(), p1.Geo()
	return Bearing(geoPoint.BearingTo(&geoPoint1))
}

func (p Point) DistanceTo(p1 Point) float64 {
	geoPoint, geoPoint1 := p.Geo(), p1.Geo()

	// geoPoint.GreatCircleDistance(&geoPoint1) returns distance in kilometers
	return 1000 * geoPoint.GreatCircleDistance(&geoPoint1)
}

func (p Point) DirectionTo(p1 Point) Direction {
	geoPoint, geoPoint1 := p.Geo(), p1.Geo()
	return Direction{
		Bearing(geoPoint.BearingTo(&geoPoint1)),
		1000 * geoPoint.GreatCircleDistance(&geoPoint1),
	}
}
