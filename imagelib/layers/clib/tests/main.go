package main

//#include "image_gray.h"
import "C"

import (
	"github.com/pavlo67/common/common/imagelib"
	"github.com/pavlo67/common/common/imagelib/pix"
	"github.com/pavlo67/imagelib/layers"
	"image"
	"log"
	"unsafe"

	"github.com/pavlo67/imagelib/layers/convolution"
)

func main() {

	brClassRange := pix.Value(4)

	imgGray := image.NewGray(image.Rect(0, 0, 1000, 1000))
	if imgGray == nil {
		log.Fatal("imgGray == nil")
	}
	imgGray.Pix[1000] = 255
	imgGray.Pix[1001] = 128
	imgGray.Pix[2001] = 28

	lyr := layers.Layer{
		Gray:     *imgGray,
		Settings: imagelib.Settings{DPM: 1},
	}

	maskBrClasses, err := convolution.RGBBrightnessClasses(brClassRange)
	if err != nil {
		log.Fatal(err)
	} else if maskBrClasses == nil {
		log.Fatal("maskBrClasses == nil")
	}

	lyrBrClasses, err := convolution.Layer(&lyr, maskBrClasses, 1, true)
	if err != nil {
		log.Fatal(err)
	} else if lyrBrClasses == nil {
		log.Fatal("lyrBrClasses == nil")
	}

	C.img_br_classes(toCImageGray(*imgGray), nil)
}

func toCImageGray(imgGray image.Gray) C.ImageGray {
	return C.ImageGray{
		x_width:  C.int32_t(imgGray.Rect.Max.X - imgGray.Rect.Min.X),
		y_height: C.int32_t(imgGray.Rect.Max.Y - imgGray.Rect.Min.Y),
		x_min:    C.int32_t(imgGray.Rect.Min.X),
		y_min:    C.int32_t(imgGray.Rect.Min.Y),
		stride:   C.int32_t(imgGray.Stride),
		pixels:   (*C.uchar)(unsafe.Pointer(&imgGray.Pix[0])),
	}
}
