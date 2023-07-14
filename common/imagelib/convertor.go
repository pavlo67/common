package imagelib

import (
	"fmt"
	"image"
)

const PixMax = 0xFF
const NumColorsRGB = 3
const NumColorsRGBA = 4

const ChRed = 0
const ChGreen = 1
const ChBlue = 2

const onRGBToGray = "on RGBToGray()"

func RGBToGray(rgba image.RGBA, colorNum int) (*image.Gray, error) {
	if colorNum >= NumColorsRGB {
		return nil, fmt.Errorf(onRGBToGray+": wrong color to get: %d", colorNum)
	}

	xWidth, yHeight := rgba.Rect.Max.X-rgba.Rect.Min.X, rgba.Rect.Max.Y-rgba.Rect.Min.Y
	if xWidth <= 0 || yHeight <= 0 {
		return nil, fmt.Errorf(onRGBToGray+": empty img.Rect (%#v)", rgba.Rect)
	}

	gray := image.Gray{
		Pix:    make([]uint8, xWidth*yHeight),
		Stride: xWidth,
		Rect:   image.Rectangle{Max: image.Point{X: xWidth, Y: yHeight}},
	}

	if colorNum >= 0 {
		for y := 0; y < yHeight; y++ {
			rgbaStride := y*rgba.Stride + colorNum
			grayStride := y * xWidth
			for x := 0; x < xWidth; x++ {
				gray.Pix[grayStride+x] = rgba.Pix[rgbaStride+x*NumColorsRGBA]
			}
		}
	} else {
		for y := 0; y < yHeight; y++ {
			for x := 0; x < xWidth; x++ {
				gray.Set(x, y, rgba.At(x, y))
			}
		}
	}

	return &gray, nil
}

func GrayValueAvg(img *image.RGBA, x, y float64, csSelected int) float64 {
	dx, dy := x-float64(int(x)), y-float64(int(y))
	if dx > 0 {
		if dy > 0 {
			return float64(img.Pix[int(y)*img.Stride+int(x)*NumColorsRGBA+csSelected])*(1-dx)*(1-dy) +
				float64(img.Pix[(int(y)+1)*img.Stride+int(x)*NumColorsRGBA+csSelected])*(1-dx)*dy +
				float64(img.Pix[(int(y)+1)*img.Stride+(int(x)+1)*NumColorsRGBA+csSelected])*dx*dy +
				float64(img.Pix[int(y)*img.Stride+(int(x)+1)*NumColorsRGBA+csSelected])*dx*(1-dy)
		}
		return float64(img.Pix[int(y)*img.Stride+int(x)*NumColorsRGBA+csSelected])*(1-dx) +
			float64(img.Pix[int(y)*img.Stride+(int(x)+1)*NumColorsRGBA+csSelected])*dx
	} else {
		if dy > 0 {
			return float64(img.Pix[int(y)*img.Stride+int(x)*NumColorsRGBA+csSelected])*(1-dy) +
				float64(img.Pix[(int(y)+1)*img.Stride+int(x)*NumColorsRGBA+csSelected])*dy
		}
		return float64(img.Pix[int(y)*img.Stride+int(x)*NumColorsRGBA+csSelected])
	}
}

func GrayValue(img *image.RGBA, x, y, csSelected int) float64 {
	return float64(img.Pix[y*img.Stride+x*NumColorsRGBA+csSelected])
}
