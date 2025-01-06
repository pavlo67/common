package contours

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/imagelib"
	"github.com/pavlo67/common/common/imagelib/coloring"
)

func (cntrs Contours) GetMasks(lineWidth int, addInfo bool, fontFile, titleKey string) []imagelib.GetMask {
	getMasks := make([]imagelib.GetMask, len(cntrs))
	for i, contour := range cntrs {
		getMasks[i] = ContourGetMask{
			Contour:    contour,
			LineWidth:  lineWidth,
			ColorNamed: &coloring.RoundAbout[i%len(coloring.RoundAbout)],
			AddInfo:    addInfo,
			FontFile:   fontFile,
			TitleKey:   titleKey,
		}
		titleKey = ""
	}

	return getMasks
}

var _ imagelib.GetMask = &ContourGetMask{}
var _ imagelib.GetMask = ContourGetMask{}

type ContourGetMask struct {
	Contour
	LineWidth  int
	ColorNamed *coloring.ColorNamed
	AddInfo    bool
	FontFile   string
	TitleKey   string
}

func (contourGetMask ContourGetMask) Color() *coloring.ColorNamed {
	return contourGetMask.ColorNamed
}

func (contourGetMask ContourGetMask) Mask(clr color.Color, opts common.Map) imagelib.MasksOneColor {
	pCh := contourGetMask.Contour.Points
	if len(pCh) <= 0 {
		return nil
	}

	maskOneColor := imagelib.MaskOneColor{Color: clr}

	for i := 0; i < len(pCh); i++ {
		line := imagelib.Line(imagelib.Segment(pCh[i], pCh[(i+1)%len(pCh)]), contourGetMask.LineWidth)
		maskOneColor.Points = append(maskOneColor.Points, line...)
	}

	if contourGetMask.FontFile != "" {
		maskOneColor.Marker = &imagelib.MarkerText{
			FontFile: contourGetMask.FontFile,
			Text:     []string{strconv.Itoa(contourGetMask.Contour.N)},
			Point:    pCh[0],
		}
	}

	return imagelib.MasksOneColor{maskOneColor}
}

func (contourGetMask ContourGetMask) Info(colorNamed coloring.ColorNamed) string {
	if !contourGetMask.AddInfo {
		return ""
	}

	cntr := contourGetMask.Contour

	var title, details string

	if contourGetMask.TitleKey != "" {
		title = "\nCONTOURS: " + contourGetMask.TitleKey
	}
	//if contourGetMask.DetailedInfo {
	//	details = fmt.Sprintf("%v", cntr.Points)
	//}
	details = fmt.Sprintf("%v", cntr.Points)

	return title + fmt.Sprintf(
		"\ncntr#%d (%s, points = %d (%s))",
		cntr.N, colorNamed.Name, len(cntr.Points), details)

}
