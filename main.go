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
)

type Response struct {
	ImageName string `json:"image_name"`
}

const checkboxSizeInPixel = 25

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
	return strength > 90
}

func isCheckbox(x, y int, formImage *image.Gray, checkboxSizeInPixel int) bool {
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

	return true
}

func findBoxes(formImage *image.Gray) []image.Rectangle {
	var response []image.Rectangle
	for y := 0; y < formImage.Bounds().Max.Y; y++ {
		for x := 0; x < formImage.Bounds().Max.X; x++ {
			if isCheckbox(x, y, formImage, checkboxSizeInPixel) {
				checkbox := image.Rect(x, y, x+checkboxSizeInPixel-1, y+checkboxSizeInPixel-1)
				response = append(response, checkbox)
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

func getCheckboxes() (string, error) {
	// TODO obtain the image from client side image upload
	formImage, fileName, err := loadTestImage()
	if err != nil {
		fmt.Println("error loading test image:", err)
		return "", err
	}
	blackAndWhiteImage, err := convertToBlackAndWhite(formImage)
	if err != nil {
		fmt.Println("error converting image to black and white:", err)
		return "", err
	}
	boxes := findBoxes(blackAndWhiteImage)
	fmt.Println("Boxes length: ", len(boxes))
	return fileName, nil
}

func checkboxHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	imageName, err := getCheckboxes()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Processing checkboxes"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{ImageName: imageName})

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
