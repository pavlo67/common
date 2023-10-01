package imagelib

type ColorSelector struct {
	NumColors   int
	Selected    int
	SelectedMax int
}

var ColorSelectorFullRGB = ColorSelector{
	NumColors:   NumColorsRGBA,
	Selected:    0,
	SelectedMax: NumColorsRGB,
}
