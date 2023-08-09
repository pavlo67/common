package imagelib

import (
	"image"
)

func RGBAveraged(imgRGB image.RGBA, canalHalfSide int, differenceMax uint8) image.RGBA {
	rect := imgRGB.Rect
	imgRGBNew := image.NewRGBA(imgRGB.Rect)

	var offset int

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		offsetX := offset
		for x := rect.Min.X; x < rect.Max.X; x++ {
			clrAvg := RGBAvg(&imgRGB, x, y, canalHalfSide, differenceMax)

			// fmt.Println(clrAvg)

			imgRGBNew.Pix[offsetX] = clrAvg[0]
			imgRGBNew.Pix[offsetX+1] = clrAvg[1]
			imgRGBNew.Pix[offsetX+2] = clrAvg[2]
			imgRGBNew.Pix[offsetX+3] = 255
			offsetX += NumColorsRGBA
		}

		offset += imgRGBNew.Stride
	}

	return *imgRGBNew
}

func RGBAvg(img *image.RGBA, x, y int, canalHalfSide int, differenceMax uint8) [NumColorsRGB]uint8 {

	if img == nil {
		return [NumColorsRGB]uint8{}
	}

	offset := (y-img.Rect.Min.Y)*img.Stride + (x-img.Rect.Min.X)*NumColorsRGBA
	clr := img.Pix[offset : offset+3]

	rectAvg := image.Rect(x-canalHalfSide, y-canalHalfSide, x+canalHalfSide+1, y+canalHalfSide+1)
	rect := img.Rect.Intersect(rectAvg)

	var cnt int
	var sum [NumColorsRGB]int

	xWidth, yHeight := rect.Max.X-rect.Min.X, rect.Max.Y-rect.Min.Y
	offsetY := (rect.Min.Y-img.Rect.Min.Y)*img.Stride + (rect.Min.X-img.Rect.Min.X-1)*NumColorsRGBA

	for y1 := 0; y1 < yHeight; y1++ {
		offsetX := offsetY

	X:
		for x1 := 0; x1 < xWidth; x1++ {
			offsetX += NumColorsRGBA

			for i := 0; i < NumColorsRGB; i++ {
				if img.Pix[offsetX+i] > clr[i] {
					if img.Pix[offsetX+i]-clr[i] > differenceMax {
						continue X
					}
				} else {
					if clr[i]-img.Pix[offsetX+i] > differenceMax {
						continue X
					}
				}
			}

			sum[0] += int(img.Pix[offsetX])
			sum[1] += int(img.Pix[offsetX+1])
			sum[2] += int(img.Pix[offsetX+2])
			cnt++
		}

		offsetY += img.Stride
	}

	if cnt <= 0 {
		return [NumColorsRGB]uint8{clr[0], clr[1], clr[2]}
	}

	return [NumColorsRGB]uint8{uint8(sum[0] / cnt), uint8(sum[1] / cnt), uint8(sum[2] / cnt)}
}

//func RGBAvg(img *image.RGBA, x, y int, canalHalfSide int) [NumColorsRGB]uint8 {
//
//	if img == nil {
//		return [NumColorsRGB]uint8{}
//	}
//
//	rectAvg := image.Rect(x-canalHalfSide, y-canalHalfSide, x+canalHalfSide+1, y+canalHalfSide+1)
//	rect := img.Rect.Intersect(rectAvg)
//
//	xWidth, yHeight := rect.Max.x-rect.Min.x, rect.Max.y-rect.Min.y
//	cnt := float64(xWidth * yHeight)
//
//	if cnt <= 0 {
//		offset := (y-img.Rect.Min.y)*img.Stride + (x-img.Rect.Min.x)*NumColorsRGBA
//		return [NumColorsRGB]uint8{img.Pix[offset],img.Pix[offset+1],img.Pix[offset+2]}
//	}
//
//	offsetY := (rect.Min.y-img.Rect.Min.y)*img.Stride + (rect.Min.x-img.Rect.Min.x)*NumColorsRGBA
//
//	var sum [NumColorsRGB]int
//
//	for y1 := 0; y1 < yHeight; y1++ {
//		offsetX := offsetY
//		for x1 := 0; x1 < xWidth; x1++ {
//			sum[0] += int(img.Pix[offsetX])
//			sum[1] += int(img.Pix[offsetX+1])
//			sum[2] += int(img.Pix[offsetX+2])
//			offsetX += NumColorsRGBA
//		}
//
//		offsetY += img.Stride
//	}
//
//	return [NumColorsRGB]uint8{
//		uint8(math.Round(float64(sum[0]) / cnt)),
//		uint8(math.Round(float64(sum[1]) / cnt)),
//		uint8(math.Round(float64(sum[2]) / cnt)),
//	}
//}
