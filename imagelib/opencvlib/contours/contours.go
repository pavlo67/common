package contours

import (
	"github.com/pavlo67/common/common/imagelib"
	"image"

	"github.com/pavlo67/common/common/mathlib/plane"

	"github.com/pavlo67/imagelib/opencvlib"
)

type Contour struct {
	N              int
	TouchingBorder bool
	Points         []image.Point
	Rect           image.Rectangle
}

type Contours []Contour

func (cntr Contour) PointsSub(pMin image.Point) []image.Point {
	imagePoints := make([]image.Point, len(cntr.Points))
	for i, p := range cntr.Points {
		imagePoints[i] = p.Sub(pMin)
	}
	return imagePoints
}

func (cntr Contour) Shorten(distanceMax float64) Contour {
	// log.Printf("SHORTENING: %s / %g", f.CrN(cntr), distanceMax)

NEXT_TRY:
	for {
		for i := 0; i < len(cntr.Points)-2; i++ {
			for j := i + 1; j <= len(cntr.Points); j++ {
				if j == len(cntr.Points) || imagelib.Distance(cntr.Points[i], cntr.Points[j]) > distanceMax {
					if j-i > 2 {
						cntr.Points = append(cntr.Points[:i+1], cntr.Points[j-1:]...)
						continue NEXT_TRY
					}
					break
				}
			}
		}

		// log.Printf("SHORTENED!!! %s", f.CrN(cntr))
		cntr.Rect = imagelib.RectangleAround(0, cntr.Points...)
		return cntr
	}
}

func (cntr Contour) NarrowestApproximation(ratioMin float64) (*plane.Rectangle, error) {
	rectangleMin := opencvlib.RectangleMin(cntr.Points)
	if rectangleMin.HalfSideY != 0 && rectangleMin.HalfSideX/rectangleMin.HalfSideY < ratioMin {
		return nil, nil
	}

	return &rectangleMin, nil
}

//func (cntr Contour) IsCloseTo(cntr1 Contour, distanceMax float64) bool {
//
//	intersect := cntr.RectangleAround.Intersect(imagelib.RectangleExpanded(cntr1.RectangleAround, int(math.Ceil(distanceMax))))
//	if intersect.Min == intersect.Max {
//		return false
//	}
//
//	// TODO!!! reuse rect/pCh generation for cntr
//
//	dist, _, _ := imagelib.PolyChain(cntr.Points).DistanceTo(imagelib.PolyChain(cntr1.Points), distanceMax)
//	return dist <= distanceMax
//}

//func (cntr Contour) ShortenCutting(distanceMax float64) Contour {
//
//NEXT_TRY:
//	for {
//		for i := 0; i < len(cntr.Points); i++ {
//			for j := i + 1; j <= len(cntr.Points); j++ {
//				if j == len(cntr.Points) || imagelib.Distance(cntr.Points[i], cntr.Points[j]) > distanceMax {
//					if j-i > 2 {
//						cntr.Points = append(cntr.Points[:i+1], cntr.Points[j-1:]...)
//						continue NEXT_TRY
//					}
//					break
//				}
//			}
//		}
//
//		// log.Printf("SHORTENED!!! %s", f.CrN(cntr))
//
//		return cntr
//	}
//}

//func (c Contour) ShortenCutting(distanceMax float64) []Contour {
//	var contoursShorten []Contour
//
//	for _, pCh := range ShortenPolyChainCutting(plane.PolyChain(c), distanceMax) {
//		if len(pCh) > 2 || (len(pCh) == 2 && pCh[0].DistanceTo(pCh[1]) > distanceMax) {
//			contoursShorten = append(contoursShorten, Contour(pCh))
//		}
//	}
//
//	return contoursShorten

//func (cntr Contour) ToOuter(fr frame.Frame) Contour {
//	return Contour{
//		N:      cntr.N,
//		Value:  cntr.Value,
//		Points: imagelib.PointsFromPolyChain(fr.PointsToOuter(imagelib.PolyChain(cntr.Points)...)...),
//	}
//}
