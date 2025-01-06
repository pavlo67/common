package _off

import (
	"fmt"
	"image"
	"sort"
	"strconv"
)

const onRGBAToTab = "on RGBAToTab()"

func RGBAToTab(imgRGB image.RGBA) (red, green, blue string) {
	rect := imgRGB.Bounds()

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		offset := (y - rect.Min.Y) * imgRGB.Stride

		for x := rect.Min.X; x < rect.Max.X; x++ {
			clr := imgRGB.Pix[offset : offset+3]
			red += strconv.Itoa(int(clr[0])) + "\t"
			green += strconv.Itoa(int(clr[1])) + "\t"
			blue += strconv.Itoa(int(clr[2])) + "\t"
			offset += coloring.NumColorsRGBA
		}

		red += "\n"
		green += "\n"
		blue += "\n"
	}

	return red, green, blue
}

type ColorStat struct {
	ColorIndex int32
	Count      int64
}

func (cs ColorStat) String() string {
	r, g, b := uint8((cs.ColorIndex&0x00ff0000)>>16), uint8((cs.ColorIndex&0x0000ff00)>>8), uint8(cs.ColorIndex&0x000000ff)
	return fmt.Sprintf("%d / %d / %d: %d", r, g, b, cs.Count)
}

const onRGBAStat = "on RGBAStat()"

func RGBAStat(imgRGB image.RGBA, colorDelta int32) []ColorStat {
	rect := imgRGB.Bounds()

	colorCounts := map[int32]int64{}

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		offset := (y - rect.Min.Y) * imgRGB.Stride

		for x := rect.Min.X; x < rect.Max.X; x++ {
			clr := imgRGB.Pix[offset : offset+3]
			colorCounts[((int32(clr[0])/colorDelta)<<16+(int32(clr[1])/colorDelta)<<8+(int32(clr[2])/colorDelta))]++
			offset += coloring.NumColorsRGBA
		}
	}

	var colorStats []ColorStat
	for colorKey, cnt := range colorCounts {
		colorStats = append(colorStats, ColorStat{colorKey, cnt})
	}
	sort.Slice(colorStats, func(i, j int) bool { return colorStats[i].Count > colorStats[j].Count })

	return colorStats
}
