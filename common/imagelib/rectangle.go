package imagelib

import (
	"image"
	"log"

	"github.com/pavlo67/common/common/geolib"
)

func Normalize(rect image.Rectangle) image.Rectangle {
	rect = rect.Canon()

	return image.Rectangle{Max: image.Point{rect.Max.X - rect.Min.X, rect.Max.Y - rect.Min.Y}}
}

func SubArea(img image.Image, area, subArea geolib.Area) image.Image {

	// TODO!!!
	// var center image.Point2
	// rectInner := image.Rectangle{image.Point2{center.XT - cfg.XWidth/2, center.YT - cfg.YHeight/2}, image.Point2{center.XT + cfg.XWidth/2, center.YT + cfg.YHeight/2}}

	log.Fatal("on imagelib/rectangle.SubArea()")

	return nil
}
