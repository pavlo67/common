package plane

type RectangleFixed struct {
	P0, P1 Point2
}

func (rectFixed RectangleFixed) Contains(p2 Point2) bool {
	minX, maxX, minY, maxY := rectFixed.P0.X, rectFixed.P1.X, rectFixed.P0.Y, rectFixed.P1.Y
	if minX > maxX {
		minX, maxX = maxX, minX
	}
	if minY > maxY {
		minY, maxY = maxY, minY
	}

	return p2.X >= minX && p2.X <= maxX && p2.Y >= minY && p2.Y <= maxY
}

type Rectangle struct {
	Position             // of the center
	HalfSideX, HalfSideY float64
}

func (rect Rectangle) Contains(p2 Point2) bool {
	p2Rotated := p2.RotateAround(rect.Point2, -rect.Rotation)
	return p2Rotated.X >= rect.Point2.X-rect.HalfSideX && p2Rotated.X <= rect.Point2.X+rect.HalfSideX &&
		p2Rotated.Y >= rect.Point2.Y-rect.HalfSideY && p2Rotated.Y <= rect.Point2.Y+rect.HalfSideY
}

func (rect Rectangle) Outer(margin float64) Rectangle {
	return Rectangle{
		Position:  rect.Position,
		HalfSideX: rect.HalfSideX + margin,
		HalfSideY: rect.HalfSideY + margin,
	}
}

//func (rect Rectangle) Intersection(pCh PolyChain) PolyChain {
//
//	log.Fatal("on Rectangle.Intersects()")
//
//	//for _, p := range pCh {
//	//	p2Rot := RotateByAngle(p, -rect.Rotation)
//	//	if p2Rot.XT >= rect.Min.XT && p2Rot.XT <= rect.MaxIn.XT && p2Rot.YT >= rect.Min.YT && p2Rot.YT <= rect.MaxIn.YT {
//	//		return true
//	//	}
//	//}
//	return nil
//}
//
//func (rect Rectangle) IntersectionArea(rect1 Rectangle) float64 {
//
//	log.Fatal("on Rectangle.IntersectionArea()")
//
//	return 0
//}

func (rect Rectangle) Points() (p00, p01, p10, p11 Point2) {
	if rect.HalfSideX == 0 && rect.HalfSideY == 0 {
		return rect.Point2, rect.Point2, rect.Point2, rect.Point2
	}

	p00Fixed, p01Fixed := Point2{-rect.HalfSideX, -rect.HalfSideY}, Point2{-rect.HalfSideX, rect.HalfSideY}
	p00_, p01_ := p00Fixed.RotateByAngle(rect.Rotation), p01Fixed.RotateByAngle(rect.Rotation)

	return Point2{p00_.X + rect.Point2.X, p00_.Y + rect.Point2.Y},
		Point2{p01_.X + rect.Point2.X, p01_.Y + rect.Point2.Y},
		Point2{-p01_.X + rect.Point2.X, -p01_.Y + rect.Point2.Y},
		Point2{-p00_.X + rect.Point2.X, -p00_.Y + rect.Point2.Y}
}

func (rect Rectangle) OuterFixed(margin float64) RectangleFixed {
	p00, p01, p10, p11 := rect.Points()

	minX, maxX, minY, maxY := p11.X, p11.X, p11.Y, p11.Y
	for _, p := range []Point2{p00, p01, p10} {
		if p.X >= maxX {
			maxX = p.X
		} else if p.X < minX {
			minX = p.X
		}
		if p.Y >= maxY {
			maxY = p.Y
		} else if p.Y < minY {
			minY = p.Y
		}
	}

	return RectangleFixed{Point2{minX - margin, minY - margin}, Point2{maxX + margin, maxY + margin}}
}
