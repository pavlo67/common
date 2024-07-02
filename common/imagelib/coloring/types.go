package coloring

type Code int8

const Gray Code = 0
const RGBA Code = 1

type ColorSelector struct {
	NumColors int
	Selected  int
}

const NumColorsRGB = 3
const NumColorsRGBA = 4

var ColorSelectorFullRGB = ColorSelector{
	NumColors: NumColorsRGBA,
	Selected:  0,
}
