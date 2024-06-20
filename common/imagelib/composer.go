package imagelib

import (
	"fmt"
	"image"
)

// TODO!!! be careful: imgs [][]image.Image is a list of image rows so it's related to the next indexing: imgs[y][x]

const onComposeImages = "on ComposeImages()"

func ComposeImages(imgs [][]image.Image) (*image.RGBA, error) {
	if len(imgs) < 1 || len(imgs[0]) < 1 {
		return nil, nil
	} else if imgs[0][0] == nil {
		return nil, fmt.Errorf("imgs[0][0] == nil / " + onComposeImages)
	}

	rect := imgs[0][0].Bounds()

	imgComposed := image.NewRGBA(image.Rect(0, 0, rect.Dx()*len(imgs[0]), rect.Dy()*len(imgs)))

	for y, imgsY := range imgs {
		if y > 0 && len(imgsY) != len(imgs[0]) {
			return nil, fmt.Errorf("len(imgs[%d]) != len(imgs[0]): %d vs %d / "+onComposeImages, y, len(imgsY), len(imgs[0]))
		}
		for x, imgXY := range imgsY {
			if imgXY == nil {
				return nil, fmt.Errorf("imgs[%d][%d] == nil / "+onComposeImages, y, x)
			} else if imgXY.Bounds() != rect {
				return nil, fmt.Errorf("imgs[%d][%d].Bounds() != imgs[0][0].Bounds(): %v vs %v / "+onComposeImages, y, x, imgXY.Bounds(), rect)
			}
			for y0 := rect.Min.Y; y0 < rect.Max.Y; y0++ {
				for x0 := rect.Min.X; x0 < rect.Max.X; x0++ {
					imgComposed.Set(
						((x-rect.Min.X)*rect.Dx() + x0 - rect.Min.X), ((y-rect.Min.Y)*rect.Dy() + y0 - rect.Min.Y),
						imgXY.At(x0, y0))
				}
			}
		}
	}

	return imgComposed, nil
}
