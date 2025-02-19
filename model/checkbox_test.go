package model

import (
	"image"
	"image/color"
	"testing"
)

func createTestGrayImage() *image.Gray {
	img := image.NewGray(image.Rect(0, 0, 100, 100))

	// White region (10,10)-(30,30)
	for y := 10; y < 30; y++ {
		for x := 10; x < 30; x++ {
			img.SetGray(x, y, color.Gray{Y: 255}) // White pixels
		}
	}

	// Black region (40,40)-(60,60)
	for y := 40; y < 60; y++ {
		for x := 40; x < 60; x++ {
			img.SetGray(x, y, color.Gray{Y: 0}) // Black pixels
		}
	}

	return img
}

func TestNewCheckboxWithMore90WhiteShouldCreateUncheckedCheckbox(t *testing.T) {
	img := createTestGrayImage()

	// Checkbox is empty (more than 90% white pixels)
	box := image.Rect(10, 10, 30, 30)
	checkbox := NewCheckbox(box, img)
	if checkbox.Status != uncheckedStatus {
		t.Errorf("Expected status: %s, got: %s", uncheckedStatus, checkbox.Status)
	}
}

func TestNewCheckboxWithLess90WhiteShouldCreateCheckedCheckbox(t *testing.T) {
	img := createTestGrayImage()

	// Checkbox is checked (less than 90% white pixels)
	box := image.Rect(40, 40, 60, 60)
	checkbox := NewCheckbox(box, img)
	if checkbox.Status != checkedStatus {
		t.Errorf("Expected status: %s, got: %s", checkedStatus, checkbox.Status)
	}
}

func TestIsAWhitePositionShouldReturnTrueForWhiteCoords(t *testing.T) {
	img := createTestGrayImage()
	if !IsAWhitePosition(10, 10, img) {
		t.Error("Expected pixel to be white")
	}
}

func TestIsAWhitePositionShouldReturnFalseForBlackCoords(t *testing.T) {
	img := createTestGrayImage()
	if IsAWhitePosition(45, 45, img) {
		t.Error("Expected pixel not to be white")
	}
}

func TestCheckedCheckboxShouldReturnGreenWhenCallingGetColor(t *testing.T) {
	checkbox := &Checkbox{Status: checkedStatus}
	if checkbox.getColor() != greenColor {
		t.Errorf("Expected color: %v, got: %v", greenColor, checkbox.getColor())
	}
}

func TestUnCheckedCheckboxShouldReturnRedWhenCallingGetColor(t *testing.T) {
	checkbox := &Checkbox{Status: uncheckedStatus}
	if checkbox.getColor() != redColor {
		t.Errorf("Expected color: %v, got: %v", redColor, checkbox.getColor())
	}
}
