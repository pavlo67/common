package imagelib

import (
	"image/color"

	"golang.org/x/image/colornames"
)

type PixDelta = int16
type PixSum = uint32

const PixMax = 0xFF

func Brightness(clr color.RGBA) uint32 {
	return uint32(clr.R) * uint32(clr.G) * uint32(clr.B)
}

type ColorNamed struct {
	color.Color
	Name string
}

var RoundAbout = []ColorNamed{
	{colornames.Black, "black"},
	{colornames.Red, "red"},
	{colornames.Aqua, "aqua"},
	{colornames.Green, "green"},
	{colornames.Blue, "blue"},

	// {colornames.White, "white"},
	// {colornames.Olive, "olive"},
	// {colornames.Orange, "orange"},
	// {colornames.Yellow, "yellow"},

}
