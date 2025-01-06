package convolution

import (
	"fmt"
	"github.com/pavlo67/common/common/imagelib/coloring"
	"github.com/pavlo67/common/common/imagelib/pix"
	"image"
	"strconv"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/imagelib/frame"
	"github.com/pavlo67/imagelib/layers"
)

var _ Mask = &RGBBrightnessClassesMask{}

type RGBBrightnessClassesMask struct {
	imgRGB     *image.RGBA
	classRange pix.Value
	classes    layers.ClassesCustom
}

const onRGBBrightnessClasses = "on RGBBrightnessClasses()"

func RGBBrightnessClasses(classRange pix.Value) (Mask, error) {
	if classRange == 0 {
		return nil, fmt.Errorf("classRange == 0 / " + onRGBBrightnessClasses)
	}

	lenTotal := int(pix.ValueMax) + 1
	classesNum := lenTotal / int(classRange)
	if lenTotal%int(classRange) > 0 {
		classesNum++
	}

	return &RGBBrightnessClassesMask{
		classRange: classRange,
		classes:    make([]int32, classesNum),
	}, nil
}

func (mask *RGBBrightnessClassesMask) Side() int {
	return 1
}

const onRGBBrightnessPrepare = "on RGBBrightnessClasses.GetNext()"

func (mask *RGBBrightnessClassesMask) Prepare(onData interface{}) error {
	switch v := onData.(type) {
	case image.RGBA:
		mask.imgRGB = &v
	case *image.RGBA:
		mask.imgRGB = v
	case frame.LayerRGBA:
		mask.imgRGB = &v.RGBA
	case *frame.LayerRGBA:
		mask.imgRGB = &v.RGBA
	case frame.Frame:
		mask.imgRGB = &v.RGBA
	case *frame.Frame:
		mask.imgRGB = &v.RGBA
	}
	if mask.imgRGB == nil {
		if onData == nil {
			return fmt.Errorf("onData == nil (%#v) / "+onRGBBrightnessPrepare, onData)
		}
		return fmt.Errorf("wrong data (%T) / "+onRGBBrightnessPrepare, onData)
	}

	mask.classes = make([]int32, len(mask.classes))

	return nil
}

func (mask RGBBrightnessClassesMask) Info() common.Map {
	return common.Map{
		"name":       "br_classes_" + strconv.Itoa(int(mask.classRange)),
		"classRange": mask.classRange,
	}
}

func (mask RGBBrightnessClassesMask) Classes() layers.ClassesCustom {
	return mask.classes
}

func (mask *RGBBrightnessClassesMask) Calculate(x, y int) pix.Value {
	offset := (y-mask.imgRGB.Rect.Min.Y)*mask.imgRGB.Stride + (x-mask.imgRGB.Rect.Min.X)*coloring.NumColorsRGBA

	brightness := pix.ValueSum(mask.imgRGB.Pix[offset]) + pix.ValueSum(mask.imgRGB.Pix[offset+1]) + pix.ValueSum(mask.imgRGB.Pix[offset+2])

	classNum := pix.Value(brightness / (3 * pix.ValueSum(mask.classRange)))
	mask.classes[classNum]++

	return mask.classRange * classNum
}
