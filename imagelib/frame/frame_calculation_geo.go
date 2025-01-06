package frame

import (
	"fmt"
	"github.com/pavlo67/common/common/imagelib"
	"image"
	"math"

	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/geolib"
	"github.com/pavlo67/common/common/mathlib/plane"
)

type PointRawGeo struct {
	Grid plane.Point2 `json:",omitempty"`
	Geo  geolib.Point `json:",omitempty"`
}

const onCalculateWithGeoPoints = "on frame.CalculateWithGeoPoints()"

// DEPRECATED: use CalculateWithPoints() instead
func CalculateWithGeoPoints(xWidth, yHeight int, ptsRaw [3]PointRawGeo, unitX, unitY float64) (*geolib.Point, plane.XToYAngle, float64, error) {

	rect := image.Rectangle{Max: image.Point{xWidth, yHeight}}
	imgPoint0 := plane.Point2{ptsRaw[0].Grid.X * unitX, ptsRaw[0].Grid.Y * unitY}
	imgPoint1 := plane.Point2{ptsRaw[1].Grid.X * unitX, ptsRaw[1].Grid.Y * unitY}
	imgPoint2 := plane.Point2{ptsRaw[2].Grid.X * unitX, ptsRaw[2].Grid.Y * unitY}

	transformations := []plane.Transformation{
		{To: imagelib.PointFramed(imgPoint0, rect), From: plane.Point2{}},
		{To: imagelib.PointFramed(imgPoint1, rect), From: ptsRaw[0].Geo.DirectionTo(ptsRaw[1].Geo).Moving()},
		{To: imagelib.PointFramed(imgPoint2, rect), From: ptsRaw[0].Geo.DirectionTo(ptsRaw[2].Geo).Moving()},
	}
	outerTriCenter := plane.Center(transformations[0].From, transformations[1].From, transformations[2].From)

	//log.Print(ptsRaw[0].Prev.DirectionTo(ptsRaw[1].Prev), ptsRaw[0].Prev.DirectionTo(ptsRaw[2].Prev))
	//log.Print(transformations[1].From, transformations[2].From)

	rotation, scale, err := plane.CalculateRotationAndScale(transformations, math.Pi/12, 0.2)
	if err != nil {
		return nil, 0, 0, errors.Wrap(err, onCalculateWithGeoPoints)
	} else if !(scale > 0) || math.IsInf(scale, 1) {
		return nil, 0, 0, fmt.Errorf("wrong scale: %f / "+onCalculateWithGeoPoints, scale)
	}

	geoTriCenter := ptsRaw[0].Geo.MovedAt(outerTriCenter)

	// log.Printf("point0 %v moved at %v ---> geoTriCenter: %v", ptsRaw[0].Prev, outerTriCenter, geoTriCenter)

	innerTriCenter := plane.Center(transformations[0].To, transformations[1].To, transformations[2].To)
	triCenterRadius := innerTriCenter.Radius() / scale
	triCenterRotation := innerTriCenter.XToYAngleFromOx() + rotation
	triCenterX, triCenterY := triCenterRadius*math.Cos(float64(triCenterRotation)), triCenterRadius*math.Sin(float64(triCenterRotation))

	geoPoint := geoTriCenter.MovedAt(plane.Point2{-triCenterX, -triCenterY})

	// log.Printf("geoTriCenter %v moved at (%f, %f) ---> geoPoint: %v", geoTriCenter, -triCenterX, -triCenterY, geoPoint)

	return &geoPoint, rotation, scale, nil
}
