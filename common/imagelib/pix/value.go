package pix

import "image"

type Value = uint8
type ValueDelta = int16 // TODO!!! be careful
type ValueSum = int32   // TODO!!! be careful
type ValueSumSigned = ValueSum

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
