package geometry

// Rotation is angle in radians (-Pi < rotation <= Pi), rotation from Ox to Oy (counter clockwise) has a positive angle
type Rotation float64

// Position moves geometry shape to Point2 and (after!) rotates it with Rotation angle around itself
type Position struct {
	Point2
	Rotation
}

//type Bearing struct {
//	Rotation          // of the coordinate system
//	Rotation    Rotation // in the coordinate system
//}
//
//func (bearing Bearing) Projections(r float64) (x, y float64) {
//	angle := float64(bearing.Rotation + bearing.Rotation)
//	return r * math.Cos(angle), r * math.Sin(angle)
//}
