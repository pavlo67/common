package _off

//// Marker -------------------------------------------------------------------------------
//
//type Marker interface {
//	Mark(draw.Image, color.Color)
//}
//
//var _ Marker = &MarkerText{}
//var _ Marker = MarkerText{}
//
//type MarkerText struct {
//	DPI      float64
//	Size     float64
//	Spacing  float64
//	FontFile string
//	Text     []string
//	image.Point
//	// logger.Operator
//}
//
//func (mt MarkerText) Mark(drawImg draw.Image, clr color.Color) {
//	if _, err := Write(drawImg, mt.Point, mt.DPI, mt.Size, mt.Spacing, mt.FontFile, clr, mt.Text); err != nil {
//		log.Printf("ERROR: on MarkerText.Mark(): %s", err)
//	}
//}
//
//func Line(s plane.Segment, width int) []image.Point {
//	begin, end := s[0].ImagePoint(), s[1].ImagePoint()
//
//	var wMin int
//	if width <= 1 {
//		width = 1
//	} else {
//		wMin = -(width / 2)
//	}
//
//	if begin == end {
//		points := make([]image.Point, width*width)
//		for x := 0; x < width; x++ {
//			for y := 0; y < width; y++ {
//				points[x*width+y] = image.Point{begin.X + wMin + x, begin.Y + wMin + y}
//			}
//		}
//		return points
//	}
//
//	if math.Abs(float64(end.X-begin.X)) >= math.Abs(float64(end.Y-begin.Y)) {
//		if begin.X > end.X {
//			begin, end = end, begin
//		}
//		deltaX := end.X - begin.X
//		k := float64(end.Y-begin.Y) / float64(deltaX)
//		points := make([]image.Point, (deltaX+1)*width)
//		for x := 0; x <= deltaX; x++ {
//			for y := 0; y < width; y++ {
//				points[x*width+y] = image.Point{begin.X + x, begin.Y + wMin + y + int(math.Round(k*float64(x)))}
//			}
//		}
//
//		return points
//	} else {
//		if begin.Y > end.Y {
//			begin, end = end, begin
//		}
//		deltaY := end.Y - begin.Y
//		k := float64(end.X-begin.X) / float64(deltaY)
//		points := make([]image.Point, (deltaY+1)*width)
//		for y := 0; y <= deltaY; y++ {
//			for x := 0; x < width; x++ {
//				points[y*width+x] = image.Point{begin.X + wMin + x + int(math.Round(k*float64(y))), begin.Y + y}
//			}
//		}
//
//		return points
//	}
//}
//
//func AddHLine(img draw.Image, x1, y, x2 int, clr color.Color) {
//	for ; x1 <= x2; x1++ {
//		img.Set(x1, y, clr)
//	}
//}
//
//func AddVLine(img draw.Image, x, y1, y2 int, clr color.Color) {
//	for ; y1 <= y2; y1++ {
//		img.Set(x, y1, clr)
//	}
//}
//
//func AddLine(img draw.Image, ls plane.Segment, clr color.Color) {
//	x1Int, y1Int, x2Int, y2Int := int(math.Round(ls[0].X)), int(math.Round(ls[0].Y)), int(math.Round(ls[1].X)), int(math.Round(ls[1].Y))
//
//	k := 0.
//	if x1Int == x2Int {
//		img.Set(x2Int, y2Int, clr)
//	} else {
//		if x2Int < x1Int {
//			x1Int, y1Int, x2Int, y2Int = x2Int, y2Int, x1Int, y1Int
//		}
//
//		k = float64(y2Int-y1Int) / float64(x2Int-x1Int)
//	}
//
//	for x := x1Int; x <= x2Int; x++ {
//		img.Set(x, y1Int+int(math.Round(float64(x-x1Int)*k)), clr)
//	}
//}
//
//func AddRectangle(img draw.Image, rect image.Rectangle, clr color.Color) {
//	AddHLine(img, rect.Min.X, rect.Min.Y, rect.Max.X, clr)
//	AddHLine(img, rect.Min.X, rect.Max.Y, rect.Max.X, clr)
//	AddVLine(img, rect.Min.X, rect.Min.Y, rect.Max.Y, clr)
//	AddVLine(img, rect.Max.X, rect.Min.Y, rect.Max.Y, clr)
//
//}
