package geolib

import (
	"fmt"
)

type XYRange [2]int

type XYRanges struct {
	Zoom   int
	XT, YT XYRange
}

func (xyRanges XYRanges) Key() string {
	return fmt.Sprintf("z%d_x%d_x%d_y%d_y%d", xyRanges.Zoom, xyRanges.XT[0], xyRanges.XT[1], xyRanges.YT[0], xyRanges.YT[1])
}

func (xyRanges XYRanges) Canon() XYRanges {
	xMin, xMax, yMin, yMax := xyRanges.XT[0], xyRanges.XT[1], xyRanges.YT[0], xyRanges.YT[1]
	if xMax < xMin {
		xMin, xMax = xMax, xMin
	}
	if yMax < yMin {
		yMin, yMax = yMax, yMin
	}

	return XYRanges{xyRanges.Zoom, XYRange{xMin, xMax}, XYRange{yMin, yMax}}
}

//func (xyRanges XYRanges) Area() Area {
//	xyRanges = xyRanges.Canon()
//
//	// TODO!!! use +1 in XYRanges max itself
//
//	return Area{
//		Tile{X: xyRanges.XT[0], Y: xyRanges.YT[0], Zoom: xyRanges.Zoom}.LeftTop(),
//		Tile{X: xyRanges.XT[1] + 1, Y: xyRanges.YT[1] + 1, Zoom: xyRanges.Zoom}.LeftTop(),
//	}
//}

func (xyRanges XYRanges) LeftTop() Point {

	xyRanges = xyRanges.Canon()

	return Tile{X: xyRanges.XT[0], Y: xyRanges.YT[0], Zoom: xyRanges.Zoom}.LeftTop()
}
