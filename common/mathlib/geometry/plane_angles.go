package geometry

// TODO!!!! be careful: all single point angles are calculated in range -pi < angle <= pi

func Vector(p0, p1 Point2) Point2 {
	return Point2{X: p1.X - p0.X, Y: p1.Y - p0.Y}
}

func RotateAround(p, center Point2, angle Rotation) Point2 {
	pToCenter := Point2{p.X - center.X, p.Y - center.Y}
	pToCenterRotated := RotateByAngle(pToCenter, angle)

	return Point2{pToCenterRotated.X + center.X, pToCenterRotated.Y + center.Y}
}
