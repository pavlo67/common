package geometry

import (
	"sort"

	"github.com/pavlo67/common/common/mathlib"
)

type PolyChain []Point2

func (pCh PolyChain) Length() float64 {
	length := 0.
	for i := 1; i < len(pCh); i++ {
		length += Distance(pCh[i-1], pCh[i])
	}
	return length
}

func (pCh PolyChain) Reversed() PolyChain {
	polyChainReversed := make([]Point2, len(pCh))
	for i, p2 := range pCh {
		polyChainReversed[len(pCh)-i-1] = p2
	}
	return polyChainReversed
}

func (pCh PolyChain) Direction(deviationMaxIn float64) *LineSegment {
	if len(pCh) < 1 {
		return nil
	} else if len(pCh) < 2 {
		return &LineSegment{pCh[0], pCh[0]}
	}

	pEnd := pCh[len(pCh)-1]
	directionLine := LineSegment{pCh[len(pCh)-2], pEnd}

	for i := len(pCh) - 3; i >= 0; i-- {
		directionLineNext := LineSegment{LineSegment{pCh[i], pCh[i+1]}.Middle(), pEnd}
		for j := i + 1; j < len(pCh)-1; j++ {
			if DistanceToLine(pCh[j], directionLineNext) > deviationMaxIn {
				return &directionLine
			}
		}
		directionLine = directionLineNext
	}

	return &directionLine
}

func ShortenPolyChain(pCh PolyChain, maxDistanceToBeJoined float64) PolyChain {
	for i := 0; i <= len(pCh)-3; i++ {
		for j := len(pCh) - 1; j >= i+2; j-- {
			if Distance(pCh[i], pCh[j]) <= maxDistanceToBeJoined {
				pCh = append(pCh[:i+1], pCh[j:]...)
				break
			}
		}
	}

	return pCh
}

func StraightenPolyChain(pCh PolyChain, minDeviation float64) PolyChain {
	if len(pCh) <= 2 {
		return pCh
	}

	straightenedPolyChain := PolyChain{pCh[0]}
	prevMedium := LineSegment{pCh[0], pCh[1]}.Middle()
	var prevMediumAdded bool

	for i := 1; i < len(pCh)-1; i++ {
		medium := LineSegment{pCh[i], pCh[i+1]}.Middle()
		if distance := DistanceToLine(pCh[i], LineSegment{prevMedium, medium}); distance > minDeviation {
			straightenedPolyChain = append(straightenedPolyChain, pCh[i])
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

	return append(straightenedPolyChain, pCh[len(pCh)-1])
}

func ApproximatePolyChain(pCh PolyChain, minDeviation float64) PolyChain {
	if len(pCh) <= 2 {
		return pCh
	}

	lineSegment := LineSegment{pCh[0], pCh[len(pCh)-1]}
	segmentLength := Distance(pCh[0], pCh[len(pCh)-1])
	segmentLengthSquare := segmentLength * segmentLength
	maxDistance := minDeviation
	var distance, ratio, maxRatio float64
	var maxI int
	for i := 1; i < len(pCh)-1; i++ {
		distanceToFirst, distanceToLast := Distance(pCh[i], pCh[0]), Distance(pCh[i], pCh[len(pCh)-1])
		if distanceToFirst > distanceToLast {
			distanceToFirst, distanceToLast = distanceToLast, distanceToFirst
		}

		if distanceToFirst*distanceToFirst+segmentLengthSquare <= distanceToLast*distanceToLast {
			distance = distanceToFirst
		} else {
			distance = DistanceToLine(pCh[i], lineSegment)
		}

		if distance > maxDistance {
			maxDistance, maxI, maxRatio = distance, i, float64(i)/float64(len(pCh)-i-1)
		} else if distance == maxDistance {
			if i < len(pCh)-i-1 {
				ratio = float64(i) / float64(len(pCh)-i-1)
			} else {
				ratio = float64(len(pCh)-i-1) / float64(i)
			}
			if ratio > maxRatio {
				maxI, maxRatio = i, float64(i)/float64(len(pCh)-i-1)
			}
		}
	}
	if maxI <= 0 {
		return lineSegment[:]
	}

	return append(ApproximatePolyChain(pCh[:maxI+1], minDeviation), ApproximatePolyChain(pCh[maxI:], minDeviation)[1:]...)
}

func CutPolyChain(pCh PolyChain, fromEndI int, axis LineSegment) PolyChain {

	numI := len(pCh)

	if fromEndI < 0 || fromEndI >= numI {
		return nil
	}

	p0 := pCh[fromEndI]
	if DistanceToLine(p0, axis) == 0 {
		return nil
	}

	cutted := []Point2{p0}

	divided := false
	for i := (fromEndI + 1) % numI; i != fromEndI; i = (i + 1) % numI {
		p1 := pCh[i]
		if DividedByLine(p0, p1, axis) != nil {
			divided = true
			s := LineSegment{pCh[i-1], p1}
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

		i := (fromEndI + numI - 1) % numI
		for {
			p1 := pCh[i]
			if DividedByLine(p0, p1, axis) != nil {
				s := LineSegment{pCh[(i+1)%numI], p1}
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
		for ; i != fromEndI; i = (i + 1) % numI {
			cutted = append(cutted, pCh[i])
		}

	}

	return cutted
}

type ProjectionPolyChainOnPolyChain struct {
	N0 int
	ProjectionOnPolyChain
}

func AveragePolyChains(pCh0, pCh1 PolyChain, distanceMaxIn float64) (ok bool, polyChain0Averaged PolyChain, polyChain1Rest []PolyChain) {

	var projections []ProjectionPolyChainOnPolyChain

	for i0, p0 := range pCh0 {
		if dist, pr := DistanceToPolyChain(p0, pCh1); dist <= distanceMaxIn {
			projections = append(projections, ProjectionPolyChainOnPolyChain{
				N0:                    i0,
				ProjectionOnPolyChain: pr,
			})
		}
	}

	var pCh1Averaged []int

	if len(projections) > 0 {
		sort.Slice(projections, func(i, j int) bool {
			return projections[i].N > projections[j].N || (projections[i].N == projections[j].N && projections[i].Position > projections[j].Position)
		})

		ok = true
		for _, pr := range projections {
			pAvg := Point2{0.5 * (pr.X + pCh0[pr.N0].X), 0.5 * (pr.Y + pCh0[pr.N0].Y)}
			if pr.Position == 0 {
				// TODO!!! if duplicated???
				pCh1[pr.N] = pAvg
				pCh1Averaged = append(pCh1Averaged, pr.N)
			} else {
				pCh1 = append(pCh1[:pr.N+1], append(PolyChain{pAvg}, pCh1[pr.N+1:]...)...)
				for i := range pCh1Averaged {
					pCh1Averaged[i]++
				}
				pCh1Averaged = append(pCh1Averaged, pr.N+1)
			}
			pCh0[pr.N0] = pAvg
		}
		projections = nil
	}

	// log.Fatalf("%v / %v", pCh0, pCh1)
	// TODO!!! if pCh1 order is reversed and projections are multiple

	for i1, p1 := range pCh1 {
		if mathlib.InInt(pCh1Averaged, i1) {
			continue
		}

		// TODO: be careful if pr.Position == 0 - this projection point on pCh0 must be averaged above
		if dist, pr := DistanceToPolyChain(p1, pCh0); dist <= distanceMaxIn && pr.Position > 0 {
			projections = append(projections, ProjectionPolyChainOnPolyChain{
				N0:                    i1,
				ProjectionOnPolyChain: pr,
			})
		}
	}

	if len(projections) > 0 {
		ok = true
		for _, pr := range projections {
			pAvg := Point2{0.5 * (pr.X + pCh1[pr.N0].X), 0.5 * (pr.Y + pCh1[pr.N0].Y)}
			pCh1[pr.N0] = pAvg
			pCh0 = append(pCh0[:pr.N+1], append(PolyChain{pAvg}, pCh0[pr.N+1:]...)...)
			pCh1Averaged = append(pCh1Averaged, pr.N0)
		}
	}

	sort.Slice(pCh1Averaged, func(i, j int) bool { return pCh1Averaged[i] < pCh1Averaged[j] })

	var n1Prev, n1Next int
	for _, n1 := range pCh1Averaged {
		if n1 > n1Next {
			polyChain1Rest = append(polyChain1Rest, pCh1[n1Prev:n1+1])
		}
		n1Prev, n1Next = n1, n1+1
	}
	if len(pCh1) > n1Next {
		polyChain1Rest = append(polyChain1Rest, pCh1[n1Prev:])
	}

	// log.Print(pCh0, polyChain1Rest, pCh1Averaged)

	return ok, pCh0, polyChain1Rest
}
