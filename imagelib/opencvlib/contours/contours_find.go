package contours

import (
	"github.com/pavlo67/common/common/imagelib"
	"image"

	"gocv.io/x/gocv"

	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/imagelib/layers"
)

const onContoursPV = "on opencvlib.contoursPV()"

func contoursPV(imgGray *image.Gray) (gocv.PointsVector, error) {
	if imgGray == nil {
		return gocv.PointsVector{}, errors.New("imgGray == nil /" + onContoursPV)
	}

	mat, err := gocv.ImageGrayToMatGray(imgGray)
	if err != nil {
		return gocv.PointsVector{}, errors.Wrap(err, onContoursPV)
	}
	defer mat.Close()

	// TODO!!! be careful: gocv.FindContours returns contour points counterclockwise
	return gocv.FindContours(mat, gocv.RetrievalExternal, gocv.ChainApproxSimple), nil

}

const onFind = "on opencvlib.FindContours()"

func Find(lyr *layers.Layer) (Contours, error) {
	if lyr == nil {
		return nil, errors.New("lyr == nil /" + onFind)
	}

	pv, err := contoursPV(&lyr.Gray)
	if err != nil {
		return nil, errors.Wrap(err, onFind)
	}

	var cntrs []Contour

	for i := 0; i < pv.Size(); i++ {
		points := imagelib.ConvertImagePoints(pv.At(i).ToPoints(), false, lyr.Rect.Min, 1)

		// TODO!!! be careful: opencvlib returns contour points in the counter-clockwise order
		pointsConverted := make([]image.Point, len(points))
		for i, p := range imagelib.ConvertImagePoints(points, false, lyr.Rect.Min, 1) {
			pointsConverted[len(pointsConverted)-i-1] = p
		}

		cntrs = append(cntrs, Contour{
			Rect:   imagelib.RectangleAround(0, pointsConverted...),
			Points: pointsConverted,
		})
	}

	return cntrs, nil
}

const onFindContaining = "on areas.FindContaining()"

func FindContaining(lyr *layers.Layer, containingAll bool, pts ...image.Point) ([]Contour, error) {
	if lyr == nil {
		return nil, errors.New("lyr == nil /" + onFindContaining)
	}

	pv, err := contoursPV(&lyr.Gray)
	if err != nil {
		return nil, errors.Wrap(err, onFindContaining)
	}

	var cntrs []Contour

CONTOURS:
	for i := 0; i < pv.Size(); i++ {
		pv := pv.At(i)
		approved := containingAll
		for _, p := range pts {
			inContour := gocv.PointPolygonTest(pv, p.Sub(lyr.Rect.Min), false)
			if containingAll {
				// contour rejected
				if inContour < 0 {
					continue CONTOURS
				}
			} else {
				// contour approved
				if inContour >= 0 {
					approved = true
					break
				}
			}
		}
		if approved {
			cntrs = append(cntrs, Contour{
				// Value:  pix.ValueMax,
				Points: imagelib.ConvertImagePoints(pv.ToPoints(), false, lyr.Rect.Min, 1),
			})
		}
	}

	return cntrs, nil
}

//const onFindContaining = "on areas.FindContaining()"
//
//func FindContaining(lyr *layers.Layer, containingAll bool, pts ...image.Point) (contours.contoursPV, error) {
//	if lyr == nil {
//		return nil, errors.New("lyr == nil /" + onFindContaining)
//	}
//
//	pv, err := contoursPV(&lyr.Gray)
//	if err != nil {
//		return nil, errors.Wrap(err, onFindContaining)
//	}
//
//	var cntrs []contours.Contour
//
//CONTOURS:
//	for i := 0; i < pv.Size(); i++ {
//		pv := pv.At(i)
//		approved := containingAll
//		for _, p := range pts {
//			inContour := gocv.PointPolygonTest(pv, p.Sub(lyr.Rect.Min), false)
//			if containingAll {
//				// contour rejected
//				if inContour < 0 {
//					continue CONTOURS
//				}
//			} else {
//				// contour approved
//				if inContour >= 0 {
//					approved = true
//					break
//				}
//			}
//		}
//		if approved {
//			cntrs = append(cntrs, contours.Contour{
//				Value:  pix.ValueMax,
//				Points: imagelib.ConvertImagePoints(pv.ToPoints(), false, lyr.Rect.Min, 1),
//			})
//		}
//	}
//
//	return cntrs, nil
//}

//const onFindClosing = "on areas.FindClosing()"
//
//func FindClosing(lyr *layers.Layer, inverseColors bool, scale, closeMaskSizePix int) ([]contours.Contour, image.Image, error) {
//
//	if lyr == nil {
//		return nil, nil, errors.New("lyr == nil /" + onFindClosing)
//	}
//	v := pix.ValueMax
//	if inverseColors {
//		lyrInversed, err := lyr.Inversed()
//		if err != nil {
//			return nil, nil, errors.Wrap(err, onFindClosing)
//		} else if lyrInversed == nil {
//			return nil, nil, errors.New("lyrInversed == nil /" + onFindClosing)
//		}
//
//		lyr, v = lyrInversed, 0
//	}
//
//
//	mat, err := gocv.ImageGrayToMatGray(&lyr.Gray)
//	if err != nil {
//		return nil, nil, errors.Wrap(err, onFindClosing)
//	}
//	defer mat.Close()
//
//	matClosed := gocv.NewMat()
//	defer matClosed.Close()
//
//	gocv.MorphologyEx(mat, &matClosed, gocv.MorphClose, gocv.GetStructuringElement(gocv.MorphRect, image.Point{X: closeMaskSizePix, Y: closeMaskSizePix}))
//
//	imgClosed, err := matClosed.ToImage()
//	if err != nil {
//		return nil, nil, errors.Wrap(err, onFindClosing)
//	} else if imgClosedGray, _ := imgClosed.(*image.Gray); imgClosedGray != nil && imgClosedGray.Rect.Min != lyr.Rect.Min {
//		imgClosedGray.Rect = lyr.Rect
//		imgClosed = imgClosedGray
//	}
//
//	// TODO!!! be careful: gocv.FindContours returns contour points counterclockwise
//	pvs := gocv.FindContours(matClosed, gocv.RetrievalExternal, gocv.ChainApproxSimple)
//
//	var contours []contours.Contour
//	for i := 0; i < pvs.Size(); i++ {
//		contours = append(contours, contours.Contour{
//			Value:  v,
//			Points: imagelib.ConvertImagePoints(pvs.At(i).ToPoints(), false, lyr.Rect.Min, scale),
//		})
//	}
//
//	return contours, imgClosed, nil
//}

//func Connect(contours []plane.Contour, distanceMax float64) []plane.Contour {
//	var contoursConnected []plane.Contour
//
//	for i, cntr := range contours {
//
//		for j := i + 1; j < len(contours); j++ {
//			if len(contours[j]) < 1 {
//				continue
//			}
//			pCh := plane.PolyChain(contours[j])
//			for ii, p := range cntr {
//				if dist, pr := plane.GetProjectionOnPolyChain(p, pCh); dist <= distanceMax {
//					toAdd := append(plane.PolyChain{pr.Point2}, pCh[pr.LinkN+1:]...)
//					if pr.pos == 0 {
//						toAdd = append(toAdd, pCh[:pr.LinkN]...)
//					} else {
//						toAdd = append(toAdd, pCh[:pr.LinkN+1]...)
//					}
//					cntr, contours[j] = append(cntr[:ii+1], append(toAdd, cntr[ii+1:]...)...), nil
//					break
//				}
//			}
//		}
//
//		contoursConnected = append(contoursConnected, cntr)
//	}
//
//	return contoursConnected
//}
