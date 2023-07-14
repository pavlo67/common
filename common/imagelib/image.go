package imagelib

import (
	"fmt"
	"image"
)

type GetImage interface {
	Image() (image.Image, string, error)
	Bounds() image.Rectangle
}

type SubImage interface {
	SubImage(image.Rectangle) image.Image
}

//// GetImage -----------------------------------------------------------------------------
//
//var _ GetImage = &GetImage{}
//
//type GetImage struct {
//	Img  image.GetImage
//	Rect image.Rectangle
//}
//
//func (imageOp *GetImage) Bounds() image.Rectangle {
//	if imageOp == nil || imageOp.Img == nil {
//		return image.Rectangle{}
//	}
//
//	return imageOp.Rect
//}
//
//func (imageOp GetImage) GetImage() (image.GetImage, string, error) {
//	return imageOp.Img, "", nil
//}

// GetImageRGBA -------------------------------------------------------------------------

var _ GetImage = &GetImageRGBA{}

type GetImageRGBA struct {
	RGBA        *image.RGBA
	ResizeRatio float64
	ImageMasks  []GetMask
}

func (imageOp *GetImageRGBA) Bounds() image.Rectangle {
	if imageOp == nil || imageOp.RGBA == nil {
		return image.Rectangle{}
	}

	return imageOp.RGBA.Rect
}

const onImage = "on GetImageRGBA.GetImage()"

func (imageOp *GetImageRGBA) Image() (image.Image, string, error) {
	if imageOp == nil || imageOp.RGBA == nil {
		return nil, "", fmt.Errorf(onImage + ": imageOp == nil || imageOp.RGBA == nil")
	}

	img, _, err := Resize(*imageOp.RGBA, imageOp.ResizeRatio)
	if err != nil {
		return nil, "", err
	} else if img == nil {
		return nil, "", fmt.Errorf("resized img == nil / " + onImage)
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

// GetImageGray -------------------------------------------------------------------------

var _ GetImage = &GetImageGray{}

type GetImageGray struct {
	Gray *image.Gray
}

func (imageOp GetImageGray) Image() (image.Image, string, error) {
	return imageOp.Gray, "", nil
}

func (imageOp *GetImageGray) Bounds() image.Rectangle {
	if imageOp == nil || imageOp.Gray == nil {
		return image.Rectangle{}
	}

	return imageOp.Gray.Rect
}
