package service

import (
	"errors"
	"github.com/xpitr256/home-vision-challenge/model"
	"image"
	"image/color"
	"os"
	"path/filepath"
)

const blackDetectionThreshold = 50

func loadTestImage() (image.Image, string, error) {
	filePath := "test/test-image.jpg"
	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", errors.New("error opening image")
	}
	fileName := filepath.Base(filePath)
	defer file.Close()
	formImage, _, err := image.Decode(file)
	if err != nil {
		return nil, "", errors.New("error decoding image")
	}
	return formImage, fileName, nil
}

func convertToBlackAndWhite(img image.Image) (*image.Gray, error) {
	whiteColor := color.Gray{Y: 255}
	blackColor := color.Gray{Y: 0}
	response := image.NewGray(img.Bounds())
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			brightness := uint8((r*299 + g*587 + b*114) / 1000 >> 8)
			// Less than 128 of brightness I considered it as "dark"
			if brightness < 128 {
				response.SetGray(x, y, blackColor)
			} else {
				response.SetGray(x, y, whiteColor)
			}
		}
	}
	return response, nil
}

func isCheckbox(x, y int, formImage *image.Gray, checkboxSizeInPixel int, lastDetected []image.Rectangle) bool {
	edges := &model.Edges{
		Top:    &model.TopEdge{},
		Bottom: &model.BottomEdge{},
		Left:   &model.LeftEdge{},
		Right:  &model.RightEdge{},
	}

	// Delegate the check to the Edges struct
	if !edges.IsStrong(x, y, checkboxSizeInPixel, formImage) {
		return false
	}

	checkbox := image.Rect(x, y, x+checkboxSizeInPixel-1, y+checkboxSizeInPixel-1)

	// Avoid adding overlapping checkboxes
	for _, last := range lastDetected {
		if last.Overlaps(checkbox) {
			return false
		}
	}
	return true
}

func findBoxes(formImage *image.Gray, checkboxSizeInPixel int) []image.Rectangle {
	var lastDetectedCheckboxes []image.Rectangle
	var response []image.Rectangle

	for y := 0; y < formImage.Bounds().Max.Y; y++ {
		for x := 0; x < formImage.Bounds().Max.X; x++ {
			if isCheckbox(x, y, formImage, checkboxSizeInPixel, lastDetectedCheckboxes) {
				checkbox := image.Rect(x, y, x+checkboxSizeInPixel-1, y+checkboxSizeInPixel-1)
				response = append(response, checkbox)

				// Keep track of the last 3 detected checkboxes
				lastDetectedCheckboxes = append(lastDetectedCheckboxes, checkbox)
				if len(lastDetectedCheckboxes) > 3 {
					lastDetectedCheckboxes = lastDetectedCheckboxes[1:]
				}
			}
		}
	}
	return response
}

func removeBlacks(formImage *image.Gray, list []image.Rectangle) []image.Rectangle {
	var result []image.Rectangle
	for _, box := range list {
		total := 0
		blackCount := 0
		for y := box.Min.Y + 1; y < box.Max.Y-1; y++ {
			for x := box.Min.X + 1; x < box.Max.X-1; x++ {
				total++
				if !model.IsAWhitePosition(x, y, formImage) {
					blackCount++
				}
			}
		}
		blackRatio := float64(blackCount) / float64(total) * 100
		if blackRatio < blackDetectionThreshold {
			result = append(result, box)
		}
	}
	return result
}

func getCheckboxesFrom(image *image.Gray, boxes []image.Rectangle) []model.Checkbox {
	response := []model.Checkbox{}
	for _, box := range boxes {
		checkBox := model.NewCheckbox(box, image)
		response = append(response, *checkBox)
	}
	return response
}

func GetCheckboxes(sizeInPixel int) ([]model.Checkbox, string, error) {
	// TODO: Replace with actual image upload from client
	formImage, fileName, err := loadTestImage()
	if err != nil {
		return nil, "", err
	}
	blackAndWhiteImage, err := convertToBlackAndWhite(formImage)
	if err != nil {
		return nil, "", err
	}
	boxes := findBoxes(blackAndWhiteImage, sizeInPixel)
	// Avoid figures with black areas that might be confused with a box
	boxes = removeBlacks(blackAndWhiteImage, boxes)
	checkboxes := getCheckboxesFrom(blackAndWhiteImage, boxes)
	return checkboxes, fileName, nil
}
