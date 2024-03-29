package plane

import (
	"fmt"
	"math"
	"slices"

	"github.com/pavlo67/common/common/mathlib/sets"
)

type PolyChain []Point2

//type ProjectionPolyChainOnPolyChain struct {
//	N0 int
//	ProjectionOnPolyChain
//}

func (pCh PolyChain) Length() float64 {
	length := 0.
	for i := 1; i < len(pCh); i++ {
		length += pCh[i-1].DistanceTo(pCh[i])
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

func (pCh PolyChain) Direction(deviationMaxIn float64) *Segment {
	if len(pCh) < 2 {
		return nil
	}

	pEnd := pCh[len(pCh)-1]
	directionLine := Segment{pCh[len(pCh)-2], pEnd}

	for i := len(pCh) - 3; i >= 0; i-- {
		directionLineNext := Segment{Segment{pCh[i], pCh[i+1]}.Middle(), pEnd}
		for j := i + 1; j < len(pCh)-1; j++ {
			if pCh[j].DistanceToLine(directionLineNext) > deviationMaxIn {
				return &directionLine
			}
		}
		directionLine = directionLineNext

		directionLineNext = Segment{pCh[i], pEnd}
		for j := i + 1; j < len(pCh)-1; j++ {
			if pCh[j].DistanceToLine(directionLineNext) > deviationMaxIn {
				return &directionLine
			}
		}
		directionLine = directionLineNext
	}

	if directionLine[0] == directionLine[1] {
		return nil
	}

	return &directionLine
}

func (p Point2) ProjectionOnLineSegment(ls Segment) (distance, projectionPosition float64) {
	d0, d1, d := p.DistanceSquare(ls[0]), p.DistanceSquare(ls[1]), ls[0].DistanceSquare(ls[1])
	var reversed bool
	if d1 < d0 {
		d0, d1 = d1, d0
		reversed = true
	}
	if d0+d < d1 {
		return math.NaN(), math.NaN()
	}

	c0 := (d0 + d - d1) / (2 * math.Sqrt(d))

	if reversed {
		return math.Sqrt(d0 - c0*c0), math.Sqrt(d) - c0
	} else {
		return math.Sqrt(d0 - c0*c0), c0
	}
}

func (p Point2) ProjectionOnPolyChain(pCh PolyChain) (float64, ProjectionOnPolyChain) {

	if len(pCh) < 1 {
		return math.NaN(), ProjectionOnPolyChain{N: -1, Position: math.NaN(), Point2: Point2{math.NaN(), math.NaN()}}
	} else if len(pCh) == 1 {
		return p.DistanceTo(pCh[0]), ProjectionOnPolyChain{Point2: pCh[0]}
	}

	minDist := math.Inf(1)
	pr := ProjectionOnPolyChain{
		N:        -1,
		Position: math.NaN(),
		Point2:   Point2{math.NaN(), math.NaN()},
	}
	var n int
	var pPr Point2

	// POINTS:
	for i, pI := range pCh[:len(pCh)-1] {
		dist, position := p.ProjectionOnLineSegment(Segment{pI, pCh[i+1]})
		if dist >= minDist || math.IsNaN(dist) {
			continue
		}

		if segmentLength := pI.DistanceTo(pCh[i+1]); segmentLength <= 0 {
			pPr, n, position = pI, i, 0
		} else if position >= segmentLength {
			pPr, n, position = pCh[i+1], i+1, 0
		} else {
			dx, dy := pCh[i+1].X-pI.X, pCh[i+1].Y-pI.Y
			pPr, n = Point2{pI.X + dx*position/segmentLength, pI.Y + dy*position/segmentLength}, i
		}

		minDist, pr = dist, ProjectionOnPolyChain{n, position, pPr}
	}

	return minDist, pr
}

func (pCh PolyChain) DistanceTo(pCh1 PolyChain, distanceMax float64) (dist float64, pr, pr1 *ProjectionOnPolyChain) {
	for n0, p0 := range pCh {
		if dist_, pr_ := p0.DistanceToPolyChain(pCh1); dist_ <= distanceMax {
			return dist_, &ProjectionOnPolyChain{N: n0, Point2: p0}, &pr_
		}
	}
	for n1, p1 := range pCh1 {
		if dist_, pr_ := p1.DistanceToPolyChain(pCh); dist_ <= distanceMax {
			return dist_, &pr_, &ProjectionOnPolyChain{N: n1, Point2: p1}
		}
	}

	return math.NaN(), nil, nil
}

func AveragePolyChains(pCh0, pCh1 PolyChain, distanceMaxIn float64, connectEnds bool) (
	ok bool, pCh0Averaged PolyChain, pCh1RestsInitial []PolyChain) {

	// log.Print(distanceMaxIn, pCh0, pCh1)
	// if distanceMaxIn == 9.87053098413958 {
	// }

	var p1Averaged []int

	for n0, p0 := range pCh0 {
		if dist, pr := p0.ProjectionOnPolyChain(pCh1); dist <= distanceMaxIn {
			pCh0[n0] = Point2{0.5 * (pr.X + p0.X), 0.5 * (pr.Y + p0.Y)}
			if pr.Position == 0 {
				pCh1[pr.N] = pCh0[n0]
				p1Averaged = append(p1Averaged, pr.N)
			} else {
				// TODO!!! be careful: appending order is essential here
				pCh1 = append(pCh1[:pr.N+1], append(PolyChain{pCh0[n0]}, pCh1[pr.N+1:]...)...)

				for i, p := range p1Averaged {
					if p > pr.N {
						p1Averaged[i]++
					}
				}
				p1Averaged = append(p1Averaged, pr.N+1)
			}
		}
	}

	for n1, p1 := range pCh1 {
		if sets.In(p1Averaged, n1) {
			continue
		}

		if dist, pr := p1.ProjectionOnPolyChain(pCh0); dist <= distanceMaxIn {
			pCh1[n1] = Point2{0.5 * (pr.X + p1.X), 0.5 * (pr.Y + p1.Y)}
			if pr.Position == 0 {
				pCh0[pr.N] = pCh1[n1]
			} else {
				// TODO!!! be careful: appending order is essential here
				pCh0 = append(pCh0[:pr.N+1], append(PolyChain{pCh1[n1]}, pCh0[pr.N+1:]...)...)

			}
			p1Averaged = append(p1Averaged, n1)
		}
	}

	slices.Sort(p1Averaged)

	// log.Print(pCh0, "\n", pCh1, "\n", p1Averaged)

	var n1Prev, n1Next int
	for _, n1 := range p1Averaged {
		if n1 > n1Next {
			pCh1RestsInitial = append(pCh1RestsInitial, pCh1[n1Prev:n1+1])
		}
		n1Prev, n1Next = n1, n1+1
	}
	if len(pCh1) > n1Next {
		pCh1RestsInitial = append(pCh1RestsInitial, pCh1[n1Prev:])
	}

	ok = len(p1Averaged) > 0

	if ok && connectEnds && len(pCh1RestsInitial) > 0 {
		var pCh1Rests []PolyChain
		p00, p01 := pCh0[0], pCh0[len(pCh0)-1]
		v0 := Point2{p01.X - p00.X, p01.Y - p00.Y}
		for _, pCh1Rest := range pCh1RestsInitial {
			p10, p11 := pCh1Rest[0], pCh1Rest[len(pCh1Rest)-1]
			if math.Abs(v0.AnglesDelta(Point2{p11.X - p10.X, p11.Y - p10.Y})) <= math.Pi/2 {
				if p01 == p10 {
					pCh0 = append(pCh0, pCh1Rest[1:]...)
				} else if p11 == p00 {
					pCh0 = append(pCh1Rest, pCh0[1:]...)
				} else {
					pCh1Rests = append(pCh1Rests, pCh1Rest)
				}
			} else {
				if p01 == p11 {
					pCh0 = append(pCh0, pCh1Rest.Reversed()[1:]...)
				} else if p00 == p10 {
					pCh0 = append(pCh0.Reversed(), pCh1Rest[1:]...)
				} else {
					pCh1Rests = append(pCh1Rests, pCh1Rest)
				}

			}
		}

		if len(pCh1Rests) != len(pCh1RestsInitial) {
			// so some non-averaged fragment of pCh1 is added to pCh0
			// TODO!!! be careful: it looks like this operation "reverses" the original pCh0/pCh1 order because pCh1Rests contains pCh0 as the last fragment
			return true, nil, append(pCh1Rests, pCh0)
		}

		return true, pCh0, pCh1Rests
	}

	// if distanceMaxIn == 9.87053098413958 {
	// if ok {
	//	log.Print(2222222222222, distanceMaxIn, pCh0, pCh1RestsInitial)
	// }

	return ok, pCh0, pCh1RestsInitial
}

type PolyChainsIntersection struct {
	Point2
	N0, N1 int
}

func PolyChainsIntersectionAny(pCh0, pCh1 PolyChain) *PolyChainsIntersection {
	for i0 := 1; i0 < len(pCh0); i0++ {
		s0 := Segment{pCh0[i0-1], pCh0[i0]}
		for i1 := 1; i1 < len(pCh1); i1++ {
			if p := SegmentsIntersection(s0, Segment{pCh1[i1-1], pCh1[i1]}); p != nil {
				return &PolyChainsIntersection{*p, i0 - 1, i1 - 1}
			}
		}
	}

	return nil
}

// it's equal to append(append(...

//pCh1New := make(PolyChain, len(pCh1)+1)
//for i := 0; i <= pr.N; i++ {
//	pCh1New[i] = pCh1[i]
//}
//pCh1New[pr.N+1] = pCh0[n0]
//for i := pr.N + 1; i < len(pCh1); i++ {
//	pCh1New[i+1] = pCh1[i]
//}
//pCh1 = pCh1New

//pCh0New := make(PolyChain, len(pCh0)+1)
//for i := 0; i <= pr.N; i++ {
//	pCh0New[i] = pCh0[i]
//}
//pCh0New[pr.N+1] = pCh1[n1]
//for i := pr.N + 1; i < len(pCh0); i++ {
//	pCh0New[i+1] = pCh0[i]
//}
//pCh0 = pCh0New

func (pCh PolyChain) Shorten(maxDistanceToBeJoined float64) PolyChain {
	for i := 0; i <= len(pCh)-3; i++ {
		for j := len(pCh) - 1; j >= i+2; j-- {
			if pCh[i].DistanceTo(pCh[j]) <= maxDistanceToBeJoined {
				pCh = append(pCh[:i+1], pCh[j:]...)
				break
			}
		}
	}

	return pCh
}

func (pCh PolyChain) Straighten(minDeviation float64) PolyChain {
	if len(pCh) <= 2 {
		return pCh
	}

	straightenedPolyChain := PolyChain{pCh[0]}
	prevMedium := Segment{pCh[0], pCh[1]}.Middle()
	var prevMediumAdded bool

	for i := 1; i < len(pCh)-1; i++ {
		medium := Segment{pCh[i], pCh[i+1]}.Middle()
		if distance := pCh[i].DistanceToLine(Segment{prevMedium, medium}); distance > minDeviation {
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

func (pCh PolyChain) Cut(fromEndI int, axis Segment) PolyChain {

	numI := len(pCh)

	if fromEndI < 0 || fromEndI >= numI {
		return nil
	}

	p0 := pCh[fromEndI]
	if p0.DistanceToLine(axis) == 0 {
		return nil
	}

	cutted := []Point2{p0}

	divided := false
	for i := (fromEndI + 1) % numI; i != fromEndI; i = (i + 1) % numI {
		p1 := pCh[i]
		if axis.DivideByLine(p0, p1) != nil {
			divided = true
			s := Segment{pCh[i-1], p1}
			p11 := s.LinesIntersection(axis)
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
			if axis.DivideByLine(p0, p1) != nil {
				s := Segment{pCh[(i+1)%numI], p1}
				p11 := s.LinesIntersection(axis)
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

func (pCh PolyChain) Approximate(minDeviation float64) PolyChain {
	if len(pCh) <= 2 {
		return pCh
	}

	lineSegment := Segment{pCh[0], pCh[len(pCh)-1]}
	segmentLength := pCh[0].DistanceTo(pCh[len(pCh)-1])
	segmentLengthSquare := segmentLength * segmentLength
	maxDistance := minDeviation
	var distance, ratio, maxRatio float64
	var maxI int
	for i := 1; i < len(pCh)-1; i++ {
		distanceToFirst, distanceToLast := pCh[i].DistanceTo(pCh[0]), pCh[i].DistanceTo(pCh[len(pCh)-1])
		if distanceToFirst > distanceToLast {
			distanceToFirst, distanceToLast = distanceToLast, distanceToFirst
		}

		if distanceToFirst*distanceToFirst+segmentLengthSquare <= distanceToLast*distanceToLast {
			distance = distanceToFirst
		} else {
			distance = pCh[i].DistanceToLine(lineSegment)
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

	return append(append(PolyChain{}, pCh[:maxI+1].Approximate(minDeviation)...), pCh[maxI:].Approximate(minDeviation)[1:]...)
}

func (pCh PolyChain) Filter() PolyChain {
	if len(pCh) < 2 {
		return pCh
	}
	var pChNew PolyChain

I:
	for i := 0; i < len(pCh); {
		pChNew = append(pChNew, pCh[i])
		for j := i + 1; j < len(pCh); j++ {
			if pCh[j] != pCh[i] {
				i = j
				continue I
			}
		}
		break
	}

	return pChNew
}

func (pCh PolyChain) ShortenWithTheSamePoints(maxDistanceToBeJoined float64) PolyChain {
	for i := 0; i <= len(pCh)-3; i++ {
		for j := len(pCh) - 1; j >= i+2; j-- {
			if pCh[i].DistanceTo(pCh[j]) <= maxDistanceToBeJoined {
				for k := i + 1; k < j; k++ {
					pCh[k] = pCh[i]
				}
				break
			}
		}
	}

	return pCh
}

func (pCh PolyChain) ShortString() string {
	if len(pCh) < 1 {
		return "[]"
	}
	var pChStr string
	for _, p := range pCh {
		pChStr += fmt.Sprintf(" {%.1f %.1f}", p.X, p.Y)
	}

	return "[" + pChStr[1:] + "]"
}
