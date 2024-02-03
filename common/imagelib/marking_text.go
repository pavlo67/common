package imagelib

import (
	"image"
	"image/color"
	"log"

	"golang.org/x/image/draw"

	"github.com/pavlo67/common/common/mathlib/plane"
)

var _ Marker = &MarkerText{}
var _ Marker = MarkerText{}

type MarkerText struct {
	DPI      float64
	Size     float64
	Spacing  float64
	FontFile string
	Text     []string
	image.Point
}

func (mt MarkerText) Mark(drawImg draw.Image, clr color.Color) {
	if _, err := Write(drawImg, mt.Point, mt.DPI, mt.Size, mt.Spacing, mt.FontFile, clr, mt.Text); err != nil {
		log.Printf("ERROR: on MarkerText.Mark(): %s", err)
	}
}

// ----------------------------------------------------------------------------------------

//var _ GetMask = &TextGetMask{}
//var _ GetMask = TextGetMask{}
//
//type TextGetMask struct {
//	*ColorNamed
//	Point     image.Point
//	Segments  []plane.Segment
//	LineWidth int
//	FontFile  string
//	Label     string
//	Title     string
//	Text      string
//}
//
//func (textGetMask TextGetMask) Color() *ColorNamed {
//	return textGetMask.ColorNamed
//}
//
//func (textGetMask TextGetMask) Mask(clr color.Color) Mask {
//	if textGetMask.FontFile == "" {
//		return nil
//	}
//
//	var points []image.Point
//	for _, segment := range textGetMask.Segments {
//		points = append(points, Line(segment, textGetMask.LineWidth)...)
//	}
//
//	return Mask{
//		MaskOneColor{
//			Color:  clr,
//			Points: points,
//			Marker: &MarkerText{
//				FontFile: textGetMask.FontFile,
//				Text:     []string{textGetMask.Label},
//				Point:    textGetMask.Point,
//			},
//		},
//	}
//}
//
//func (textGetMask TextGetMask) Info(colorNamed ColorNamed) string {
//
//	var title string
//	if textGetMask.Title != "" {
//		title = "\n" + textGetMask.Title + "\n"
//	}
//
//	return title + textGetMask.Text
//}

var _ GetMask = &SegmentsGetMask{}
var _ GetMask = SegmentsGetMask{}

type SegmentsGetMask struct {
	*ColorNamed
	Point     image.Point
	Segments  []plane.Segment
	LineWidth int
	FontFile  string
	Label     string
	Title     string
	Text      string
}

func (segmentsGetMask SegmentsGetMask) Color() *ColorNamed {
	return segmentsGetMask.ColorNamed
}

func (segmentsGetMask SegmentsGetMask) Mask(clr color.Color) Mask {
	if segmentsGetMask.FontFile == "" {
		return nil
	}

	var points []image.Point
	for _, segment := range segmentsGetMask.Segments {
		points = append(points, Line(segment, segmentsGetMask.LineWidth)...)
	}
	xMin, xMax := segmentsGetMask.Point.X-segmentsGetMask.LineWidth-2, segmentsGetMask.Point.X+segmentsGetMask.LineWidth+3
	yMin, yMax := segmentsGetMask.Point.Y-segmentsGetMask.LineWidth-2, segmentsGetMask.Point.Y+segmentsGetMask.LineWidth+3

	points = append(points, FilledRectangle(xMin, xMax, yMin, yMax)...)

	return Mask{
		MaskOneColor{
			Color:  clr,
			Points: points,
			Marker: &MarkerText{
				FontFile: segmentsGetMask.FontFile,
				Text:     []string{segmentsGetMask.Label},
				Point:    segmentsGetMask.Point,
			},
		},
	}
}

func (segmentsGetMask SegmentsGetMask) Info(colorNamed ColorNamed) string {

	var title string
	if segmentsGetMask.Title != "" {
		title = "\n" + segmentsGetMask.Title + "\n"
	}

	return title + segmentsGetMask.Text
}

func FilledRectangle(xMin, xMax, yMin, yMax int) []image.Point {
	points := make([]image.Point, 0, max(0, (xMax-xMin)*(yMax-yMin)))

	for x := xMin; x < xMax; x++ {
		for y := yMin; y < yMax; y++ {
			points = append(points, image.Point{x, y})
		}
	}

	return points
}
