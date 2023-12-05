package plane

import (
	"math"

	"github.com/pavlo67/common/common/mathlib"
)

type ProjectionOnPolyChain struct {
	N        int
	Position float64
	Point2
}

func (p Point2) DistanceToPolyChain(pCh PolyChain) (float64, ProjectionOnPolyChain) {

	if len(pCh) < 1 {
		return math.NaN(), ProjectionOnPolyChain{N: -1, Position: math.NaN(), Point2: Point2{math.NaN(), math.NaN()}}
	} else if len(pCh) == 1 {
		return p.DistanceTo(pCh[0]), ProjectionOnPolyChain{Point2: pCh[0]}
	}

	minDist := math.Inf(1)
	var pr ProjectionOnPolyChain
	var n int
	var pPr Point2

	// POINTS:
	for i, pI := range pCh[:len(pCh)-1] {
		dist, position := p.DistanceToSegment(Segment{pI, pCh[i+1]})
		if dist >= minDist {
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

		//if dist < minDist {
		//	minDist, projections = dist, []ProjectionOnPolyChain{{n, position, pPr}}
		//} else {
		//	for _, pr := range projections {
		//		if n == pr.N && position == pr.Position {
		//			continue POINTS
		//		}
		//	}
		//
		//	projections = append(projections, ProjectionOnPolyChain{n, position, pPr})
		//}
	}

	return minDist, pr
}

func AddProjectionPoint(pCh PolyChain, pr ProjectionOnPolyChain) (PolyChain, ProjectionOnPolyChain, bool) {
	if pr.N < 0 {
		return append(PolyChain{pr.Point2}, pCh...), ProjectionOnPolyChain{Point2: pr.Point2}, true
	} else if pr.N >= len(pCh) || (pr.N == len(pCh)-1 && pr.Position > 0) {
		return append(pCh, pr.Point2), ProjectionOnPolyChain{N: len(pCh), Point2: pr.Point2}, true
	} else if pr.Position == 0 {
		// TODO??? check if pr.Point2 == pCh[pr.N]
		return pCh, pr, false
	}

	return append(pCh[:pr.N+1], append(PolyChain{pr.Point2}, pCh[pr.N+1:]...)...), ProjectionOnPolyChain{N: pr.N + 1, Point2: pr.Point2}, true
}

type ProjectionOnPolyChainDirected struct {
	Distance float64
	Angle    float64
	ProjectionOnPolyChain
}

func ProjectionsOnPolyChain(polyChain PolyChain, p Point2, distanceMax float64) []ProjectionOnPolyChainDirected {
	if len(polyChain) < 1 {
		return nil

	} else if len(polyChain) == 1 {
		if distance := p.DistanceTo(polyChain[0]); distance <= distanceMax {
			return []ProjectionOnPolyChainDirected{{
				Distance:              distance,
				ProjectionOnPolyChain: ProjectionOnPolyChain{Point2: polyChain[0]},
			}}
		}
		return nil

	}

	var projections []ProjectionOnPolyChainDirected

	for n := 0; n < len(polyChain)-1; n++ {
		pn := polyChain[n]
		dist, position := p.DistanceToSegment(Segment{pn, polyChain[n+1]})

		// log.Print(dist, distanceMax)

		if dist > distanceMax {
			continue
		}

		segmLen := pn.DistanceTo(polyChain[n+1])
		if position < segmLen {
			projections = append(projections, ProjectionOnPolyChainDirected{dist, 0, ProjectionOnPolyChain{
				n, position, Point2{pn.X + (polyChain[n+1].X-pn.X)*position/segmLen, pn.Y + (polyChain[n+1].Y-pn.Y)*position/segmLen}}})
		} else if n == len(polyChain)-2 {
			projections = append(projections, ProjectionOnPolyChainDirected{dist, 0, ProjectionOnPolyChain{
				n + 1, 0, polyChain[n+1]}})
		}
	}

	return projections
}

// we check the ray directed from ls[1] to infinity
func SegmentProjectionOnPolyChain(polyChain PolyChain, ls Segment, distanceMax float64) *ProjectionOnPolyChainDirected {
	//log.Printf("%v + %v", polyChain, ls)

	if len(polyChain) < 1 {
		return nil

	}

	segmLen := ls[0].DistanceTo(ls[1])

	if len(polyChain) == 1 {
		distance0, distance1 := polyChain[0].DistanceTo(ls[0]), polyChain[0].DistanceTo(ls[1])
		if math.Abs(segmLen+distance1-distance0) <= mathlib.Eps {
			return &ProjectionOnPolyChainDirected{Distance: distance1, ProjectionOnPolyChain: ProjectionOnPolyChain{Point2: polyChain[0]}}
		}
		return nil
	}

	var pr *ProjectionOnPolyChainDirected

	for n, pn := range polyChain[:len(polyChain)-1] {
		if p := ls.DivideByLine(pn, polyChain[n+1]); p != nil {

			//log.Print(*p)

			if distance1 := p.DistanceTo(ls[1]); distance1 <= distanceMax && (pr == nil || distance1 < pr.Distance) {

				distance0 := p.DistanceTo(ls[0])

				//log.Print(distance0, distance1, segmLen)

				if math.Abs(segmLen+distance1-distance0) <= mathlib.Eps {
					position := p.DistanceTo(pn)
					if position >= pn.DistanceTo(polyChain[n+1]) {
						pr = &ProjectionOnPolyChainDirected{
							Distance:              distance1,
							ProjectionOnPolyChain: ProjectionOnPolyChain{N: n + 1, Point2: *p}}
					} else {
						pr = &ProjectionOnPolyChainDirected{
							Distance:              distance1,
							ProjectionOnPolyChain: ProjectionOnPolyChain{N: n, Position: position, Point2: *p}}
					}
				}
			}
		}
	}

	return pr
}
