package model

import "image"

const CheckboxDefaultSizeInPixels = 24

type Checkbox struct {
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Status string `json:"status"`
}

func NewCheckbox(box image.Rectangle, image *image.Gray) *Checkbox {
	status := "checked"
	isEmpty := isEmptyCheckbox(box, image)
	if isEmpty {
		status = "unchecked"
	}
	return &Checkbox{X: box.Min.X, Y: box.Min.Y, Status: status}
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
