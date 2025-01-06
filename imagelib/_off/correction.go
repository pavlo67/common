package _off

import (
	"image"
	"image/color"
)

func CorrectChannels(imgRGB image.RGBA, dr, dg, db int8) image.RGBA {
	rect := imgRGB.Rect.Canon()

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			rgba, _ := imgRGB.At(x, y).(color.RGBA)

			//if rgba == nil {
			//	log.Fatal("rgba == nil / on CorrectChannels()")
			//}

			r1 := int(rgba.R) + int(dr)

			if r1 < 0 {
				r1 = 0
			} else if r1 > 0xFF {
				r1 = 0xFF
			}

			g1 := int(rgba.G) + int(dg)
			if g1 < 0 {
				g1 = 0
			} else if g1 > 0xFF {
				g1 = 0xFF
			}

			b1 := int(rgba.B) + int(db)
			if b1 < 0 {
				b1 = 0
			} else if b1 > 0xFF {
				b1 = 0xFF
			}

			imgRGB.Set(x, y, color.RGBA{uint8(r1), uint8(g1), uint8(b1), rgba.A})
		}
	}

	return imgRGB
}
