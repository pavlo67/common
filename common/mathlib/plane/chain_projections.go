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

		// fmt.Printf("{%g %g} {%g %g} {%g %g} --> %g\n", p.X, p.Y, pI.X, pI.Y, pCh[i+1].X, pCh[i+1].Y, dist)

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
		//		if n == pr.n && position == pr.pos {
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

func EndProjection(pCh PolyChain, start bool) ProjectionOnPolyChain {
	var pr ProjectionOnPolyChain
	if start {
		pr.Point2 = pCh[0]
	} else {
		pr.N, pr.Point2 = len(pCh)-1, pCh[len(pCh)-1]
	}

	return pr
}

func ProjectionBetween(pr0, pr1, pr ProjectionOnPolyChain) bool {
	if pr0.N > pr1.N || (pr0.N == pr1.N && pr0.Position > pr1.Position) {
		pr0, pr1 = pr1, pr0
	}

	return (pr.N > pr0.N || (pr.N == pr0.N && pr.Position >= pr0.Position)) &&
		(pr.N < pr1.N || (pr.N == pr1.N && pr.Position <= pr1.Position))
}

func CutWithProjection(pCh PolyChain, pr ProjectionOnPolyChain, fromStart bool) PolyChain {
	if pr.N < 0 || pr.N >= len(pCh) {
		return nil
	}

	if fromStart {
		if pr.Position == 0 {
			return append(PolyChain{}, pCh[:pr.N+1]...)
		}
		return append(append(PolyChain{}, pCh[:pr.N+1]...), pr.Point2)
	}
	if pr.Position == 0 {
		return append(PolyChain{}, pCh[pr.N:]...)
	}
	return append(PolyChain{pr.Point2}, pCh[pr.N+1:]...)
}

func CutWithProjections(pCh PolyChain, pr0, pr1 ProjectionOnPolyChain) PolyChain {
	if pr0.N < 0 || pr0.N >= len(pCh) || pr1.N < 0 || pr1.N >= len(pCh) {
		return nil
	}

	var reversed bool
	if pr0.N > pr1.N || (pr0.N == pr1.N && pr0.Position > pr1.Position) {
		reversed, pr0, pr1 = true, pr1, pr0

	}

	if pr1.Position == 0 {
		pCh = append(PolyChain{}, pCh[:pr1.N+1]...)
	} else {
		pCh = append(append(PolyChain{}, pCh[:pr1.N+1]...), pr1.Point2)
	}

	if pr0.Position == 0 {
		pCh = append(PolyChain{}, pCh[pr0.N:]...)
	} else {
		pCh = append(PolyChain{pr0.Point2}, pCh[pr0.N+1:]...)
	}
	if reversed {
		return pCh.Reversed()
	}

	return pCh
}

//func DivideByProjection(pCh plane.PolyChain, pr plane.ProjectionOnPolyChain) []plane.PolyChain {
//	if pr.n < 0 || pr.n >= len(pCh) {
//		return nil
//	}
//
//	if pr.pos > 0 {
//		return []plane.PolyChain{
//			append(pCh[:pr.n+1], pr.Point2).Reversed(),
//			append(plane.PolyChain{pr.Point2}, pCh[pr.n+1:]...),
//		}
//	}
//
//	pChs := []plane.PolyChain{pCh[pr.n:]}
//	if pr.n > 0 {
//		pChs = append(pChs, pCh[:pr.n+1].Reversed())
//	}
//	return pChs
//}
