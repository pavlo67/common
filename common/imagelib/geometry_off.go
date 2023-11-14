package imagelib

//func Direction(el1, el2 image.Point) float64 {
//	dx := float64(el2.X) - float64(el1.X)
//	dy := float64(el2.Y) - float64(el1.Y)
//
//	if dx == 0 {
//		if dy > 0 {
//			return 0
//		} else if dy < 0 {
//			return 180
//		} else {
//			math.NaN()
//		}
//	}
//
//	direction_ := 180 * math.Atan(dy/dx) / math.Pi
//	if dx > 0 {
//		return direction_
//	} else if dy > 0 {
//		return 180 + direction_
//	}
//
//	return -180 + direction_
//}
//

//func AverageAlongOx(points2 []numlib.Point2) ([]image.ImagePoint, image.Rect) {
//	if len(points2) < 1 {
//		return nil, image.Rect{}
//	}
//
//	sort.Slice(points2, func(i, j int) bool { return points2[i].Position < points2[j].Position })
//	pX := int(points2[0].Position)
//	if points2[0].Position < 0 {
//		pX--
//	}
//	xBase := -pX
//	yBase := 0
//	var yBaseI int
//	for _, p := range points2 {
//		if p.Y >= 0 {
//			yBaseI = int(p.Y)
//		} else {
//			yBaseI = -int(p.Y)
//		}
//
//		if yBaseI > yBase {
//			yBase = yBaseI
//		}
//	}
//
//	points := make([]image.ImagePoint, len(points2))
//
//	var yPlus, yMinus []float64
//	for _, p := range points2 {
//		pXNext := int(p.Position)
//		if p.Position < 0 {
//			pXNext--
//		}
//		if pXNext != pX {
//			points = append(points, AveragedPoint(pX+xBase, yBase+1, yPlus, yMinus))
//			yPlus, yMinus, pX = nil, nil, pXNext
//		}
//
//		if p.Y > 0 {
//			yPlus = append(yPlus, p.Y)
//		} else if p.Y < 0 {
//			yMinus = append(yMinus, p.Y)
//		} else {
//			yPlus = append(yPlus, p.Y)
//			yMinus = append(yMinus, p.Y)
//		}
//	}
//	if len(yPlus)+len(yMinus) > 0 {
//		points = append(points, AveragedPoint(pX+xBase, yBase, yPlus, yMinus))
//	}
//	toX := points[len(points)-1].Position + 1
//
//	for i := len(points) - 1; i >= 0; i-- {
//		points = append(points, image.ImagePoint{Position: points[i].Position, Y: yBase - points[i].Y})
//	}
//
//	return points, image.Rect{Max: image.ImagePoint{toX, yBase * 2}}
//}
//
//func AveragedPoint(x, yBase int, yPlus, yMinus []float64) image.ImagePoint {
//	var yPlusAvg, yMinusAvg float64
//
//	if len(yPlus) > 0 {
//		for _, y := range yPlus {
//			yPlusAvg += y
//		}
//		yPlusAvg /= float64(len(yPlus))
//	}
//
//	if len(yMinus) > 0 {
//		for _, y := range yMinus {
//			yMinusAvg += y
//		}
//		yMinusAvg /= float64(len(yMinus))
//	}
//
//	return image.ImagePoint{x, yBase + int(math.Round((yPlusAvg-yMinusAvg)/2))}
//}
