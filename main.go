package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Response struct {
	ImageName string `json:"image_name"`
}

func getCheckboxes() (string, error) {
	// TODO obtain the image from client side image upload
	filePath := "test/test-image.jpg"
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("error opening image:", err)
		return "", errors.New("error opening image")
	}
	fileName := filepath.Base(filePath)
	defer file.Close()
	formImage, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("error decoding image:", err)
		return "", errors.New("error decoding image")
	}
	fmt.Println("Image bounds: ", formImage.Bounds())
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
