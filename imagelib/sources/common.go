package sources

import (
	"github.com/pavlo67/common/common/geolib"
	"github.com/pavlo67/common/common/mathlib/plane"
)

type TestInfo struct {
	NFrom, NTo int
	NDivider   int
}

const TestInfoFilename = "test_info.json"
const VideoInfoFilename = "video_info.json"
const FramesDescriptionsFilename = "frames.jlist"
const FramesAllDescriptionsFilename = "frames_all.jlist"
const RePGNStr = `^(\d{4})\.pgm$`
const RePNGStr = `^(\d{4})\.png$`
const RePGNPNGStr = `^(\d{4})\.(?:pgm|png)$`

func InterpolatedDescriptions(descr, descrNext Description) []Description {
	var descrsAll []Description
	if divider := descrNext.N - descr.N; divider > 1 {
		moving := descr.GeoPoint.DirectionTo(*descrNext.GeoPoint).Moving()
		stepDX, stepDY := moving.X/float64(divider), moving.Y/float64(divider)

		// TODO!!! be careful: no .Canon(); no .XToYAngle (it's opposite to Bearing)
		stepBearing := (descrNext.Bearing - descr.Bearing) / geolib.Bearing(divider)

		stepDPM := (descrNext.DPM - descr.DPM) / float64(divider)
		for j := 1; j < divider; j++ {
			geoPointRef := new(geolib.Point)
			*geoPointRef = descr.GeoPoint.MovedAt(plane.Point2{float64(j) * stepDX, float64(j) * stepDY})
			descrsAll = append(descrsAll, Description{
				N:        descr.N + j,
				GeoPoint: geoPointRef,
				Bearing:  (descr.Bearing + stepBearing*geolib.Bearing(j)).Canon(),
				DPM:      descr.DPM + stepDPM*float64(j),
			})
		}
	}

	return descrsAll
}
