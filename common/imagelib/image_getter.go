package imagelib

import (
	"fmt"
	"image"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/logger"

	"github.com/pavlo67/common/common/imagelib/coloring"
)

type SubImage interface {
	SubImage(image.Rectangle) image.Image
}

type Bounded interface {
	Bounds() image.Rectangle
}

type Described interface {
	Bounds() image.Rectangle
	Description() Settings
}

// Imager -------------------------------------------------------------------------

var _ logger.GetImage = &GetImage{}

func Get(img image.Image, getMasks ...GetMask) logger.GetImage {
	return &GetImage{
		Images:     []image.Image{img},
		ImageMasks: getMasks,
	}
}

type GetImage struct {
	Images     []image.Image
	ImageMasks []GetMask
}

func (op *GetImage) Bounds() image.Rectangle {
	if op == nil || len(op.Images) < 1 {
		return image.Rectangle{}
	}

	return op.Images[0].Bounds()
}

const onImage = "on Imager.Image()"

func (op *GetImage) Image(opts common.Map) (image.Image, string, error) {
	if op == nil || len(op.Images) < 1 {
		return nil, "", fmt.Errorf(onImage + ": op == nil || len(op.Images) == 0")
	}

	img := ImageToRGBACopied(op.Images[0])
	for _, imgToAdd := range op.Images[1:] {
		rect := imgToAdd.Bounds()
		for x := rect.Min.X; x < rect.Max.X; x++ {
			for y := rect.Min.Y; y < rect.Max.Y; y++ {
				img.Set(x, y, imgToAdd.At(x, y))
			}
		}
	}

	var info string
	var mask MasksOneColor
	for i, maskI := range op.ImageMasks {
		colorNamed := maskI.Color()
		if colorNamed == nil || colorNamed.Color == nil {
			colorNamed = &coloring.RoundAbout[i%len(coloring.RoundAbout)]
		}
		mask = append(mask, maskI.Mask(colorNamed.Color, opts)...)
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
