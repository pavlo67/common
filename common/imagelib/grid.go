package imagelib

import (
	"fmt"
	"image/color"

	"github.com/pavlo67/common/common/mathlib/plane"

	"golang.org/x/image/draw"
)

type GridMarking struct {
	DividerX   int
	DividerY   int
	GridPoints []plane.Point2
}

func CrossMark(img draw.Image, x, y, size int, clr color.Color) {
	AddVLine(img, x, y-size, y+size, clr)
	AddHLine(img, x-size, y, x+size, clr)
}

const onAddGrid = "on imagelib.MarkGrid()"

func MarkGrid(img draw.Image, gridMarking GridMarking, clrGrid, clrPoints color.Color, crossSize int) error {
	dividerX, dividerY := gridMarking.DividerX, gridMarking.DividerY

	if dividerX <= 0 || dividerY <= 0 {
		return fmt.Errorf("wrong dividers: dividerX = %d, dividerY = %d / "+onAddGrid, dividerX, dividerY)
	}

	rect := img.Bounds()
	size := rect.Size()

	dx, dy := size.X/dividerX, size.Y/dividerY
	if dx <= 0 || dy <= 0 {
		return fmt.Errorf("wrong dividers: dividerX = %d, dividerY = %d on size = %v / "+onAddGrid, dividerX, dividerY, size)
	}

	for x := rect.Min.X; x < rect.Max.X; x += dx {
		AddVLine(img, x, rect.Min.Y, rect.Max.Y-1, clrGrid)
	}
	for y := rect.Min.Y; y < rect.Max.Y; y += dy {
		AddHLine(img, rect.Min.X, y, rect.Max.X-1, clrGrid)
	}

	for _, gp := range gridMarking.GridPoints {

		CrossMark(img, int(float64(dx)*gp.X), int(float64(dy)*gp.Y), crossSize, clrPoints)
	}

	return nil
}
