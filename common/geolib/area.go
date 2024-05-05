package geolib

import (
	"github.com/pavlo67/common/common/mathlib/plane"
)

func XYRangesAround(geoPoint Point, zoom int, hsX, hsY float64) XYRanges {

	tileMin, tileMax := geoPoint.MovedAt(plane.Point2{-hsX, hsY}).Tile(zoom), geoPoint.MovedAt(plane.Point2{hsX, -hsY}).Tile(zoom)

	if tileMax.X < tileMin.X {
		tileMin.X, tileMax.X = tileMax.X, tileMin.X
	}
	if tileMax.Y < tileMin.Y {
		tileMin.Y, tileMax.Y = tileMax.Y, tileMin.Y
	}

	return XYRanges{Zoom: zoom, XT: XYRange{tileMin.X, tileMax.X}, YT: XYRange{tileMin.Y, tileMax.Y}}

}

func PointInRanges(geoPoint Point, xyRanges XYRanges, tileSide int) plane.Point2 {

	tile := geoPoint.Tile(xyRanges.Zoom)
	leftTop := tile.LeftTop()
	rightBottom := Tile{tile.X + 1, tile.Y + 1, tile.Zoom}.LeftTop()

	x := float64(tileSide) * (float64(tile.X-xyRanges.XT[0]) + float64(geoPoint.Lon-leftTop.Lon)/float64(rightBottom.Lon-leftTop.Lon))
	y := float64(tileSide) * (float64(tile.Y-xyRanges.YT[0]) + float64(geoPoint.Lat-leftTop.Lat)/float64(rightBottom.Lat-leftTop.Lat))

	return plane.Point2{x, y}
}
