package geometry

import (
	"math"
)

// TODO!!!! be careful: all single point angles are calculated in range -pi < angle <= pi

func Angle(p Point2) float64 {
	if p.X == 0 {
		if p.Y > 0 {
			return math.Pi / 2
		} else if p.Y < 0 {
			return -math.Pi / 2
		} else {
			return math.NaN()
		}
	} else if p.X >= 0 {
		return math.Atan(p.Y / p.X)
	} else if p.Y >= 0 {
		return math.Atan(p.Y/p.X) + math.Pi
	} else {
		return math.Atan(p.Y/p.X) - math.Pi
	}
}

func Vector(p0, p1 Point2) Point2 {
	return Point2{X: p1.X - p0.X, Y: p1.Y - p0.Y}
}
