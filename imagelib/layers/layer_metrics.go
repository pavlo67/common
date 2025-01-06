package layers

func (lyr *Layer) MinMax() {
	if lyr == nil {
		return
	}

	xWidth, yHeight := lyr.Rect.Max.X-lyr.Rect.Min.X, lyr.Rect.Max.Y-lyr.Rect.Min.Y

	minValue, maxValue := lyr.Pix[0], lyr.Pix[0]
	var offset int
	for y := 0; y < yHeight; y++ {
		for x := 0; x < xWidth; x++ {
			if v := lyr.Pix[offset+x]; v < minValue {
				minValue = v
			} else if v > maxValue {
				maxValue = v
			}
		}

		offset += lyr.Stride
	}
	lyr.Min, lyr.Max = minValue, maxValue
}

//func CountGrayMetrics(img image.GrayWide) Metrics {
//	xWidth := img.Rect.Max.X - img.Rect.Min.X
//	yHeight := img.Rect.Max.Y - img.Rect.Min.Y
//
//	grayMetrics := Metrics{}
//
//	var cnt, sum, sumSquared uint64
//
//	for y := 0; y < yHeight; y++ {
//		stride := y * img.Stride
//		for x := 0; x < xWidth; x++ {
//			val := img.Pix[stride+x]
//			sum += uint64(val)
//			sumSquared += uint64(val) * uint64(val)
//			cnt++
//		}
//	}
//
//	grayMetrics.Avg = Value(float64(sum) / float64(cnt))
//	grayMetrics.Criterion = Value(math.Sqrt(float64(sumSquared) / float64(cnt)))
//
//	return grayMetrics
//}
