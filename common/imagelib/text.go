package imagelib

import (
	"image"
	"os"

	"github.com/golang/freetype/truetype"

	"github.com/golang/freetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
)

const dpiDefault = 72.
const fontfileDefault = "_fonts/LiberationMono-Regular.ttf"

var f *truetype.Font

func Write(drawImage draw.Image, dpi, size, spacing float64, fontfile string, imgClr image.Image, text []string) (int32, error) {

	if f == nil {
		if fontfile == "" {
			fontfile = fontfileDefault
		}

		fontBytes, err := os.ReadFile(fontfile)
		if err != nil {
			return 0, err
		}

		f, err = freetype.ParseFont(fontBytes)
		if err != nil {
			return 0, err
		}
	}

	if dpi <= 0 {
		dpi = dpiDefault
	}

	ctx := freetype.NewContext()
	ctx.SetDPI(dpi)
	ctx.SetFont(f)
	ctx.SetFontSize(size)
	ctx.SetClip(drawImage.Bounds())
	ctx.SetDst(drawImage)
	ctx.SetSrc(imgClr)
	ctx.SetHinting(font.HintingFull)

	pt := freetype.Pt(10, 10+int(ctx.PointToFixed(size)>>6))

	for _, t := range text {
		_, err := ctx.DrawString(t, pt)
		if err != nil {
			return int32(pt.Y), err
		}
		pt.Y += ctx.PointToFixed(size * spacing)

	}

	return int32(pt.Y), nil
}
