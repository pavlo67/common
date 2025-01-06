package layers

import (
	"github.com/pavlo67/common/common/imagelib"
	"github.com/pavlo67/common/common/imagelib/pix"
	"image"
)

var _ imagelib.Described = &Layer{}

type Layer struct {
	image.Gray
	imagelib.Settings
	Metrics
}

type Layers map[string]*Layer

func (lyr Layer) Description() imagelib.Settings {
	return lyr.Settings
}

func (lyr *Layer) Offset(x, y int) int {
	return (y-lyr.Rect.Min.Y)*lyr.Stride + (x - lyr.Rect.Min.X)
}

func (lyr Layer) Length() int64 {
	return int64(lyr.Rect.Max.Y-lyr.Rect.Min.Y) * int64(lyr.Rect.Max.X-lyr.Rect.Min.X)
}

type Metrics struct {
	Min, Max      pix.Value
	ClassesCustom ClassesCustom
	Classes256    Classes256
	Criterion     float64 // TODO!!! keep it for shifting
}

type Classes256 [256]int32

type ClassesCustom []int32

func (classesCustom ClassesCustom) Range() pix.Value {
	if len(classesCustom) <= 0 {
		return 0

	}

	r := pix.Value(256 / len(classesCustom))
	if pix.Value(256%len(classesCustom)) >= r {
		return 0
	}

	return r
}
