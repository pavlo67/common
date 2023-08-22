package imagelib

import (
	"fmt"
	"image"
)

const onComposeImages = "on imagelib.ComposeImages()"

func ComposeImages(imgs [][]image.Image) (*image.RGBA, error) {
	if len(imgs) < 1 || len(imgs[0]) < 1 {
		return nil, nil
	} else if imgs[0][0] == nil {
		return nil, fmt.Errorf("imgs[0][0] == nil / " + onComposeImages)
	}

	rect := imgs[0][0].Bounds()

	imgComposed := image.NewRGBA(image.Rect(0, 0, rect.Dx()*len(imgs), rect.Dy()*len(imgs[0])))

	for x, imgsX := range imgs {
		if x > 0 && len(imgsX) != len(imgs[0]) {
			return nil, fmt.Errorf("len(imgs[%d]) != len(imgs[0]): %d vs %d / "+onComposeImages, x, len(imgsX), len(imgs[0]))
		}
		for y, imgXY := range imgsX {
			if imgXY == nil {
				return nil, fmt.Errorf("imgs[%d][%d] == nil / "+onComposeImages, x, y)
			} else if imgXY.Bounds() != rect {
				return nil, fmt.Errorf("imgs[%d][%d].Bounds() != imgs[0][0].Bounds(): %v vs %v / "+onComposeImages, x, y, imgXY.Bounds(), rect)
			}
			for y0 := rect.Min.Y; y0 < rect.Max.Y; y0++ {
				for x0 := rect.Min.X; x0 < rect.Max.X; x0++ {
					imgComposed.Set(
						((x-rect.Min.X)*rect.Dx() + x0 - rect.Min.X),
						((y-rect.Min.Y)*rect.Dy() + y0 - rect.Min.X),
						imgXY.At(x0, y0))
				}
			}
		}
	}

	return imgComposed, nil
}
