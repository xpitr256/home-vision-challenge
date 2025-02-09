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
	fmt.Println(blackAndWhiteImage.Bounds())
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
