package imagelib

// DEPRECATED
//func GrayAddHLine(gray image.Gray, x1, y, x2 int, clr color.Color) {
//	for ; x1 <= x2; x1++ {
//		gray.Set(x1, y, clr)
//	}
//}
//
//// DEPRECATED
//func GrayAddVLine(gray image.Gray, x, y1, y2 int, clr color.Color) {
//	for ; y1 <= y2; y1++ {
//		gray.Set(x, y1, clr)
//	}
//}
//
//// DEPRECATED
//func GrayAddRectangle(gray image.Gray, rect image.Rectangle, clr color.Color) {
//	GrayAddHLine(gray, rect.Min.X, rect.Min.Y, rect.Max.X, clr)
//	GrayAddHLine(gray, rect.Min.X, rect.Max.Y, rect.Max.X, clr)
//	GrayAddVLine(gray, rect.Min.X, rect.Min.Y, rect.Max.Y, clr)
//	GrayAddVLine(gray, rect.Max.X, rect.Min.Y, rect.Max.Y, clr)
//
//}
//
//// DEPRECATED
//func RGBAAddHLine(rgba image.RGBA, x1, y, x2 int, clr color.Color) {
//	for ; x1 <= x2; x1++ {
//		rgba.Set(x1, y, clr)
//	}
//}
//
//// DEPRECATED
//func RGBAAddVLine(rgba image.RGBA, x, y1, y2 int, clr color.Color) {
//	for ; y1 <= y2; y1++ {
//		rgba.Set(x, y1, clr)
//	}
//}
//
//// DEPRECATED
//func RGBAAddRectangle(rgba image.RGBA, rect image.Rectangle, clr color.Color) {
//	RGBAAddHLine(rgba, rect.Min.X, rect.Min.Y, rect.Max.X, clr)
//	RGBAAddHLine(rgba, rect.Min.X, rect.Max.Y, rect.Max.X, clr)
//	RGBAAddVLine(rgba, rect.Min.X, rect.Min.Y, rect.Max.Y, clr)
//	RGBAAddVLine(rgba, rect.Max.X, rect.Min.Y, rect.Max.Y, clr)
//
//}
//
//// DEPRECATED
//func GrayAddRuler(gray image.Gray, numOfPixels, numOfMeters uint, dpm float64, clr color.Color) {
//	for x := gray.Rect.Min.X; x < gray.Rect.Max.X; x += int(numOfPixels) {
//		GrayAddVLine(gray, x, gray.Rect.Min.Y, (gray.Rect.Min.Y+gray.Rect.Max.Y)/2, clr)
//	}
//
//	if dpm > 0 {
//		numOfMeterPixels := int(math.Round(float64(numOfMeters) * dpm))
//
//		for x := gray.Rect.Min.X; x < gray.Rect.Max.X; x += numOfMeterPixels {
//			GrayAddVLine(gray, x, (gray.Rect.Min.Y+gray.Rect.Max.Y)/2, gray.Rect.Max.Y, clr)
//		}
//	}
//}
