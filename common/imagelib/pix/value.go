package pix

import "image"

type Value = uint8
type ValueSum = uint32  // TODO!!! be careful
type ValueDelta = int16 // TODO!!! be careful

const ValueMax Value = 0xFF
const ValueMiddle Value = 0x7F

type MinMax struct {
	Min Value
	Max Value
}

type Point struct {
	image.Point
	Value
}
