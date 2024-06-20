package imagelib

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

const NumColorsRGB = 3
const NumColorsRGBA = 4

const ChRed = 0
const ChGreen = 1
const ChBlue = 2

const onImageToGray = "on ImageToGray()"

func ImageToGray(img image.Image) (*image.Gray, error) {
	if img == nil {
		return nil, fmt.Errorf(onImageToGray + ": nil img")
	}

	if grayPtr, ok := img.(*image.Gray); ok {
		if grayPtr != nil {
			return grayPtr, nil
		}
		return nil, fmt.Errorf(onImageToGray + ": nil img.(*Gray)")
	}

	rect := img.Bounds()
	gray := image.NewGray(rect)

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			gray.Set(x, y, img.At(x, y))
		}
	}

	return gray, nil
}

func ImageToGrayCopied(img image.Image) *image.Gray {
	if img == nil {
		return nil
	}

	rect := img.Bounds()
	gray := image.NewGray(rect)

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			gray.Set(x, y, img.At(x, y))
		}
	}

	return gray
}

const onImageToRGBA = "on ImageToRGBA()"

func ImageToRGBA(img image.Image) (*image.RGBA, error) {
	if img == nil {
		return nil, fmt.Errorf(onImageToRGBA + ": nil img")
	}

	if rgbaPtr, ok := img.(*image.RGBA); ok {
		if rgbaPtr != nil {
			return rgbaPtr, nil
		}
		return nil, fmt.Errorf(onImageToRGBA + ": nil img.(*RGBA)")
	}

	rect := img.Bounds()
	rgba := image.NewRGBA(rect)

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}

	return rgba, nil
}

func ImageToRGBACopied(img image.Image) *image.RGBA {
	if img == nil {
		return nil
	}

	rect := img.Bounds()
	rgba := image.NewRGBA(rect)

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}

	return rgba
}

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
			rgbaStride := (y-rgba.Rect.Min.Y)*rgba.Stride + colorNum
			grayStride := y * xWidth
			for x := 0; x < xWidth; x++ {
				gray.Pix[grayStride+x] = rgba.Pix[rgbaStride+(x-rgba.Rect.Min.X)*NumColorsRGBA]
			}
		}
	} else {
		for y := 0; y < yHeight; y++ {
			for x := 0; x < xWidth; x++ {
				gray.Set(x+rgba.Rect.Min.X, y+rgba.Rect.Min.Y, rgba.At(x, y))
			}
		}
	}

	return &gray, nil
}
