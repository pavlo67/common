package geolib

type Area [2]Point

func (area Area) XYRanges(zoom int) XYRanges {
	tileMin, tileMax := area[0].Tile(zoom), area[1].Tile(zoom)

	if tileMax.X < tileMin.X {
		tileMin.X, tileMax.X = tileMax.X, tileMin.X
	}
	if tileMax.Y < tileMin.Y {
		tileMin.Y, tileMax.Y = tileMax.Y, tileMin.Y
	}

	return XYRanges{Zoom: zoom, XT: XYRange{tileMin.X, tileMax.X}, YT: XYRange{tileMin.Y, tileMax.Y}}
}
