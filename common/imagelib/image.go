package imagelib

import (
	"fmt"
	"image"
)

type Imager interface {
	Image() (image.Image, string, error)
	Bounds() image.Rectangle
}

type SubImage interface {
	SubImage(image.Rectangle) image.Image
}

// Imager -------------------------------------------------------------------------

var _ Imager = &GetImage{}

type GetImage struct {
	Images     []image.Image
	ImageMasks []GetMask
}

func (imageOp *GetImage) Bounds() image.Rectangle {
	if imageOp == nil || len(imageOp.Images) < 1 {
		return image.Rectangle{}
	}

	return imageOp.Images[0].Bounds()
}

const onImage = "on Imager.Imager()"

func (imageOp *GetImage) Image() (image.Image, string, error) {
	if imageOp == nil || len(imageOp.Images) < 1 {
		return nil, "", fmt.Errorf(onImage + ": imageOp == nil || len(imageOp.Images) == 0")
	}

	img := ImageToRGBACopied(imageOp.Images[0])
	for _, imgToAdd := range imageOp.Images[1:] {
		rect := imgToAdd.Bounds()
		for x := rect.Min.X; x < rect.Max.X; x++ {
			for y := rect.Min.Y; y < rect.Max.Y; y++ {
				img.Set(x, y, imgToAdd.At(x, y))
			}
		}
	}

	var info string
	var mask Mask
	for i, maskI := range imageOp.ImageMasks {
		colorNamed := maskI.Color()
		if colorNamed == nil || colorNamed.Color == nil {
			colorNamed = &RoundAbout[i%len(RoundAbout)]
		}
		mask = append(mask, maskI.Mask(colorNamed.Color)...)
		info += maskI.Info(*colorNamed)
	}

	mask.ShowOn(img)

	return img, info, nil
}

//// GetImageGray -------------------------------------------------------------------------
//
//var _ Imager = &GetImageGray{}
//
//type GetImageGray struct {
//	Gray *image.Gray
//}
//
//func (imageOp GetImageGray) Image() (image.Image, string, error) {
//	return imageOp.Gray, "", nil
//}
//
//func (imageOp *GetImageGray) Bounds() image.Rectangle {
//	if imageOp == nil || imageOp.Gray == nil {
//		return image.Rectangle{}
//	}
//
//	return imageOp.Gray.Rect
//}
