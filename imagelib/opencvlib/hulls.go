package opencvlib

import (
	"github.com/pavlo67/common/common/imagelib"
	"image"
	"math"

	"gocv.io/x/gocv"

	"github.com/pavlo67/common/common/mathlib/plane"
)

func RectangleMin(points []image.Point) plane.Rectangle {
	rotatedRect := gocv.MinAreaRect(gocv.NewPointVectorFromPoints(points))

	width, height, angle := rotatedRect.Width, rotatedRect.Height, plane.XToYAngle(math.Pi*rotatedRect.Angle/180)
	if width < height {
		width, height, angle = height, width, angle+math.Pi*0.5
	}

	return plane.Rectangle{
		RectangleXY: plane.RectangleXY{
			Point2:    imagelib.PolyChain([]image.Point{rotatedRect.Center})[0],
			HalfSideX: float64(width) * 0.5,
			HalfSideY: float64(height) * 0.5,
		},
		XToYAngle: angle.Canon(),
	}
}

func ConvexHull(points []image.Point) []int {
	matHull := gocv.NewMat()
	defer matHull.Close()

	gocv.ConvexHull(gocv.NewPointVectorFromPoints(points), &matHull, false, false)

	// TODO: be careful!!!
	// if !returnPoints: matHull.T() == CV32S   (== []int == list of contour indices, required for gocv.ConvexityDefects())
	// if  returnPoints: matHull.T() == CV32SC2 (== []image.Point???)

	convexHullIndices := make([]int, matHull.Size()[0])

	for i := 0; i < len(convexHullIndices); i++ {
		convexHullIndices[i] = int(matHull.GetIntAt(i, 0))
		// convexHullIndices, convexHullPoints := make([]int, hullLength), make([]image.Point, hullLength)

		if i > 0 && convexHullIndices[i] < convexHullIndices[i-1] {
			// The convex hull indices are not monotonous, which can be in the case when the input contour contains self-intersections
		}
	}

	// pr.ConvexHullArea = gocv.ContourArea(gocv.NewPointVectorFromPoints(pr.ConvexHullPoints)) / (ls.DPM * ls.DPM)

	return convexHullIndices
}
