package imagelib

import (
	"image"
)

func Normalize(rect image.Rectangle) image.Rectangle {
	rect = rect.Canon()

	return image.Rectangle{Max: image.Point{rect.Max.X - rect.Min.X, rect.Max.Y - rect.Min.Y}}
}
