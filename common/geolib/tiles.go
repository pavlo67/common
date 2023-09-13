package geolib

import (
	"fmt"
	"math"
)

type Tile struct {
	X, Y int
	Zoom int
}

type XYRange [2]int

type XYRanges struct {
	Zoom   int
	XT, YT XYRange
}

func (xyRanges XYRanges) Key() string {
	return fmt.Sprintf("z%d_x%d_x%d_y%d_y%d", xyRanges.Zoom, xyRanges.XT[0], xyRanges.XT[1], xyRanges.YT[0], xyRanges.YT[1])
}

func (xyRanges XYRanges) Area() Area {
	return Area{
		Tile{X: xyRanges.XT[0], Y: xyRanges.YT[0], Zoom: xyRanges.Zoom}.Point(),
		Tile{X: xyRanges.XT[1] + 1, Y: xyRanges.YT[1] + 1, Zoom: xyRanges.Zoom}.Point(),
	}
}

func (xyRanges XYRanges) Point() Point {
	return Tile{X: xyRanges.XT[0], Y: xyRanges.YT[0], Zoom: xyRanges.Zoom}.Point()
}

// https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames
// left-top corner of the tile

func (tile Tile) Point() Point {
	n := math.Pow(2, float64(tile.Zoom))
	latAngle := math.Atan(math.Sinh(math.Pi * (1 - 2*float64(tile.Y)/n)))

	return Point{
		Lat: Degrees(latAngle * 180 / math.Pi),
		Lon: Degrees(float64(tile.X)/n*360 - 180),
	}
}
