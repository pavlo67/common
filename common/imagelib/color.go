package imagelib

import (
	"image/color"

	"golang.org/x/image/colornames"
)

func Brightness(clr color.RGBA) uint32 {
	return uint32(clr.R) * uint32(clr.G) * uint32(clr.B)
}

type ColorNamed struct {
	color.Color
	Name string
}

var RoundAbout = []ColorNamed{
	{colornames.Red, "red"},
	{colornames.Orange, "orange"},
	{colornames.Yellow, "yellow"},
	{colornames.Green, "green"},
	{colornames.Blue, "blue"},
	//{colornames.Pink, "pink"},
	//{colornames.Violet, "violet"},
}
