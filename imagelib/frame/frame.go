package frame

import (
	"image"
	"math"

	"github.com/pavlo67/common/common/imagelib"
	"github.com/pavlo67/common/common/imagelib/coloring"
	"github.com/pavlo67/common/common/imagelib/pix"
	"github.com/pavlo67/common/common/mathlib"
	"github.com/pavlo67/common/common/mathlib/plane"
)

type ValueRGBA [coloring.NumColorsRGBA]pix.Value

var _ imagelib.Described = &LayerRGBA{}

type LayerRGBA struct {
	image.RGBA
	imagelib.Settings
}

func (lyrRGB LayerRGBA) Description() imagelib.Settings {
	return lyrRGB.Settings
}

type Frame struct {
	LayerRGBA
	plane.Position
}

// TODO!!! be careful:
// the origin of image.Rectangle is in the left top corner, Ox to right, Oy to bottom
// the geometry origin is in the rectangle center (when frame.pos is zero), geometry Ox to right, geometry Oy to top

func (fr Frame) PlaneRectangle() plane.Rectangle {
	var halfSideX, halfSideY float64

	if fr.DPM > 0 && !math.IsInf(fr.DPM, 1) {
		rect := fr.RGBA.Rect
		halfSideX, halfSideY = max(0, float64(rect.Max.X-rect.Min.X)/(2*fr.DPM)), max(0, float64(rect.Max.Y-rect.Min.Y)/(2*fr.DPM))
	}

	return plane.Rectangle{RectangleXY: plane.RectangleXY{fr.Point2, halfSideX, halfSideY}, XToYAngle: fr.XToYAngle}
}

func (fr Frame) PointsToOuter(pChInner ...plane.Point2) plane.PolyChain {
	rect := fr.RGBA.Rect
	center := plane.Point2{0.5 * float64(rect.Min.X+rect.Max.X-1), 0.5 * float64(rect.Min.Y+rect.Max.Y-1)}
	pChOuter := make(plane.PolyChain, len(pChInner))

	if !(fr.DPM > 0) || math.IsInf(fr.DPM, 1) {
		for i := range pChInner {
			pChOuter[i].X, pChOuter[i].Y = math.NaN(), math.NaN()
		}

	} else {
		for i, p := range pChInner {
			radiusOuter := math.Sqrt((p.X-center.X)*(p.X-center.X)+(p.Y-center.Y)*(p.Y-center.Y)) / fr.DPM
			if radiusOuter <= mathlib.EPS {
				pChOuter[i] = fr.Point2

			} else {
				angleInner := plane.Point2{p.X - center.X, p.Y - center.Y}.XToYAngleFromOx()
				angleOuter := fr.XToYAngle - angleInner
				pChOuter[i] = plane.Point2{fr.Point2.X + radiusOuter*math.Cos(float64(angleOuter)), fr.Point2.Y + radiusOuter*math.Sin(float64(angleOuter))}
			}
		}
	}

	return pChOuter
}

func (fr Frame) PointToInner(p2Outer plane.Point2) plane.Point2 {
	radiusInner := math.Sqrt((p2Outer.X-fr.X)*(p2Outer.X-fr.X)+(p2Outer.Y-fr.Y)*(p2Outer.Y-fr.Y)) * fr.DPM
	var angleOuter plane.XToYAngle
	if radiusInner > mathlib.EPS {
		angleOuter = plane.Point2{p2Outer.X - fr.X, p2Outer.Y - fr.Y}.XToYAngleFromOx()
	}
	angleInner := fr.XToYAngle - angleOuter
	rect := fr.RGBA.Rect
	center := plane.Point2{0.5 * float64(rect.Min.X+rect.Max.X-1), 0.5 * float64(rect.Min.Y+rect.Max.Y-1)}

	return plane.Point2{center.X + radiusInner*math.Cos(float64(angleInner)), center.Y + radiusInner*math.Sin(float64(angleInner))}
}

// MovingToInner calculates the frame moving over fixed image (inner --> outer)
func (fr Frame) MovingToInner(movingOuter plane.Point2) plane.Point2 {
	if !(fr.DPM > 0) || math.IsInf(fr.DPM, 1) {
		return plane.Point2{math.NaN(), math.NaN()}
	}

	movingRadiusInner := math.Sqrt(movingOuter.X*movingOuter.X+movingOuter.Y*movingOuter.Y) * fr.DPM
	movingAngleInner := fr.XToYAngle

	if movingRadiusInner > mathlib.EPS {
		movingAngleInner -= movingOuter.XToYAngleFromOx()
	}

	return plane.Point2{movingRadiusInner * math.Cos(float64(movingAngleInner)), movingRadiusInner * math.Sin(float64(movingAngleInner))}
}

// MovingToOuter calculates the frame moving over fixed image (outer --> inner)
func (fr Frame) MovingToOuter(movingInner plane.Point2) plane.Point2 {
	if !(fr.DPM > 0) || math.IsInf(fr.DPM, 1) {
		return plane.Point2{math.NaN(), math.NaN()}
	}

	movingRadiusOuter := math.Sqrt(movingInner.X*movingInner.X+movingInner.Y*movingInner.Y) / fr.DPM
	movingAngleOuter := fr.XToYAngle
	if movingRadiusOuter > mathlib.EPS {
		movingAngleOuter -= movingInner.XToYAngleFromOx()
	}

	return plane.Point2{movingRadiusOuter * math.Cos(float64(movingAngleOuter)), movingRadiusOuter * math.Sin(float64(movingAngleOuter))}
}
