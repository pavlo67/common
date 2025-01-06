package convolution

import (
	"fmt"
	"github.com/pavlo67/common/common/imagelib/pix"
	"image"
	"strconv"

	"github.com/pavlo67/common/common"

	"github.com/pavlo67/imagelib/layers"
)

var _ Mask = &GrayBrightnessClassesMask{}

type GrayBrightnessClassesMask struct {
	imgGray    *image.Gray
	classRange pix.Value
	classes    layers.ClassesCustom
}

const onGrayBrightnessClasses = "on GrayBrightnessClasses()"

func GrayBrightnessClasses(classRange pix.Value) (Mask, error) {
	if classRange == 0 {
		return nil, fmt.Errorf("classRange == 0 / " + onGrayBrightnessClasses)
	}
	classesNum := (int(pix.ValueMax) + 1) / int(classRange)
	if pix.ValueMax%classRange > 0 {
		classesNum++
	}

	return &GrayBrightnessClassesMask{
		classRange: classRange,
		classes:    make([]int32, classesNum),
	}, nil
}

func (mask *GrayBrightnessClassesMask) Side() int {
	return 1
}

const onGrayBrightnessPrepare = "on GrayBrightnessClasses.GetNext()"

func (mask *GrayBrightnessClassesMask) Prepare(onData interface{}) error {
	switch v := onData.(type) {
	case image.Gray:
		mask.imgGray = &v
	case *image.Gray:
		mask.imgGray = v
	case layers.Layer:
		mask.imgGray = &v.Gray
	case *layers.Layer:
		mask.imgGray = &v.Gray
	}
	if mask.imgGray == nil {
		if onData == nil {
			return fmt.Errorf("onData == nil (%#v) / "+onGrayBrightnessPrepare, onData)
		}
		return fmt.Errorf("wrong data (%T) / "+onGrayBrightnessPrepare, onData)
	}

	mask.classes = make([]int32, len(mask.classes))

	return nil
}

func (mask GrayBrightnessClassesMask) Info() common.Map {
	return common.Map{
		"name":       "br_classes_" + strconv.Itoa(int(mask.classRange)),
		"classRange": mask.classRange,
	}
}

func (mask GrayBrightnessClassesMask) Classes() layers.ClassesCustom {
	return mask.classes
}

func (mask *GrayBrightnessClassesMask) Calculate(x, y int) pix.Value {
	offset := (y-mask.imgGray.Rect.Min.Y)*mask.imgGray.Stride + (x - mask.imgGray.Rect.Min.X)

	classNum := mask.imgGray.Pix[offset] / mask.classRange
	mask.classes[classNum]++

	return mask.classRange * classNum
}
