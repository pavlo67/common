package geometry

import (
	"math"
)

type PolyChain []Point2

func (polyChain PolyChain) Length() float64 {
	length := 0.
	for i := 1; i < len(polyChain); i++ {
		length += Distance(polyChain[i-1], polyChain[i])
	}
	return length
}

func (polyChain PolyChain) Reversed() PolyChain {
	polyChainReversed := make([]Point2, len(polyChain))
	for i, p2 := range polyChain {
		polyChainReversed[len(polyChain)-i-1] = p2
	}
	return polyChainReversed
}

func StraightenPolyChain(polyChain PolyChain, minDeviation float64) PolyChain {
	if len(polyChain) <= 2 {
		return polyChain
	}

	straightenedPolyChain := PolyChain{polyChain[0]}
	prevMedium := LineSegment{polyChain[0], polyChain[1]}.Middle()
	var prevMediumAdded bool

	for i := 1; i < len(polyChain)-1; i++ {
		medium := LineSegment{polyChain[i], polyChain[i+1]}.Middle()
		if distance := DistanceToLine(polyChain[i], LineSegment{prevMedium, medium}); distance > minDeviation {
			straightenedPolyChain = append(straightenedPolyChain, polyChain[i])
			prevMediumAdded = false
		} else if prevMediumAdded {
			straightenedPolyChain = append(straightenedPolyChain, medium)
			prevMediumAdded = true
		} else {
			straightenedPolyChain = append(straightenedPolyChain, prevMedium, medium)
			prevMediumAdded = true
		}
		prevMedium = medium
	}

	return append(straightenedPolyChain, polyChain[len(polyChain)-1])
}

func ApproximatePolyChain(polyChain PolyChain, minDeviation float64) PolyChain {
	if len(polyChain) <= 2 {
		return polyChain
	}

	lineSegment := LineSegment{polyChain[0], polyChain[len(polyChain)-1]}
	segmentLength := Distance(polyChain[0], polyChain[len(polyChain)-1])
	segmentLengthSquare := segmentLength * segmentLength
	maxDistance := minDeviation
	var distance, ratio, maxRatio float64
	var maxI int
	for i := 1; i < len(polyChain)-1; i++ {
		distanceToFirst, distanceToLast := Distance(polyChain[i], polyChain[0]), Distance(polyChain[i], polyChain[len(polyChain)-1])
		if distanceToFirst > distanceToLast {
			distanceToFirst, distanceToLast = distanceToLast, distanceToFirst
		}

		if distanceToFirst*distanceToFirst+segmentLengthSquare <= distanceToLast*distanceToLast {
			distance = distanceToFirst
		} else {
			distance = DistanceToLine(polyChain[i], lineSegment)
		}

		if distance > maxDistance {
			maxDistance, maxI, maxRatio = distance, i, float64(i)/float64(len(polyChain)-i-1)
		} else if distance == maxDistance {
			if i < len(polyChain)-i-1 {
				ratio = float64(i) / float64(len(polyChain)-i-1)
			} else {
				ratio = float64(len(polyChain)-i-1) / float64(i)
			}
			if ratio > maxRatio {
				maxI, maxRatio = i, float64(i)/float64(len(polyChain)-i-1)
			}
		}
	}
	if maxI <= 0 {
		return lineSegment[:]
	}

	return append(ApproximatePolyChain(polyChain[:maxI+1], minDeviation), ApproximatePolyChain(polyChain[maxI:], minDeviation)[1:]...)
}

func CutPolyChain(polyChain PolyChain, endI int, axis LineSegment) PolyChain {

	numI := len(polyChain)

	if endI < 0 || endI >= numI {
		return nil
	}

	p0 := polyChain[endI]
	if DistanceToLine(p0, axis) == 0 {
		return nil
	}

	cutted := []Point2{p0}

	divided := false
	for i := (endI + 1) % numI; i != endI; i = (i + 1) % numI {
		p1 := polyChain[i]
		if DividedByLine(p0, p1, axis) {
			divided = true
			s := LineSegment{polyChain[i-1], p1}
			p11 := LinesIntersection(s, axis)
			if p11 != nil {
				cutted = append(cutted, *p11)
			} else {
				// TODO!!! but why???
				cutted = append(cutted, s.Middle())
			}
			break
		} else {
			cutted = append(cutted, p1)
		}
	}

	if divided {

		i := (endI + numI - 1) % numI
		for {
			p1 := polyChain[i]
			if DividedByLine(p0, p1, axis) {
				s := LineSegment{polyChain[(i+1)%numI], p1}
				p11 := LinesIntersection(s, axis)
				if p11 != nil {
					cutted = append(cutted, *p11)
				} else {
					// TODO!!! but why???
					cutted = append(cutted, s.Middle())
				}
				break
			}
			i = (i + numI - 1) % numI
		}

		i = (i + 1) % numI
		for ; i != endI; i = (i + 1) % numI {
			cutted = append(cutted, polyChain[i])
		}

	}

	return cutted
}

func AveragePolyChains(polyChain0, polyChain1 PolyChain, distanceMaxIn float64) (ok bool, polyChain0Averaged PolyChain, polyChain1Rest []PolyChain) {

	for i0, p0 := range polyChain0 {
		dist, projections := DistanceToPolyChain(polyChain1, p0)

		// log.Print(i0, p0)

		if dist <= distanceMaxIn {

			// log.Print(i0, p0, dist, projections)

			pr := projections[0]
			pAvg := Point2{0.5 * (pr.X + p0.X), 0.5 * (pr.Y + p0.Y)}
			polyChain0[i0] = pAvg
			if pr.Position == 0 {
				polyChain1[pr.N] = pAvg
			} else {
				polyChain1 = append(polyChain1[:pr.N+1], append(PolyChain{pAvg}, polyChain1[pr.N+1:]...)...)
			}
			ok = true
		}
	}

	// log.Fatalf("%v / %v", polyChain0, polyChain1)

	// TODO!!! if polyChain1 order is reversed and projections are multiple

	var nextI1 int
	for i1, p1 := range polyChain1 {
		dist, projections := DistanceToPolyChain(polyChain0, p1)
		if dist <= distanceMaxIn {
			pr := projections[0]
			pAvg := Point2{0.5 * (pr.X + p1.X), 0.5 * (pr.Y + p1.Y)}
			polyChain1[i1] = pAvg
			if pr.Position == 0 {
				polyChain0[pr.N] = pAvg
			} else {
				polyChain0 = append(polyChain0[:pr.N+1], append(PolyChain{pAvg}, polyChain0[pr.N+1:]...)...)
			}
			if i1 > nextI1 {
				i1RestFrom := nextI1
				if i1RestFrom > 0 {
					i1RestFrom--
				}
				polyChain1Rest = append(polyChain1Rest, polyChain1[i1RestFrom:i1+1])
			}
			ok = true
			nextI1 = i1 + 1
		}
	}
	if nextI1 < len(polyChain1) {
		i1RestFrom := nextI1
		if i1RestFrom > 0 {
			i1RestFrom--
		}
		polyChain1Rest = append(polyChain1Rest, polyChain1[i1RestFrom:])
	}

	return ok, polyChain0, polyChain1Rest
}

type ProjectionOnPolyChain struct {
	N        int
	Position float64
	Point2
}

func DistanceToPolyChain(polyChain PolyChain, p Point2) (float64, []ProjectionOnPolyChain) {
	if len(polyChain) < 1 {
		return 0, nil
	}

	minDist, projections := Distance(p, polyChain[0]), []ProjectionOnPolyChain{{N: 0, Position: 0, Point2: polyChain[0]}}
	var p2 Point2
	var n int

POINTS:
	for i, p0 := range polyChain[:len(polyChain)-1] {
		dist, position := DistanceToLineSegment(p, LineSegment{p0, polyChain[i+1]})
		if dist > minDist {
			continue
		}

		if segmentLength := Distance(p0, polyChain[i+1]); segmentLength <= 0 {
			p2, n = p0, i
		} else if position >= segmentLength {
			p2, n, position = polyChain[i+1], i+1, 0
		} else {
			dx, dy := polyChain[i+1].X-p0.X, polyChain[i+1].Y-p0.Y
			p2, n = Point2{p0.X + dx*position/segmentLength, p0.Y + dy*position/segmentLength}, i
		}

		if dist < minDist {
			minDist, projections = dist, []ProjectionOnPolyChain{{n, position, p2}}
		} else {
			for _, pr := range projections {
				if n == pr.N && position == pr.Position {
					continue POINTS
				}
			}

			projections = append(projections, ProjectionOnPolyChain{n, position, p2})
		}
	}

	return minDist, projections
}

func DistanceToLineSegment(p Point2, ls LineSegment) (float64, float64) {
	d0, d1, d := DistanceSquare(p, ls[0]), DistanceSquare(p, ls[1]), DistanceSquare(ls[0], ls[1])
	var reversed bool
	if d1 < d0 {
		d0, d1, ls = d1, d0, LineSegment{ls[1], ls[0]}
		reversed = true
	}
	if d0+d <= d1 {
		if reversed {
			return math.Sqrt(d0), math.Sqrt(d)
		} else {
			return math.Sqrt(d0), 0
		}
	}

	c0 := (d0 + d - d1) / (2 * math.Sqrt(d))

	if reversed {
		return math.Sqrt(d0 - c0*c0), math.Sqrt(d) - c0
	} else {
		return math.Sqrt(d0 - c0*c0), c0
	}
}
