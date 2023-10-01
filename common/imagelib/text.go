package imagelib

import (
	"image"
	"image/color"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
)

const DPIDefault = 72.
const SizeDefault = 18.
const SpacingDefault = 1.5
const FontfileDefault = "_fonts/LiberationMono-Regular.ttf"

var f *truetype.Font

func Write(drawImage draw.Image, point image.Point, dpi, size, spacing float64, fontFile string, clr color.Color, text []string) (int32, error) {

	if f == nil || fontFile != "" {
		if fontFile == "" {
			fontFile = FontfileDefault
		}

		fontBytes, err := os.ReadFile(fontFile)
		if err != nil {
			return 0, err
		}

		f, err = freetype.ParseFont(fontBytes)
		if err != nil {
			return 0, err
		}
	}

	if dpi <= 0 {
		dpi = DPIDefault
	}
	if size <= 0 {
		size = SizeDefault
	}
	if spacing <= 0 {
		spacing = SpacingDefault
	}

	ctx := freetype.NewContext()
	ctx.SetDPI(dpi)
	ctx.SetFont(f)
	ctx.SetFontSize(size)
	ctx.SetClip(drawImage.Bounds())
	ctx.SetDst(drawImage)
	ctx.SetSrc(image.NewUniform(clr))
	ctx.SetHinting(font.HintingFull)

	// (10, 10) for start from left top corner with some margin
	pt := freetype.Pt(point.X, point.Y+int(ctx.PointToFixed(size)>>6))

	for _, t := range text {
		_, err := ctx.DrawString(t, pt)
		if err != nil {
			return int32(pt.Y), err
		}
		pt.Y += ctx.PointToFixed(size * spacing)

	}

	return int32(pt.Y), nil
}
