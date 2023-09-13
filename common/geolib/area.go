package geolib

import (
	"log"
)

type Area [2]Point

//func (area *Area) Canon() {
//	if area[0].Lat > area[1].Lat {
//		area[0].Lat, area[1].Lat = area[1].Lat, area[0].Lat
//	}
//	if area[0].Lon > area[1].Lon {
//		area[0].Lon, area[1].Lon = area[1].Lon, area[0].Lon
//	}
//}

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

func InAreas([]Area, Area) bool {
	log.Fatal("/crud/locations.InAreas()")
	return false
}

type Ranges struct {
	Lat, Lon [2]Degrees
}

func (area Area) Ranges() Ranges {
	var ranges Ranges

	lat0, lon0 := area[0].Lat, area[0].Lon
	lat1, lon1 := area[1].Lat, area[1].Lon

	if lat0 < lat1 {
		ranges.Lat = [2]Degrees{lat0, lat1}
	} else {
		ranges.Lat = [2]Degrees{lat1, lat0}
	}

	if lon0 < lon1 {
		ranges.Lon = [2]Degrees{lon0, lon1}
	} else {
		ranges.Lon = [2]Degrees{lon1, lon0}
	}

	return ranges
}

func (area Area) Sides() (xSide, ySide float64) {
	ranges := area.Ranges()

	p01 := Point{ranges.Lat[0], ranges.Lon[1]}

	return Point{ranges.Lat[0], ranges.Lon[0]}.DistanceTo(p01), Point{ranges.Lat[1], ranges.Lon[1]}.DistanceTo(p01)
}
