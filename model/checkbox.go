package model

import (
	"image"
	"image/color"
)

const CheckboxDefaultSizeInPixels = 24

const (
	checkedStatus   = "checked"
	uncheckedStatus = "unchecked"
)

var (
	greenColor = color.RGBA{R: 0, G: 180, B: 90, A: 255}
	redColor   = color.RGBA{R: 220, G: 50, B: 50, A: 255}
)

type Checkbox struct {
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Status string `json:"status"`
	Box    image.Rectangle
}

func NewCheckbox(box image.Rectangle, image *image.Gray) *Checkbox {
	status := checkedStatus
	isEmpty := isEmptyCheckbox(box, image)
	if isEmpty {
		status = uncheckedStatus
	}
	return &Checkbox{X: box.Min.X, Y: box.Min.Y, Status: status, Box: box}
}

func isEmptyCheckbox(box image.Rectangle, image *image.Gray) bool {
	total := 0
	empties := 0
	// Avoid considering border pixels
	for y := box.Min.Y + 1; y < box.Max.Y-1; y++ {
		for x := box.Min.X + 1; x < box.Max.X-1; x++ {
			total++
			if IsAWhitePosition(x, y, image) {
				empties++
			}
		}
	}
	return (float64(empties) / float64(total) * 100) > 90
}

func IsAWhitePosition(x, y int, image *image.Gray) bool {
	return image.GrayAt(x, y).Y == 255
}

func (c *Checkbox) getColor() color.RGBA {
	if c.Status == checkedStatus {
		return greenColor
	}
	return redColor
}
