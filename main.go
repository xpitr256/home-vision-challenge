package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type CheckboxResponse struct {
	ImageName       string     `json:"image_name"`
	TotalDetections int        `json:"total_detections"`
	Checkboxes      []Checkbox `json:"checkboxes"`
	SizeInPixels    int        `json:"size_in_pixels"`
}

const (
	CheckboxDefaultSizeInPixels = 24
	blackDetectionThreshold     = 50
)

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

func isAWhitePosition(x, y int, image *image.Gray) bool {
	return image.GrayAt(x, y).Y == 255
}

func isEmptyCheckbox(box image.Rectangle, image *image.Gray) bool {
	total := 0
	empties := 0
	// Avoid considering border pixels
	for y := box.Min.Y + 1; y < box.Max.Y-1; y++ {
		for x := box.Min.X + 1; x < box.Max.X-1; x++ {
			total++
			if isAWhitePosition(x, y, image) {
				empties++
			}
		}
	}
	return (float64(empties) / float64(total) * 100) > 90
}

// Edge is the interface that all edge types will implement
type Edge interface {
	IsStrong(x, y, size int, img *image.Gray) bool
}

type TopEdge struct{}

func (te *TopEdge) IsStrong(x, y, size int, img *image.Gray) bool {
	totalPixels, blackPixels := 0, 0
	for i := 0; i < size; i++ {
		if x+i < img.Bounds().Max.X && img.GrayAt(x+i, y).Y == 0 {
			blackPixels++
		}
		totalPixels++
	}
	if totalPixels == 0 {
		return false
	}
	return isStrongBoundary(float64(blackPixels) / float64(totalPixels) * 100)
}

type BottomEdge struct{}

func (be *BottomEdge) IsStrong(x, y, size int, img *image.Gray) bool {
	totalPixels, blackPixels := 0, 0
	for i := 0; i < size; i++ {
		if x+i < img.Bounds().Max.X && (img.GrayAt(x+i, y+size-1).Y == 0 || img.GrayAt(x+i, y+size-2).Y == 0) {
			blackPixels++
		}
		totalPixels++
	}
	if totalPixels == 0 {
		return false
	}
	return isStrongBoundary(float64(blackPixels) / float64(totalPixels) * 100)
}

type LeftEdge struct{}

func (le *LeftEdge) IsStrong(x, y, size int, img *image.Gray) bool {
	totalPixels, blackPixels := 0, 0
	for i := 0; i < size; i++ {
		if y+i < img.Bounds().Max.Y && img.GrayAt(x, y+i).Y == 0 {
			blackPixels++
		}
		totalPixels++
	}
	if totalPixels == 0 {
		return false
	}
	return isStrongBoundary(float64(blackPixels) / float64(totalPixels) * 100)
}

type RightEdge struct{}

func (re *RightEdge) IsStrong(x, y, size int, img *image.Gray) bool {
	totalPixels, blackPixels := 0, 0
	for i := 0; i < size; i++ {
		if y+i < img.Bounds().Max.Y && (img.GrayAt(x+size-1, y+i).Y == 0 || img.GrayAt(x+size-2, y+i).Y == 0) {
			blackPixels++
		}
		totalPixels++
	}
	if totalPixels == 0 {
		return false
	}
	return isStrongBoundary(float64(blackPixels) / float64(totalPixels) * 100)
}

type Edges struct {
	Top    Edge
	Bottom Edge
	Left   Edge
	Right  Edge
}

func (e *Edges) IsStrong(x, y, size int, img *image.Gray) bool {
	return e.Top.IsStrong(x, y, size, img) && e.Bottom.IsStrong(x, y, size, img) && e.Left.IsStrong(x, y, size, img) && e.Right.IsStrong(x, y, size, img)
}

func isStrongBoundary(strength float64) bool {
	return strength > 82
}

func isCheckbox(x, y int, formImage *image.Gray, checkboxSizeInPixel int, lastDetected []image.Rectangle) bool {
	edges := &Edges{
		Top:    &TopEdge{},
		Bottom: &BottomEdge{},
		Left:   &LeftEdge{},
		Right:  &RightEdge{},
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
func getCheckboxesFrom(image *image.Gray, boxes []image.Rectangle) []Checkbox {
	response := []Checkbox{}
	for _, box := range boxes {
		checkBox := NewCheckbox(box, image)
		response = append(response, *checkBox)
	}
	return response
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

func removeBlacks(formImage *image.Gray, list []image.Rectangle) []image.Rectangle {
	var result []image.Rectangle
	for _, box := range list {
		total := 0
		blackCount := 0
		for y := box.Min.Y + 1; y < box.Max.Y-1; y++ {
			for x := box.Min.X + 1; x < box.Max.X-1; x++ {
				total++
				if !isAWhitePosition(x, y, formImage) {
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

func getCheckboxSize(r *http.Request) (int, error) {
	checkboxSizeInPixels := CheckboxDefaultSizeInPixels
	checkboxSizeInPixelsStr := r.URL.Query().Get("size")
	if checkboxSizeInPixelsStr != "" {
		converted, err := strconv.Atoi(checkboxSizeInPixelsStr)
		if err != nil {
			return 0, errors.New("invalid 'size' parameter, must be an integer between 1 and 200")
		}
		checkboxSizeInPixels = converted
	}
	if checkboxSizeInPixels <= 0 || checkboxSizeInPixels > 200 {
		return 0, errors.New("'size' parameter must be between 1 and 200")
	}
	return checkboxSizeInPixels, nil
}

func getCheckboxes(sizeInPixel int) ([]Checkbox, string, error) {
	// TODO obtain the image from client side image upload
	formImage, fileName, err := loadTestImage()
	if err != nil {
		fmt.Println("error loading test image:", err)
		return nil, "", err
	}
	blackAndWhiteImage, err := convertToBlackAndWhite(formImage)
	if err != nil {
		fmt.Println("error converting image to black and white:", err)
		return nil, "", err
	}
	boxes := findBoxes(blackAndWhiteImage, sizeInPixel)
	// Avoid figures with black areas that might be confused with a box
	boxes = removeBlacks(blackAndWhiteImage, boxes)
	checkboxes := getCheckboxesFrom(blackAndWhiteImage, boxes)
	return checkboxes, fileName, nil
}

func checkboxHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	checkboxSizeInPixels, err := getCheckboxSize(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	checkboxes, imageName, err := getCheckboxes(checkboxSizeInPixels)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Processing checkboxes"})
		return
	}
	w.WriteHeader(http.StatusOK)
	response := CheckboxResponse{
		ImageName:       imageName,
		TotalDetections: len(checkboxes),
		Checkboxes:      checkboxes,
		SizeInPixels:    checkboxSizeInPixels,
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/checkbox", checkboxHandler)

	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
