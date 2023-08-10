package geometry

import (
	"math"

	"github.com/pavlo67/common/common/mathlib"
)

type ProjectionOnPolyChain struct {
	N        int
	Position float64
	Point2
}

func DistanceToPolyChain(p Point2, pCh PolyChain) (float64, ProjectionOnPolyChain) {

	if len(pCh) < 1 {
		return math.NaN(), ProjectionOnPolyChain{N: -1, Position: math.NaN(), Point2: Point2{math.NaN(), math.NaN()}}
	} else if len(pCh) == 1 {
		return Distance(p, pCh[0]), ProjectionOnPolyChain{Point2: pCh[0]}
	}

	minDist := math.Inf(1)
	var pr ProjectionOnPolyChain
	var n int
	var pPr Point2

	// POINTS:
	for i, pI := range pCh[:len(pCh)-1] {
		dist, position := DistanceToLineSegment(p, LineSegment{pI, pCh[i+1]})
		if dist >= minDist {
			continue
		}

		if segmentLength := Distance(pI, pCh[i+1]); segmentLength <= 0 {
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

const onDivideByProjection = "on AddProjectionPoint()"

func AddProjectionPoint(pCh PolyChain, pr ProjectionOnPolyChain) (PolyChain, ProjectionOnPolyChain) {
	if pr.N < 0 {
		return append(PolyChain{pr.Point2}, pCh...), ProjectionOnPolyChain{Point2: pr.Point2}
	} else if pr.N >= len(pCh) || (pr.N == len(pCh)-1 && pr.Position > 0) {
		return append(pCh, pr.Point2), ProjectionOnPolyChain{N: len(pCh), Point2: pr.Point2}
	} else if pr.Position == 0 {
		// TODO??? check if pr.Point2 == pCh[pr.N]
		return pCh, pr
	}

	return append(pCh[:pr.N+1], append(PolyChain{pr.Point2}, pCh[pr.N+1:]...)...), ProjectionOnPolyChain{N: pr.N + 1, Point2: pr.Point2}
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
		if distance := Distance(p, polyChain[0]); distance <= distanceMax {
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
		dist, position := DistanceToLineSegment(p, LineSegment{pn, polyChain[n+1]})

		// log.Print(dist, distanceMax)

		if dist > distanceMax {
			continue
		}

		segmLen := Distance(pn, polyChain[n+1])
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
func SegmentProjectionOnPolyChain(polyChain PolyChain, ls LineSegment, distanceMax float64) *ProjectionOnPolyChainDirected {
	//log.Printf("%v + %v", polyChain, ls)

	if len(polyChain) < 1 {
		return nil

	}

	segmLen := Distance(ls[0], ls[1])

	if len(polyChain) == 1 {
		distance0, distance1 := Distance(polyChain[0], ls[0]), Distance(polyChain[0], ls[1])
		if math.Abs(segmLen+distance1-distance0) <= mathlib.Eps {
			return &ProjectionOnPolyChainDirected{Distance: distance1, ProjectionOnPolyChain: ProjectionOnPolyChain{Point2: polyChain[0]}}
		}
		return nil
	}

	var pr *ProjectionOnPolyChainDirected

	for n, pn := range polyChain[:len(polyChain)-1] {
		if p := DividedByLine(pn, polyChain[n+1], ls); p != nil {

			//log.Print(*p)

			if distance1 := Distance(*p, ls[1]); distance1 <= distanceMax && (pr == nil || distance1 < pr.Distance) {

				distance0 := Distance(*p, ls[0])

				//log.Print(distance0, distance1, segmLen)

				if math.Abs(segmLen+distance1-distance0) <= mathlib.Eps {
					position := Distance(pn, *p)
					if position >= Distance(pn, polyChain[n+1]) {
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
