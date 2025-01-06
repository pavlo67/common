package convolution

//import (
//	"fmt"
//	"image"
//
//	"github.com/pavlo67/common/common"
//	"github.com/pavlo67/common/common/imagelib"
//	"github.com/pavlo67/maps/mapping/imager/layer"
//)
//
//var _ convolution.Mask = &classesMask{}
//
//type classesMask struct {
//	imgRGB    *image.RGBA
//	rangeLow  uint8
//	rangeHigh uint8
//}
//
//const onClasses = "on Classes()"
//
//func Classes(deltaMax uint8) (Mask, error) {
//	return &classesMask{
//		rangeLow:  deltaMax,
//		rangeHigh: 0xFF - deltaMax,
//	}, nil
//}
//
//func (mask *classesMask) Side() int {
//	return 1
//}
//
//const onClassesPrepare = "on Classes.GetNext()"
//
//func (mask *classesMask) GetNext(onData interface{}) error {
//	switch v := onData.(type) {
//	case image.RGBA:
//		mask.imgRGB = &v
//	case *image.RGBA:
//		mask.imgRGB = v
//case imager.LayerRGBA:
//mask.imgRGB = &v.RGBA
//case *imager.LayerRGBA:
//mask.imgRGB = &v.RGBA
//case imager.Frame:
//mask.imgRGB = &v.RGBA
//case *imager.Frame:
//mask.imgRGB = &v.RGBA
//	}
//	if mask.imgRGB == nil {
//		if onData == nil {
//			return fmt.Errorf("onData == nil (%#v) / "+onClassesPrepare, onData)
//		}
//		return fmt.Errorf("wrong data (%T) / "+onClassesPrepare, onData)
//	}
//
//	return nil
//}
//
//func (mask classesMask) Info() common.Mappings {
//	return common.Mappings{
//		"name":     "classes",
//		"rangeLow": mask.rangeLow,
//	}
//}
//
//func (mask *classesMask) InitStat(thrClose float64) error {
//	return nil
//}
//
//func (mask classesMask) Classes() interface{} {
//	return nil
//}
//
//func (mask classesMask) Calculate(x, y int) layer.Value {
//	offset := (y-mask.imgRGB.Rect.Min.Y)*mask.imgRGB.Stride + (x-mask.imgRGB.Rect.Min.X)*coloring.NumColorsRGBA
//	clr := mask.imgRGB.Pix[offset : offset+3]
//
//	r, g, b := clr[0], clr[1], clr[2]
//
//	if r >= mask.rangeHigh && r >= mask.rangeHigh && r >= mask.rangeHigh {
//		return 0xE0
//	} else if r <= mask.rangeLow && g <= mask.rangeLow && b <= mask.rangeLow {
//		return 0x00
//	} else if r >= g {
//		if b >= r {
//			return 0x20
//		} else if b >= g {
//			return 0x40
//		} else {
//			return 0x60
//		}
//	} else {
//		if b >= g {
//			return 0x80
//		} else if b >= r {
//			return 0xA0
//		} else {
//			return 0xC0
//		}
//	}
//
//}
