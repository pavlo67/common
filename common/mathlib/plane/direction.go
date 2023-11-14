package plane

func Middle(p0, p1 Point2) Point2 {
	return Point2{(p0.X + p1.X) / 2, (p0.Y + p1.Y) / 2}
}

func PointOnRay(p0 Point2, r float64) Point2 {
	r0 := p0.Radius()
	if r0 == 0 {
		return Point2{}
	}

	return Point2{p0.X * r / r0, p0.Y * r / r0}
}

func Direction(pCh PolyChain, deviationMaxIn float64) (*Segment, int) {
	if len(pCh) < 2 {
		return nil, 0
	}

	pEnd := pCh[len(pCh)-1]
	directionLine := Segment{pCh[len(pCh)-2], pEnd}
	n := 1

	for i := len(pCh) - 3; i >= 0; i-- {
		directionLineNext := Segment{Segment{pCh[i], pCh[i+1]}.Middle(), pEnd}
		for j := i + 1; j < len(pCh)-1; j++ {
			if pCh[j].DistanceToLine(directionLineNext) > deviationMaxIn {
				return &directionLine, n
			}
		}
		directionLine = directionLineNext

		directionLineNext = Segment{pCh[i], pEnd}
		for j := i + 1; j < len(pCh)-1; j++ {
			if pCh[j].DistanceToLine(directionLineNext) > deviationMaxIn {
				return &directionLine, n
			}
		}
		directionLine = directionLineNext
		n++

	}

	if directionLine[0] == directionLine[1] {
		return nil, 0
	}

	return &directionLine, n
}
